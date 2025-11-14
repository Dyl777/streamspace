package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/streamspace/streamspace/api/internal/db"
	"github.com/streamspace/streamspace/api/internal/models"
)

// PluginHandler handles plugin-related HTTP requests
type PluginHandler struct {
	db *db.Database
}

// NewPluginHandler creates a new plugin handler
func NewPluginHandler(database *db.Database) *PluginHandler {
	return &PluginHandler{db: database}
}

// RegisterRoutes registers plugin routes
func (h *PluginHandler) RegisterRoutes(r *gin.RouterGroup) {
	plugins := r.Group("/plugins")
	{
		// Plugin catalog
		plugins.GET("/catalog", h.BrowsePluginCatalog)
		plugins.GET("/catalog/:id", h.GetCatalogPlugin)
		plugins.POST("/catalog/:id/rate", h.RatePlugin)
		plugins.POST("/catalog/:id/install", h.InstallPlugin)

		// Installed plugins
		plugins.GET("", h.ListInstalledPlugins)
		plugins.GET("/:id", h.GetInstalledPlugin)
		plugins.PATCH("/:id", h.UpdateInstalledPlugin)
		plugins.DELETE("/:id", h.UninstallPlugin)
		plugins.POST("/:id/enable", h.EnablePlugin)
		plugins.POST("/:id/disable", h.DisablePlugin)
	}
}

// BrowsePluginCatalog browses available plugins
func (h *PluginHandler) BrowsePluginCatalog(c *gin.Context) {
	category := c.Query("category")
	pluginType := c.Query("type")
	search := c.Query("search")
	sortBy := c.DefaultQuery("sort", "popular") // popular, rating, newest, name

	query := `
		SELECT
			cp.id, cp.repository_id, cp.name, cp.version, cp.display_name,
			cp.description, cp.category, cp.plugin_type, cp.icon_url,
			cp.manifest, cp.tags, cp.install_count, cp.avg_rating, cp.rating_count,
			cp.created_at, cp.updated_at,
			r.id as repo_id, r.name as repo_name, r.url as repo_url, r.type as repo_type
		FROM catalog_plugins cp
		JOIN catalog_repositories r ON cp.repository_id = r.id
		WHERE 1=1
	`

	args := []interface{}{}
	argIndex := 1

	if category != "" {
		query += ` AND cp.category = $` + strconv.Itoa(argIndex)
		args = append(args, category)
		argIndex++
	}

	if pluginType != "" {
		query += ` AND cp.plugin_type = $` + strconv.Itoa(argIndex)
		args = append(args, pluginType)
		argIndex++
	}

	if search != "" {
		query += ` AND (cp.display_name ILIKE $` + strconv.Itoa(argIndex) +
			` OR cp.description ILIKE $` + strconv.Itoa(argIndex) +
			` OR $` + strconv.Itoa(argIndex) + ` = ANY(cp.tags))`
		args = append(args, "%"+search+"%")
		argIndex++
	}

	// Sorting
	switch sortBy {
	case "popular":
		query += ` ORDER BY cp.install_count DESC, cp.avg_rating DESC`
	case "rating":
		query += ` ORDER BY cp.avg_rating DESC, cp.rating_count DESC`
	case "newest":
		query += ` ORDER BY cp.created_at DESC`
	case "name":
		query += ` ORDER BY cp.display_name ASC`
	default:
		query += ` ORDER BY cp.install_count DESC`
	}

	rows, err := h.db.DB().Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch plugins", "details": err.Error()})
		return
	}
	defer rows.Close()

	var plugins []models.CatalogPlugin
	for rows.Next() {
		var plugin models.CatalogPlugin
		var manifestJSON []byte
		var tags sql.NullString

		err := rows.Scan(
			&plugin.ID, &plugin.RepositoryID, &plugin.Name, &plugin.Version,
			&plugin.DisplayName, &plugin.Description, &plugin.Category, &plugin.PluginType,
			&plugin.IconURL, &manifestJSON, &tags, &plugin.InstallCount,
			&plugin.AvgRating, &plugin.RatingCount, &plugin.CreatedAt, &plugin.UpdatedAt,
			&plugin.Repository.ID, &plugin.Repository.Name, &plugin.Repository.URL, &plugin.Repository.Type,
		)
		if err != nil {
			continue
		}

		// Parse manifest
		if len(manifestJSON) > 0 {
			json.Unmarshal(manifestJSON, &plugin.Manifest)
		}

		// Parse tags
		if tags.Valid {
			// PostgreSQL array format: {tag1,tag2,tag3}
			tagsStr := tags.String
			if len(tagsStr) > 2 {
				tagsStr = tagsStr[1 : len(tagsStr)-1] // Remove { }
				json.Unmarshal([]byte(`["`+tagsStr+`"]`), &plugin.Tags)
			}
		}

		plugins = append(plugins, plugin)
	}

	c.JSON(http.StatusOK, gin.H{
		"plugins": plugins,
		"total":   len(plugins),
	})
}

