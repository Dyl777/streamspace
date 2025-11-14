package api

import (
	"bufio"
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	templateGVR = schema.GroupVersionResource{
		Group:    "stream.streamspace.io",
		Version:  "v1alpha1",
		Resource: "templates",
	}
)

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Get allowed origins from environment variable
		allowedOrigins := os.Getenv("ALLOWED_ORIGINS")

		// If not set, allow localhost for development only
		if allowedOrigins == "" {
			allowedOrigins = "http://localhost:3000,http://localhost:5173"
		}

		// Special case: "*" means allow all (use with caution)
		if allowedOrigins == "*" {
			log.Println("WARNING: WebSocket accepting connections from all origins")
			return true
		}

		// Check if request origin is in allowed list
		origin := r.Header.Get("Origin")
		for _, allowed := range strings.Split(allowedOrigins, ",") {
			if strings.TrimSpace(allowed) == origin {
				return true
			}
		}

		log.Printf("WebSocket connection rejected from origin: %s", origin)
		return false
	},
}

// ============================================================================
// Health & Version Endpoints
// ============================================================================

// Health returns health status
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"service": "streamspace-api",
	})
}

// Version returns API version
func (h *Handler) Version(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version": "v0.1.0",
		"api": "v1",
		"phase": "2.2",
	})
}

// ============================================================================
// Stub Methods (To Be Implemented)
// ============================================================================

// UpdateTemplate updates a template (admin only)
func (h *Handler) UpdateTemplate(c *gin.Context) {
	templateName := c.Param("id")
	if templateName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "template id required"})
		return
	}

	var updateReq struct {
		DisplayName      *string  `json:"displayName"`
		Description      *string  `json:"description"`
		Icon             *string  `json:"icon"`
		Tags             []string `json:"tags"`
		DefaultResources *struct {
			Memory string `json:"memory"`
			CPU    string `json:"cpu"`
		} `json:"defaultResources"`
	}

	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing template
	template, err := h.k8sClient.GetTemplate(c.Request.Context(), h.namespace, templateName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}

	// Apply updates
	if updateReq.DisplayName != nil {
		template.DisplayName = *updateReq.DisplayName
	}
	if updateReq.Description != nil {
		template.Description = *updateReq.Description
	}
	if updateReq.Icon != nil {
		template.Icon = *updateReq.Icon
	}
	if updateReq.Tags != nil {
		template.Tags = updateReq.Tags
	}
	if updateReq.DefaultResources != nil {
		template.DefaultResources.Memory = updateReq.DefaultResources.Memory
		template.DefaultResources.CPU = updateReq.DefaultResources.CPU
	}

	// Update template in Kubernetes using dynamic client
	obj := h.k8sClient.GetDynamicClient().Resource(templateGVR).Namespace(h.namespace)
	unstructuredTemplate, err := obj.Get(c.Request.Context(), templateName, metav1.GetOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update spec fields
	spec := unstructuredTemplate.Object["spec"].(map[string]interface{})
	spec["displayName"] = template.DisplayName
	spec["description"] = template.Description
	spec["icon"] = template.Icon
	spec["tags"] = template.Tags
	if updateReq.DefaultResources != nil {
		spec["defaultResources"] = map[string]interface{}{
			"memory": template.DefaultResources.Memory,
			"cpu":    template.DefaultResources.CPU,
		}
	}

	_, err = obj.Update(c.Request.Context(), unstructuredTemplate, metav1.UpdateOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Template updated successfully",
		"template": template,
	})
}

// ListNodes returns cluster nodes
// Note: This is now implemented in handlers/nodes.go via NodeHandler
// This stub remains for backwards compatibility with old routes
func (h *Handler) ListNodes(c *gin.Context) {
	nodes, err := h.k8sClient.GetNodes(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nodes)
}

// ListPods returns pods in namespace
func (h *Handler) ListPods(c *gin.Context) {
	namespace := c.Query("namespace")
	if namespace == "" {
		namespace = h.namespace
	}

	pods, err := h.k8sClient.GetPods(c.Request.Context(), namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pods)
}

