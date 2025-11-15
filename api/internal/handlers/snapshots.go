package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/streamspace/streamspace/api/internal/db"
	"github.com/streamspace/streamspace/api/internal/k8s"
)

// SnapshotsHandler handles session snapshot and restore operations
type SnapshotsHandler struct {
	db *db.Database
}

// NewSnapshotsHandler creates a new snapshots handler
func NewSnapshotsHandler(database *db.Database) *SnapshotsHandler {
	return &SnapshotsHandler{
		db: database,
	}
}

// Snapshot represents a session snapshot
type Snapshot struct {
	ID           string                 `json:"id"`
	SessionID    string                 `json:"sessionId"`
	UserID       string                 `json:"userId"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description,omitempty"`
	Type         string                 `json:"type"`   // manual, automatic, scheduled
	Status       string                 `json:"status"` // creating, available, restoring, failed, deleted
	StoragePath  string                 `json:"storagePath,omitempty"`
	SizeBytes    int64                  `json:"sizeBytes"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt    time.Time              `json:"createdAt"`
	CompletedAt  *time.Time             `json:"completedAt,omitempty"`
	ExpiresAt    *time.Time             `json:"expiresAt,omitempty"`
	ErrorMessage string                 `json:"errorMessage,omitempty"`
}

// RegisterRoutes registers snapshot routes
func (h *SnapshotsHandler) RegisterRoutes(router *gin.RouterGroup) {
	snapshots := router.Group("/sessions/:sessionId/snapshots")
	{
		// Snapshot management
		snapshots.GET("", h.ListSnapshots)
		snapshots.POST("", h.CreateSnapshot)
		snapshots.GET("/:snapshotId", h.GetSnapshot)
		snapshots.DELETE("/:snapshotId", h.DeleteSnapshot)

		// Restore operations
		snapshots.POST("/:snapshotId/restore", h.RestoreSnapshot)
		snapshots.GET("/:snapshotId/restore/status", h.GetRestoreStatus)

		// Snapshot configuration
		snapshots.GET("/config", h.GetSnapshotConfig)
		snapshots.PUT("/config", h.UpdateSnapshotConfig)
	}

	// User's all snapshots across sessions
	router.GET("/snapshots", h.ListAllUserSnapshots)
	router.GET("/snapshots/stats", h.GetSnapshotStats)
}

// ListSnapshots returns all snapshots for a session
func (h *SnapshotsHandler) ListSnapshots(c *gin.Context) {
	sessionID := c.Param("sessionId")
	userID, _ := c.Get("userID")
	userIDStr := userID.(string)

	ctx := context.Background()

	// Verify session ownership
	if !h.verifySessionOwnership(ctx, sessionID, userIDStr) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this session"})
		return
	}

	rows, err := h.db.DB().QueryContext(ctx, `
		SELECT id, session_id, user_id, name, description, type, status, storage_path,
		       size_bytes, metadata, created_at, completed_at, expires_at, error_message
		FROM session_snapshots
		WHERE session_id = $1
		ORDER BY created_at DESC
	`, sessionID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list snapshots"})
		return
	}
	defer rows.Close()

	snapshots := []Snapshot{}
	for rows.Next() {
		var s Snapshot
		var description, storagePath, errorMessage sql.NullString
		var completedAt, expiresAt sql.NullTime
		var metadataJSON []byte

		if err := rows.Scan(&s.ID, &s.SessionID, &s.UserID, &s.Name, &description, &s.Type, &s.Status, &storagePath, &s.SizeBytes, &metadataJSON, &s.CreatedAt, &completedAt, &expiresAt, &errorMessage); err == nil {
			if description.Valid {
				s.Description = description.String
			}
			if storagePath.Valid {
				s.StoragePath = storagePath.String
			}
			if errorMessage.Valid {
				s.ErrorMessage = errorMessage.String
			}
			if completedAt.Valid {
				s.CompletedAt = &completedAt.Time
			}
			if expiresAt.Valid {
				s.ExpiresAt = &expiresAt.Time
			}
			if len(metadataJSON) > 0 {
				json.Unmarshal(metadataJSON, &s.Metadata)
			}

			snapshots = append(snapshots, s)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"snapshots": snapshots,
		"count":     len(snapshots),
		"sessionId": sessionID,
	})
}

// CreateSnapshot creates a new snapshot of a session
func (h *SnapshotsHandler) CreateSnapshot(c *gin.Context) {
	sessionID := c.Param("sessionId")
	userID, _ := c.Get("userID")
	userIDStr := userID.(string)

	var req struct {
		Name        string                 `json:"name" binding:"required"`
		Description string                 `json:"description"`
		Type        string                 `json:"type"`      // manual, automatic
		ExpiresIn   string                 `json:"expiresIn"` // duration like "7d", "30d", "90d"
		Metadata    map[string]interface{} `json:"metadata"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()

	// Verify session ownership
	if !h.verifySessionOwnership(ctx, sessionID, userIDStr) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this session"})
		return
	}

	// Default to manual type
	if req.Type == "" {
		req.Type = "manual"
	}

	// Calculate expiration
	var expiresAt *time.Time
	if req.ExpiresIn != "" {
		duration, err := time.ParseDuration(req.ExpiresIn)
		if err == nil {
			expiry := time.Now().Add(duration)
			expiresAt = &expiry
		}
	}

	snapshotID := fmt.Sprintf("snap_%s_%d", sessionID, time.Now().UnixNano())
	metadataJSON, _ := json.Marshal(req.Metadata)

	// Get storage path
	storagePath := h.getSnapshotStoragePath(sessionID, snapshotID)

	_, err := h.db.DB().ExecContext(ctx, `
		INSERT INTO session_snapshots (id, session_id, user_id, name, description, type, status, storage_path, metadata, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, 'creating', $7, $8, $9)
	`, snapshotID, sessionID, userIDStr, req.Name, req.Description, req.Type, storagePath, metadataJSON, expiresAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create snapshot"})
		return
	}

	// Trigger async snapshot creation
	go h.createSnapshotAsync(snapshotID, sessionID, storagePath)

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Snapshot creation initiated",
		"snapshotId": snapshotID,
		"status":     "creating",
	})
}