// GetCatalogPlugin gets a specific plugin from the catalog
func (h *PluginHandler) GetCatalogPlugin(c *gin.Context) {
	id := c.Param("id")

	query := `
		SELECT
			cp.id, cp.repository_id, cp.name, cp.version, cp.display_name,
			cp.description, cp.category, cp.plugin_type, cp.icon_url,
			cp.manifest, cp.tags, cp.install_count, cp.avg_rating, cp.rating_count,
			cp.created_at, cp.updated_at,
			r.id as repo_id, r.name as repo_name, r.url as repo_url, r.type as repo_type
		FROM catalog_plugins cp
		JOIN catalog_repositories r ON cp.repository_id = r.id
		WHERE cp.id = $1
	`

	var plugin models.CatalogPlugin
	var manifestJSON []byte
	var tags sql.NullString

	err := h.db.DB().QueryRow(query, id).Scan(
		&plugin.ID, &plugin.RepositoryID, &plugin.Name, &plugin.Version,
		&plugin.DisplayName, &plugin.Description, &plugin.Category, &plugin.PluginType,
		&plugin.IconURL, &manifestJSON, &tags, &plugin.InstallCount,
		&plugin.AvgRating, &plugin.RatingCount, &plugin.CreatedAt, &plugin.UpdatedAt,
		&plugin.Repository.ID, &plugin.Repository.Name, &plugin.Repository.URL, &plugin.Repository.Type,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plugin not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch plugin", "details": err.Error()})
		return
	}

	// Parse manifest
	if len(manifestJSON) > 0 {
		json.Unmarshal(manifestJSON, &plugin.Manifest)
	}

	// Parse tags
	if tags.Valid {
		tagsStr := tags.String
		if len(tagsStr) > 2 {
			tagsStr = tagsStr[1 : len(tagsStr)-1]
			json.Unmarshal([]byte(`["`+tagsStr+`"]`), &plugin.Tags)
		}
	}

	// Get view count and update stats
	go func() {
		h.db.DB().Exec(`
			INSERT INTO plugin_stats (plugin_id, view_count, last_viewed_at)
			VALUES ($1, 1, $2)
			ON CONFLICT (plugin_id) DO UPDATE
			SET view_count = plugin_stats.view_count + 1,
			    last_viewed_at = $2,
			    updated_at = $2
		`, plugin.ID, time.Now())
	}()

	c.JSON(http.StatusOK, plugin)
}

// RatePlugin rates a plugin
func (h *PluginHandler) RatePlugin(c *gin.Context) {
	pluginID := c.Param("id")
	userID := c.GetString("user_id") // From auth middleware

	var req models.RatePluginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	if req.Rating < 1 || req.Rating > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Rating must be between 1 and 5"})
		return
	}

	// Insert or update rating
	_, err := h.db.DB().Exec(`
		INSERT INTO plugin_ratings (plugin_id, user_id, rating, review)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (plugin_id, user_id) DO UPDATE
		SET rating = $3, review = $4, updated_at = NOW()
	`, pluginID, userID, req.Rating, req.Review)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save rating", "details": err.Error()})
		return
	}

	// Update plugin average rating
	h.db.DB().Exec(`
		UPDATE catalog_plugins
		SET avg_rating = (SELECT AVG(rating) FROM plugin_ratings WHERE plugin_id = $1),
		    rating_count = (SELECT COUNT(*) FROM plugin_ratings WHERE plugin_id = $1),
		    updated_at = NOW()
		WHERE id = $1
	`, pluginID)

	c.JSON(http.StatusOK, gin.H{"message": "Rating submitted successfully"})
}