// ListDeployments returns deployments
func (h *Handler) ListDeployments(c *gin.Context) {
	namespace := c.Query("namespace")
	if namespace == "" {
		namespace = h.namespace
	}

	deployments, err := h.k8sClient.GetClientset().AppsV1().Deployments(namespace).List(c.Request.Context(), metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, deployments)
}

// ListServices returns services
func (h *Handler) ListServices(c *gin.Context) {
	namespace := c.Query("namespace")
	if namespace == "" {
		namespace = h.namespace
	}

	services, err := h.k8sClient.GetServices(c.Request.Context(), namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, services)
}

// ListNamespaces returns namespaces
func (h *Handler) ListNamespaces(c *gin.Context) {
	namespaces, err := h.k8sClient.GetNamespaces(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, namespaces)
}

// CreateResource creates a K8s resource
func (h *Handler) CreateResource(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

// UpdateResource updates a K8s resource
func (h *Handler) UpdateResource(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

// DeleteResource deletes a K8s resource
func (h *Handler) DeleteResource(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

// GetPodLogs returns pod logs
func (h *Handler) GetPodLogs(c *gin.Context) {
	namespace := c.Query("namespace")
	if namespace == "" {
		namespace = h.namespace
	}
	podName := c.Query("pod")
	if podName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "pod query parameter required"})
		return
	}

	// Parse optional parameters
	tailLines := int64(100) // Default to last 100 lines
	follow := c.Query("follow") == "true"

	// Get pod logs
	opts := &corev1.PodLogOptions{
		TailLines: &tailLines,
		Follow:    follow,
	}

	req := h.k8sClient.GetClientset().CoreV1().Pods(namespace).GetLogs(podName, opts)
	stream, err := req.Stream(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer stream.Close()

	// If following logs, stream them
	if follow {
		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.Header("Transfer-Encoding", "chunked")
		c.Status(http.StatusOK)

		scanner := bufio.NewScanner(stream)
		for scanner.Scan() {
			c.Writer.Write([]byte(scanner.Text() + "\n"))
			c.Writer.Flush()
		}
		return
	}

	// Otherwise return all logs
	logs, err := io.ReadAll(stream)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, string(logs))
}

// GetConfig returns configuration
func (h *Handler) GetConfig(c *gin.Context) {
	// Get configuration from streamspace-config ConfigMap
	configMap, err := h.k8sClient.GetClientset().CoreV1().ConfigMaps(h.namespace).Get(
		c.Request.Context(),
		"streamspace-config",
		metav1.GetOptions{},
	)

	if err != nil {
		// Return default config if ConfigMap doesn't exist
		c.JSON(http.StatusOK, gin.H{
			"namespace":     h.namespace,
			"ingressDomain": os.Getenv("INGRESS_DOMAIN"),
			"hibernation": gin.H{
				"enabled":           true,
				"defaultIdleTimeout": "30m",
			},
			"resources": gin.H{
				"defaultMemory": "2Gi",
				"defaultCPU":    "1000m",
			},
		})
		return
	}

	c.JSON(http.StatusOK, configMap.Data)
}

// UpdateConfig updates configuration
func (h *Handler) UpdateConfig(c *gin.Context) {
	var config map[string]string
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get or create ConfigMap
	configMap, err := h.k8sClient.GetClientset().CoreV1().ConfigMaps(h.namespace).Get(
		c.Request.Context(),
		"streamspace-config",
		metav1.GetOptions{},
	)

	if err != nil {
		// Create new ConfigMap
		configMap = &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "streamspace-config",
				Namespace: h.namespace,
			},
			Data: config,
		}

		_, err = h.k8sClient.GetClientset().CoreV1().ConfigMaps(h.namespace).Create(
			c.Request.Context(),
			configMap,
			metav1.CreateOptions{},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		// Update existing ConfigMap
		configMap.Data = config
		_, err = h.k8sClient.GetClientset().CoreV1().ConfigMaps(h.namespace).Update(
			c.Request.Context(),
			configMap,
			metav1.UpdateOptions{},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Configuration updated successfully",
		"config":  config,
	})
}