// GetSnapshot retrieves a specific snapshot
func (h *SnapshotsHandler) GetSnapshot(c *gin.Context) {
	sessionID := c.Param("sessionId")
	snapshotID := c.Param("snapshotId")
	userID, _ := c.Get("userID")
	userIDStr := userID.(string)

	ctx := context.Background()

	// Verify session ownership
	if !h.verifySessionOwnership(ctx, sessionID, userIDStr) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this session"})
		return
	}

	var s Snapshot
	var description, storagePath, errorMessage sql.NullString
	var completedAt, expiresAt sql.NullTime
	var metadataJSON []byte

	err := h.db.DB().QueryRowContext(ctx, `
		SELECT id, session_id, user_id, name, description, type, status, storage_path,
		       size_bytes, metadata, created_at, completed_at, expires_at, error_message
		FROM session_snapshots
		WHERE id = $1 AND session_id = $2
	`, snapshotID, sessionID).Scan(&s.ID, &s.SessionID, &s.UserID, &s.Name, &description, &s.Type, &s.Status, &storagePath, &s.SizeBytes, &metadataJSON, &s.CreatedAt, &completedAt, &expiresAt, &errorMessage)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Snapshot not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get snapshot"})
		return
	}

	if description.Valid {
		s.Description = description.String
	}
	if storagePath.Valid {
		s.StoragePath = storagePath.String
	}
	if errorMessage.Valid {
		s.ErrorMessage = errorMessage.String
	}
	if completedAt.Valid {
		s.CompletedAt = &completedAt.Time
	}
	if expiresAt.Valid {
		s.ExpiresAt = &expiresAt.Time
	}
	if len(metadataJSON) > 0 {
		json.Unmarshal(metadataJSON, &s.Metadata)
	}

	c.JSON(http.StatusOK, s)
}

// DeleteSnapshot deletes a snapshot
func (h *SnapshotsHandler) DeleteSnapshot(c *gin.Context) {
	sessionID := c.Param("sessionId")
	snapshotID := c.Param("snapshotId")
	userID, _ := c.Get("userID")
	userIDStr := userID.(string)

	ctx := context.Background()

	// Verify session ownership
	if !h.verifySessionOwnership(ctx, sessionID, userIDStr) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this session"})
		return
	}

	// Get storage path before deleting
	var storagePath sql.NullString
	h.db.DB().QueryRowContext(ctx, `SELECT storage_path FROM session_snapshots WHERE id = $1`, snapshotID).Scan(&storagePath)

	// Delete from database
	_, err := h.db.DB().ExecContext(ctx, `
		UPDATE session_snapshots
		SET status = 'deleted', updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND session_id = $2
	`, snapshotID, sessionID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete snapshot"})
		return
	}

	// Delete physical files asynchronously
	if storagePath.Valid && storagePath.String != "" {
		go h.deleteSnapshotFiles(storagePath.String)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Snapshot deleted",
		"snapshotId": snapshotID,
	})
}