// InstallPlugin installs a plugin from the catalog
func (h *PluginHandler) InstallPlugin(c *gin.Context) {
	catalogPluginID := c.Param("id")
	userID := c.GetString("user_id")

	var req models.InstallPluginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		req.Config = json.RawMessage("{}")
	}

	// Get catalog plugin details
	var catalogPlugin models.CatalogPlugin
	var manifestJSON []byte
	err := h.db.DB().QueryRow(`
		SELECT id, name, version, display_name, description, plugin_type, icon_url, manifest
		FROM catalog_plugins
		WHERE id = $1
	`, catalogPluginID).Scan(
		&catalogPlugin.ID, &catalogPlugin.Name, &catalogPlugin.Version,
		&catalogPlugin.DisplayName, &catalogPlugin.Description,
		&catalogPlugin.PluginType, &catalogPlugin.IconURL, &manifestJSON,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plugin not found in catalog"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch plugin", "details": err.Error()})
		return
	}

	// Parse manifest
	if len(manifestJSON) > 0 {
		json.Unmarshal(manifestJSON, &catalogPlugin.Manifest)
	}

	// Check if already installed
	var existingID int
	err = h.db.DB().QueryRow(`
		SELECT id FROM installed_plugins WHERE name = $1
	`, catalogPlugin.Name).Scan(&existingID)

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Plugin already installed", "pluginId": existingID})
		return
	}

	// Install plugin
	var installedID int
	err = h.db.DB().QueryRow(`
		INSERT INTO installed_plugins (catalog_plugin_id, name, version, enabled, config, installed_by)
		VALUES ($1, $2, $3, true, $4, $5)
		RETURNING id
	`, catalogPlugin.ID, catalogPlugin.Name, catalogPlugin.Version, req.Config, userID).Scan(&installedID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to install plugin", "details": err.Error()})
		return
	}

	// Update install count
	go func() {
		h.db.DB().Exec(`
			UPDATE catalog_plugins
			SET install_count = install_count + 1
			WHERE id = $1
		`, catalogPlugin.ID)

		h.db.DB().Exec(`
			INSERT INTO plugin_stats (plugin_id, install_count, last_installed_at)
			VALUES ($1, 1, $2)
			ON CONFLICT (plugin_id) DO UPDATE
			SET install_count = plugin_stats.install_count + 1,
			    last_installed_at = $2,
			    updated_at = $2
		`, catalogPlugin.ID, time.Now())
	}()

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Plugin installed successfully",
		"pluginId": installedID,
	})
}

// ListInstalledPlugins lists all installed plugins
func (h *PluginHandler) ListInstalledPlugins(c *gin.Context) {
	enabledOnly := c.Query("enabled") == "true"

	query := `
		SELECT
			ip.id, ip.catalog_plugin_id, ip.name, ip.version, ip.enabled,
			ip.config, ip.installed_by, ip.installed_at, ip.updated_at,
			cp.display_name, cp.description, cp.plugin_type, cp.icon_url, cp.manifest
		FROM installed_plugins ip
		LEFT JOIN catalog_plugins cp ON ip.catalog_plugin_id = cp.id
	`

	if enabledOnly {
		query += ` WHERE ip.enabled = true`
	}

	query += ` ORDER BY ip.installed_at DESC`

	rows, err := h.db.DB().Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch plugins", "details": err.Error()})
		return
	}
	defer rows.Close()

	var plugins []models.InstalledPlugin
	for rows.Next() {
		var plugin models.InstalledPlugin
		var catalogPluginID sql.NullInt64
		var displayName, description, pluginType, iconURL sql.NullString
		var manifestJSON []byte

		err := rows.Scan(
			&plugin.ID, &catalogPluginID, &plugin.Name, &plugin.Version, &plugin.Enabled,
			&plugin.Config, &plugin.InstalledBy, &plugin.InstalledAt, &plugin.UpdatedAt,
			&displayName, &description, &pluginType, &iconURL, &manifestJSON,
		)
		if err != nil {
			continue
		}

		if catalogPluginID.Valid {
			id := int(catalogPluginID.Int64)
			plugin.CatalogPluginID = &id
		}

		if displayName.Valid {
			plugin.DisplayName = displayName.String
		}
		if description.Valid {
			plugin.Description = description.String
		}
		if pluginType.Valid {
			plugin.PluginType = pluginType.String
		}
		if iconURL.Valid {
			plugin.IconURL = iconURL.String
		}

		if len(manifestJSON) > 0 {
			var manifest models.PluginManifest
			if json.Unmarshal(manifestJSON, &manifest) == nil {
				plugin.Manifest = &manifest
			}
		}

		plugins = append(plugins, plugin)
	}

	c.JSON(http.StatusOK, gin.H{
		"plugins": plugins,
		"total":   len(plugins),
	})
}