// ListUsers returns all users
func (h *Handler) ListUsers(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

// CreateUser creates a new user
func (h *Handler) CreateUser(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

// GetCurrentUser returns current user info
func (h *Handler) GetCurrentUser(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

// GetUser returns user by ID
func (h *Handler) GetUser(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

// UpdateUser updates user
func (h *Handler) UpdateUser(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

// GetUserSessions returns sessions for a user
func (h *Handler) GetUserSessions(c *gin.Context) {
	userID := c.Param("id")
	c.Redirect(http.StatusTemporaryRedirect, "/api/v1/sessions?user="+userID)
}

// GetMetrics returns metrics
func (h *Handler) GetMetrics(c *gin.Context) {
	stats := h.connTracker.GetStats()
	c.JSON(http.StatusOK, stats)
}

// ============================================================================
// WebSocket Endpoints
// ============================================================================

// SessionsWebSocket handles WebSocket for real-time session updates
func (h *Handler) SessionsWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade WebSocket connection: %v", err)
		return
	}

	h.wsManager.HandleSessionsWebSocket(conn)
}

// ClusterWebSocket handles WebSocket for real-time cluster updates
func (h *Handler) ClusterWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade WebSocket connection: %v", err)
		return
	}

	h.wsManager.HandleMetricsWebSocket(conn)
}

// LogsWebSocket handles WebSocket for streaming pod logs
func (h *Handler) LogsWebSocket(c *gin.Context) {
	namespace := c.Param("namespace")
	podName := c.Param("pod")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade WebSocket connection: %v", err)
		return
	}

	h.wsManager.HandleLogsWebSocket(conn, namespace, podName)
}

// ============================================================================
// Catalog/Repository Endpoints (Additional)
// ============================================================================

// BrowseCatalog returns catalog templates (alias for ListCatalogTemplates)
func (h *Handler) BrowseCatalog(c *gin.Context) {
	h.ListCatalogTemplates(c)
}

// InstallTemplate installs a template from catalog (alias for InstallCatalogTemplate)
func (h *Handler) InstallTemplate(c *gin.Context) {
	catalogID := c.Query("id")
	if catalogID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id query parameter required"})
		return
	}

	c.Params = append(c.Params, gin.Param{Key: "id", Value: catalogID})
	h.InstallCatalogTemplate(c)
}

// SyncCatalog triggers sync for all repositories
func (h *Handler) SyncCatalog(c *gin.Context) {
	go func() {
		if err := h.syncService.SyncAllRepositories(c.Request.Context()); err != nil {
			log.Printf("Catalog sync failed: %v", err)
		}
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Catalog sync triggered",
		"status":  "syncing",
	})
}

// RemoveRepository removes a repository (alias for DeleteRepository)
func (h *Handler) RemoveRepository(c *gin.Context) {
	h.DeleteRepository(c)
}

// ============================================================================
// Webhook Endpoint for Repository Auto-Sync
// ============================================================================

// WebhookRepositorySync handles webhooks from Git providers for auto-sync
func (h *Handler) WebhookRepositorySync(c *gin.Context) {
	var webhook struct {
		RepositoryURL string `json:"repository_url"`
		Branch        string `json:"branch"`
		Ref           string `json:"ref"`
	}

	if err := c.ShouldBindJSON(&webhook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find repository by URL
	ctx := c.Request.Context()
	var repoID int
	err := h.db.DB().QueryRowContext(ctx, `
		SELECT id FROM repositories WHERE url = $1
	`, webhook.RepositoryURL).Scan(&repoID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository not found"})
		return
	}

	// Trigger sync in background
	go func() {
		if err := h.syncService.SyncRepository(ctx, repoID); err != nil {
			log.Printf("Webhook-triggered sync failed for repository %d: %v", repoID, err)
		} else {
			log.Printf("Webhook-triggered sync completed for repository %d", repoID)
		}
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"message":      "Webhook received, sync triggered",
		"repository":   webhook.RepositoryURL,
		"repositoryID": repoID,
	})
}