// RestoreSnapshot restores a session from a snapshot
func (h *SnapshotsHandler) RestoreSnapshot(c *gin.Context) {
	sessionID := c.Param("sessionId")
	snapshotID := c.Param("snapshotId")
	userID, _ := c.Get("userID")
	userIDStr := userID.(string)

	var req struct {
		TargetSessionID string `json:"targetSessionId"` // Optional: create new session or restore to existing
	}

	c.ShouldBindJSON(&req)

	ctx := context.Background()

	// Verify session ownership
	if !h.verifySessionOwnership(ctx, sessionID, userIDStr) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this session"})
		return
	}

	// Get snapshot details
	var status, storagePath string
	err := h.db.DB().QueryRowContext(ctx, `
		SELECT status, storage_path FROM session_snapshots WHERE id = $1
	`, snapshotID).Scan(&status, &storagePath)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Snapshot not found"})
		return
	}

	if status != "available" {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Snapshot not available (status: %s)", status)})
		return
	}

	// Determine target session
	targetSession := sessionID
	if req.TargetSessionID != "" {
		targetSession = req.TargetSessionID
	}

	// Create restore job
	restoreID := fmt.Sprintf("restore_%d", time.Now().UnixNano())

	_, err = h.db.DB().ExecContext(ctx, `
		INSERT INTO snapshot_restore_jobs (id, snapshot_id, session_id, target_session_id, user_id, status)
		VALUES ($1, $2, $3, $4, $5, 'pending')
	`, restoreID, snapshotID, sessionID, targetSession, userIDStr)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create restore job"})
		return
	}

	// Trigger async restore
	go h.restoreSnapshotAsync(restoreID, snapshotID, sessionID, targetSession, storagePath)

	c.JSON(http.StatusAccepted, gin.H{
		"message":         "Restore job initiated",
		"restoreJobId":    restoreID,
		"snapshotId":      snapshotID,
		"targetSessionId": targetSession,
		"status":          "pending",
	})
}

// GetRestoreStatus returns the status of a restore operation
func (h *SnapshotsHandler) GetRestoreStatus(c *gin.Context) {
	snapshotID := c.Param("snapshotId")

	ctx := context.Background()

	rows, err := h.db.DB().QueryContext(ctx, `
		SELECT id, status, started_at, completed_at, error_message
		FROM snapshot_restore_jobs
		WHERE snapshot_id = $1
		ORDER BY started_at DESC
		LIMIT 10
	`, snapshotID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get restore status"})
		return
	}
	defer rows.Close()

	jobs := []map[string]interface{}{}
	for rows.Next() {
		var id, status string
		var errorMessage sql.NullString
		var startedAt time.Time
		var completedAt sql.NullTime

		if err := rows.Scan(&id, &status, &startedAt, &completedAt, &errorMessage); err == nil {
			job := map[string]interface{}{
				"id":        id,
				"status":    status,
				"startedAt": startedAt,
			}
			if completedAt.Valid {
				job["completedAt"] = completedAt.Time
			}
			if errorMessage.Valid {
				job["errorMessage"] = errorMessage.String
			}
			jobs = append(jobs, job)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"restoreJobs": jobs,
		"count":       len(jobs),
	})
}