// GetInstalledPlugin gets a specific installed plugin
func (h *PluginHandler) GetInstalledPlugin(c *gin.Context) {
	id := c.Param("id")

	query := `
		SELECT
			ip.id, ip.catalog_plugin_id, ip.name, ip.version, ip.enabled,
			ip.config, ip.installed_by, ip.installed_at, ip.updated_at,
			cp.display_name, cp.description, cp.plugin_type, cp.icon_url, cp.manifest
		FROM installed_plugins ip
		LEFT JOIN catalog_plugins cp ON ip.catalog_plugin_id = cp.id
		WHERE ip.id = $1
	`

	var plugin models.InstalledPlugin
	var catalogPluginID sql.NullInt64
	var displayName, description, pluginType, iconURL sql.NullString
	var manifestJSON []byte

	err := h.db.DB().QueryRow(query, id).Scan(
		&plugin.ID, &catalogPluginID, &plugin.Name, &plugin.Version, &plugin.Enabled,
		&plugin.Config, &plugin.InstalledBy, &plugin.InstalledAt, &plugin.UpdatedAt,
		&displayName, &description, &pluginType, &iconURL, &manifestJSON,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plugin not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch plugin", "details": err.Error()})
		return
	}

	if catalogPluginID.Valid {
		id := int(catalogPluginID.Int64)
		plugin.CatalogPluginID = &id
	}

	if displayName.Valid {
		plugin.DisplayName = displayName.String
	}
	if description.Valid {
		plugin.Description = description.String
	}
	if pluginType.Valid {
		plugin.PluginType = pluginType.String
	}
	if iconURL.Valid {
		plugin.IconURL = iconURL.String
	}

	if len(manifestJSON) > 0 {
		var manifest models.PluginManifest
		if json.Unmarshal(manifestJSON, &manifest) == nil {
			plugin.Manifest = &manifest
		}
	}

	c.JSON(http.StatusOK, plugin)
}

// UpdateInstalledPlugin updates a plugin's configuration
func (h *PluginHandler) UpdateInstalledPlugin(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdatePluginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	query := `UPDATE installed_plugins SET `
	args := []interface{}{}
	argIndex := 1

	if req.Enabled != nil {
		query += `enabled = $` + strconv.Itoa(argIndex) + `, `
		args = append(args, *req.Enabled)
		argIndex++
	}

	if req.Config != nil {
		query += `config = $` + strconv.Itoa(argIndex) + `, `
		args = append(args, req.Config)
		argIndex++
	}

	query += `updated_at = NOW() WHERE id = $` + strconv.Itoa(argIndex)
	args = append(args, id)

	result, err := h.db.DB().Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update plugin", "details": err.Error()})
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plugin not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Plugin updated successfully"})
}

// UninstallPlugin uninstalls a plugin
func (h *PluginHandler) UninstallPlugin(c *gin.Context) {
	id := c.Param("id")

	result, err := h.db.DB().Exec(`DELETE FROM installed_plugins WHERE id = $1`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to uninstall plugin", "details": err.Error()})
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plugin not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Plugin uninstalled successfully"})
}

// EnablePlugin enables a plugin
func (h *PluginHandler) EnablePlugin(c *gin.Context) {
	id := c.Param("id")

	result, err := h.db.DB().Exec(`
		UPDATE installed_plugins
		SET enabled = true, updated_at = NOW()
		WHERE id = $1
	`, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enable plugin", "details": err.Error()})
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plugin not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Plugin enabled successfully"})
}

// DisablePlugin disables a plugin
func (h *PluginHandler) DisablePlugin(c *gin.Context) {
	id := c.Param("id")

	result, err := h.db.DB().Exec(`
		UPDATE installed_plugins
		SET enabled = false, updated_at = NOW()
		WHERE id = $1
	`, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to disable plugin", "details": err.Error()})
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plugin not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Plugin disabled successfully"})
}