// ListAllUserSnapshots lists all snapshots for the authenticated user
func (h *SnapshotsHandler) ListAllUserSnapshots(c *gin.Context) {
	userID, _ := c.Get("userID")
	userIDStr := userID.(string)

	ctx := context.Background()

	rows, err := h.db.DB().QueryContext(ctx, `
		SELECT id, session_id, user_id, name, description, type, status, size_bytes, created_at, expires_at
		FROM session_snapshots
		WHERE user_id = $1 AND status != 'deleted'
		ORDER BY created_at DESC
		LIMIT 100
	`, userIDStr)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list snapshots"})
		return
	}
	defer rows.Close()

	snapshots := []map[string]interface{}{}
	for rows.Next() {
		var id, sessionID, userID, name, snapshotType, status string
		var description sql.NullString
		var sizeBytes int64
		var createdAt time.Time
		var expiresAt sql.NullTime

		if err := rows.Scan(&id, &sessionID, &userID, &name, &description, &snapshotType, &status, &sizeBytes, &createdAt, &expiresAt); err == nil {
			snapshot := map[string]interface{}{
				"id":        id,
				"sessionId": sessionID,
				"name":      name,
				"type":      snapshotType,
				"status":    status,
				"sizeBytes": sizeBytes,
				"createdAt": createdAt,
			}
			if description.Valid {
				snapshot["description"] = description.String
			}
			if expiresAt.Valid {
				snapshot["expiresAt"] = expiresAt.Time
			}
			snapshots = append(snapshots, snapshot)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"snapshots": snapshots,
		"count":     len(snapshots),
	})
}

// GetSnapshotStats returns snapshot statistics
func (h *SnapshotsHandler) GetSnapshotStats(c *gin.Context) {
	userID, _ := c.Get("userID")
	userIDStr := userID.(string)

	ctx := context.Background()

	var totalCount, availableCount, totalSize int64

	h.db.DB().QueryRowContext(ctx, `
		SELECT
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE status = 'available') as available,
			COALESCE(SUM(size_bytes), 0) as total_size
		FROM session_snapshots
		WHERE user_id = $1 AND status != 'deleted'
	`, userIDStr).Scan(&totalCount, &availableCount, &totalSize)

	c.JSON(http.StatusOK, gin.H{
		"totalSnapshots":     totalCount,
		"availableSnapshots": availableCount,
		"totalSizeBytes":     totalSize,
		"totalSizeGB":        float64(totalSize) / (1024 * 1024 * 1024),
	})
}

// GetSnapshotConfig returns snapshot configuration for a session
func (h *SnapshotsHandler) GetSnapshotConfig(c *gin.Context) {
	sessionID := c.Param("sessionId")

	ctx := context.Background()

	var configJSON []byte
	err := h.db.DB().QueryRowContext(ctx, `
		SELECT snapshot_config FROM sessions WHERE id = $1
	`, sessionID).Scan(&configJSON)

	if err == sql.ErrNoRows {
		// Return default config
		c.JSON(http.StatusOK, h.getDefaultSnapshotConfig())
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get snapshot config"})
		return
	}

	var config map[string]interface{}
	json.Unmarshal(configJSON, &config)

	c.JSON(http.StatusOK, config)
}

// UpdateSnapshotConfig updates snapshot configuration
func (h *SnapshotsHandler) UpdateSnapshotConfig(c *gin.Context) {
	sessionID := c.Param("sessionId")

	var config map[string]interface{}
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()

	configJSON, _ := json.Marshal(config)

	_, err := h.db.DB().ExecContext(ctx, `
		UPDATE sessions SET snapshot_config = $1 WHERE id = $2
	`, configJSON, sessionID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update snapshot config"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Snapshot configuration updated",
		"config":  config,
	})
}

// Helper functions

func (h *SnapshotsHandler) verifySessionOwnership(ctx context.Context, sessionID, userID string) bool {
	var ownerID string
	err := h.db.DB().QueryRowContext(ctx, `SELECT user_id FROM sessions WHERE id = $1`, sessionID).Scan(&ownerID)
	return err == nil && ownerID == userID
}

func (h *SnapshotsHandler) getSnapshotStoragePath(sessionID, snapshotID string) string {
	baseDir := os.Getenv("SNAPSHOT_STORAGE_PATH")
	if baseDir == "" {
		baseDir = "/data/snapshots"
	}
	return filepath.Join(baseDir, sessionID, snapshotID)
}

func (h *SnapshotsHandler) createSnapshotAsync(snapshotID, sessionID, storagePath string) {
	ctx := context.Background()

	// Update status to creating
	h.db.DB().ExecContext(ctx, `
		UPDATE session_snapshots SET status = 'creating', updated_at = CURRENT_TIMESTAMP WHERE id = $1
	`, snapshotID)

	// Real snapshot creation implementation
	sizeBytes, err := h.performSnapshotCreation(ctx, sessionID, snapshotID, storagePath)

	if err != nil {
		// Mark as failed
		log.Printf("[ERROR] Snapshot creation failed for %s: %v", snapshotID, err)
		h.db.DB().ExecContext(ctx, `
			UPDATE session_snapshots
			SET status = 'failed', error_message = $1, updated_at = CURRENT_TIMESTAMP
			WHERE id = $2
		`, err.Error(), snapshotID)
		return
	}

	// Mark as available with real size
	_, err = h.db.DB().ExecContext(ctx, `
		UPDATE session_snapshots
		SET status = 'available', size_bytes = $1, completed_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`, sizeBytes, snapshotID)

	if err != nil {
		log.Printf("[ERROR] Failed to update snapshot status: %v", err)
	} else {
		log.Printf("[INFO] Snapshot %s created successfully (%d bytes)", snapshotID, sizeBytes)
	}
}

func (h *SnapshotsHandler) performSnapshotCreation(ctx context.Context, sessionID, snapshotID, storagePath string) (int64, error) {
	// Get Kubernetes client
	k8sClient, err := k8s.NewClient()
	if err != nil {
		return 0, fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	// Get namespace from environment
	namespace := os.Getenv("NAMESPACE")
	if namespace == "" {
		namespace = "streamspace"
	}

	// Get session to find pod name
	session, err := k8sClient.GetSession(ctx, namespace, sessionID)
	if err != nil {
		return 0, fmt.Errorf("failed to get session: %w", err)
	}

	if session.Status.PodName == "" {
		return 0, fmt.Errorf("session pod not found or not running")
	}

	podName := session.Status.PodName
	log.Printf("[INFO] Creating snapshot for session %s (pod: %s)", sessionID, podName)

	// Create snapshot storage directory
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		return 0, fmt.Errorf("failed to create snapshot directory: %w", err)
	}

	// Define tar file path
	tarFileName := fmt.Sprintf("snapshot-%s.tar.gz", snapshotID)
	tarFilePath := filepath.Join(storagePath, tarFileName)

	// Execute tar command inside the pod to create compressed archive
	// Using kubectl exec to tar /config directory (where session data is stored)
	log.Printf("[INFO] Executing tar command in pod %s", podName)

	kubectlPath := os.Getenv("KUBECTL_PATH")
	if kubectlPath == "" {
		kubectlPath = "kubectl"
	}

	// Create tar.gz archive of /config directory
	tarCmd := exec.CommandContext(ctx,
		kubectlPath, "exec", "-n", namespace, podName, "--",
		"tar", "-czf", "-", "-C", "/config", ".",
	)

	// Create output file
	outFile, err := os.Create(tarFilePath)
	if err != nil {
		return 0, fmt.Errorf("failed to create tar file: %w", err)
	}
	defer outFile.Close()

	// Capture stderr for error messages
	var stderr bytes.Buffer
	tarCmd.Stdout = outFile
	tarCmd.Stderr = &stderr

	// Execute tar command
	if err := tarCmd.Run(); err != nil {
		return 0, fmt.Errorf("failed to execute tar command: %w, stderr: %s", err, stderr.String())
	}

	// Get file size
	fileInfo, err := os.Stat(tarFilePath)
	if err != nil {
		return 0, fmt.Errorf("failed to get snapshot file size: %w", err)
	}

	sizeBytes := fileInfo.Size()
	log.Printf("[INFO] Snapshot created: %s (%d bytes, %.2f MB)", tarFilePath, sizeBytes, float64(sizeBytes)/(1024*1024))

	// Create metadata file
	metadataPath := filepath.Join(storagePath, "metadata.json")
	metadata := map[string]interface{}{
		"snapshot_id": snapshotID,
		"session_id":  sessionID,
		"pod_name":    podName,
		"created_at":  time.Now().UTC().Format(time.RFC3339),
		"size_bytes":  sizeBytes,
		"tar_file":    tarFileName,
		"compression": "gzip",
		"source_path": "/config",
	}

	metadataJSON, err := json.MarshalIndent(metadata, "", "  ")
	if err == nil {
		os.WriteFile(metadataPath, metadataJSON, 0644)
	}

	return sizeBytes, nil
}

func (h *SnapshotsHandler) restoreSnapshotAsync(restoreID, snapshotID, sessionID, targetSession, storagePath string) {
	ctx := context.Background()

	// Update restore job status
	h.db.DB().ExecContext(ctx, `
		UPDATE snapshot_restore_jobs SET status = 'in_progress', started_at = CURRENT_TIMESTAMP WHERE id = $1
	`, restoreID)

	// Real snapshot restore implementation
	err := h.performSnapshotRestore(ctx, snapshotID, targetSession, storagePath)

	if err != nil {
		log.Printf("[ERROR] Snapshot restore failed for %s: %v", restoreID, err)
		h.db.DB().ExecContext(ctx, `
			UPDATE snapshot_restore_jobs
			SET status = 'failed', error_message = $1, completed_at = CURRENT_TIMESTAMP
			WHERE id = $2
		`, err.Error(), restoreID)
		return
	}

	// Mark as completed
	_, err = h.db.DB().ExecContext(ctx, `
		UPDATE snapshot_restore_jobs
		SET status = 'completed', completed_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`, restoreID)

	if err != nil {
		log.Printf("[ERROR] Failed to update restore job status: %v", err)
	} else {
		log.Printf("[INFO] Snapshot %s restored successfully to session %s", snapshotID, targetSession)
	}
}

func (h *SnapshotsHandler) performSnapshotRestore(ctx context.Context, snapshotID, targetSession, storagePath string) error {
	// Get Kubernetes client
	k8sClient, err := k8s.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	// Get namespace from environment
	namespace := os.Getenv("NAMESPACE")
	if namespace == "" {
		namespace = "streamspace"
	}

	// Get target session to find pod name
	session, err := k8sClient.GetSession(ctx, namespace, targetSession)
	if err != nil {
		return fmt.Errorf("failed to get target session: %w", err)
	}

	if session.Status.PodName == "" {
		return fmt.Errorf("target session pod not found or not running")
	}

	podName := session.Status.PodName
	log.Printf("[INFO] Restoring snapshot %s to session %s (pod: %s)", snapshotID, targetSession, podName)

	// Find the snapshot tar file
	tarFileName := fmt.Sprintf("snapshot-%s.tar.gz", snapshotID)
	tarFilePath := filepath.Join(storagePath, tarFileName)

	// Verify snapshot file exists
	if _, err := os.Stat(tarFilePath); err != nil {
		return fmt.Errorf("snapshot file not found: %w", err)
	}

	log.Printf("[INFO] Found snapshot file: %s", tarFilePath)

	kubectlPath := os.Getenv("KUBECTL_PATH")
	if kubectlPath == "" {
		kubectlPath = "kubectl"
	}

	// Step 1: Backup existing data (optional but recommended)
	log.Printf("[INFO] Creating backup of existing data before restore")
	backupCmd := exec.CommandContext(ctx,
		kubectlPath, "exec", "-n", namespace, podName, "--",
		"sh", "-c", "cd /config && tar -czf /tmp/pre-restore-backup.tar.gz . 2>/dev/null || true",
	)
	var backupStderr bytes.Buffer
	backupCmd.Stderr = &backupStderr
	if err := backupCmd.Run(); err != nil {
		log.Printf("[WARN] Failed to create pre-restore backup: %v, stderr: %s", err, backupStderr.String())
		// Continue anyway - backup is optional
	}

	// Step 2: Clear existing /config directory
	log.Printf("[INFO] Clearing existing data in /config")
	clearCmd := exec.CommandContext(ctx,
		kubectlPath, "exec", "-n", namespace, podName, "--",
		"sh", "-c", "rm -rf /config/* /config/.[!.]* 2>/dev/null || true",
	)
	var clearStderr bytes.Buffer
	clearCmd.Stderr = &clearStderr
	if err := clearCmd.Run(); err != nil {
		log.Printf("[WARN] Error clearing directory: %v, stderr: %s", err, clearStderr.String())
		// Continue anyway - some files may be locked
	}

	// Step 3: Extract snapshot tar file into pod
	log.Printf("[INFO] Extracting snapshot into pod")

	// Open tar file for reading
	tarFile, err := os.Open(tarFilePath)
	if err != nil {
		return fmt.Errorf("failed to open tar file: %w", err)
	}
	defer tarFile.Close()

	// Extract tar.gz to /config directory
	extractCmd := exec.CommandContext(ctx,
		kubectlPath, "exec", "-i", "-n", namespace, podName, "--",
		"tar", "-xzf", "-", "-C", "/config",
	)
	extractCmd.Stdin = tarFile

	var extractStderr bytes.Buffer
	extractCmd.Stderr = &extractStderr

	if err := extractCmd.Run(); err != nil {
		return fmt.Errorf("failed to extract snapshot: %w, stderr: %s", err, extractStderr.String())
	}

	log.Printf("[INFO] Snapshot extracted successfully")

	// Step 4: Verify restoration by checking file count
	verifyCmd := exec.CommandContext(ctx,
		kubectlPath, "exec", "-n", namespace, podName, "--",
		"sh", "-c", "find /config -type f | wc -l",
	)

	var verifyOutput bytes.Buffer
	verifyCmd.Stdout = &verifyOutput
	if err := verifyCmd.Run(); err == nil {
		log.Printf("[INFO] Restored file count: %s", strings.TrimSpace(verifyOutput.String()))
	}

	// Step 5: Fix permissions
	log.Printf("[INFO] Fixing permissions on restored files")
	chownCmd := exec.CommandContext(ctx,
		kubectlPath, "exec", "-n", namespace, podName, "--",
		"sh", "-c", "chown -R 1000:1000 /config 2>/dev/null || true",
	)
	var chownStderr bytes.Buffer
	chownCmd.Stderr = &chownStderr
	if err := chownCmd.Run(); err != nil {
		log.Printf("[WARN] Failed to fix permissions: %v, stderr: %s", err, chownStderr.String())
		// Continue - permissions may already be correct
	}

	log.Printf("[INFO] Snapshot restore completed successfully")
	return nil
}

func (h *SnapshotsHandler) deleteSnapshotFiles(storagePath string) {
	// Check if storage path is empty or invalid
	if storagePath == "" {
		log.Printf("Warning: Cannot delete snapshot files - empty storage path")
		return
	}

	// Security check: Ensure path is within snapshot storage directory
	baseDir := os.Getenv("SNAPSHOT_STORAGE_PATH")
	if baseDir == "" {
		baseDir = "/data/snapshots"
	}

	// Resolve absolute paths to prevent directory traversal
	absStoragePath, err := filepath.Abs(storagePath)
	if err != nil {
		log.Printf("Error resolving snapshot path %s: %v", storagePath, err)
		return
	}

	absBaseDir, err := filepath.Abs(baseDir)
	if err != nil {
		log.Printf("Error resolving base directory %s: %v", baseDir, err)
		return
	}

	// Ensure the storage path is within the base directory
	if !isSubPath(absStoragePath, absBaseDir) {
		log.Printf("Security violation: Attempt to delete files outside snapshot storage: %s", absStoragePath)
		return
	}

	// Check if path exists
	if _, err := os.Stat(absStoragePath); os.IsNotExist(err) {
		// Path doesn't exist, nothing to delete (already cleaned up or never created)
		log.Printf("Snapshot path does not exist (already deleted): %s", absStoragePath)
		return
	}

	// Delete the snapshot directory and all its contents
	err = os.RemoveAll(absStoragePath)
	if err != nil {
		log.Printf("Error deleting snapshot files at %s: %v", absStoragePath, err)
		return
	}

	log.Printf("Successfully deleted snapshot files at %s", absStoragePath)
}

// isSubPath checks if the child path is within the parent path
func isSubPath(child, parent string) bool {
	// Clean and resolve paths
	cleanChild := filepath.Clean(child)
	cleanParent := filepath.Clean(parent)

	// Check if child starts with parent
	rel, err := filepath.Rel(cleanParent, cleanChild)
	if err != nil {
		return false
	}

	// If relative path starts with "..", child is outside parent
	return !filepath.IsAbs(rel) && !containsDotDot(rel)
}

// containsDotDot checks if a path contains ".." components
func containsDotDot(path string) bool {
	// Split path by separator and check each component
	parts := strings.Split(path, string(filepath.Separator))
	for _, part := range parts {
		if part == ".." {
			return true
		}
	}
	return false
}

func (h *SnapshotsHandler) getDefaultSnapshotConfig() map[string]interface{} {
	return map[string]interface{}{
		"automaticSnapshots": map[string]interface{}{
			"enabled":  false,
			"schedule": "0 2 * * *", // Daily at 2 AM
		},
		"retention": map[string]interface{}{
			"maxSnapshots":      10,
			"retentionDays":     30,
			"deleteExpiredAuto": true,
		},
		"compression": map[string]interface{}{
			"enabled": true,
			"level":   6,
		},
	}
}
