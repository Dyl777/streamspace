package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// TimeoutConfig holds configuration for request timeouts
type TimeoutConfig struct {
	// Timeout is the maximum duration for the entire request
	Timeout time.Duration

	// ErrorMessage is the message returned when timeout occurs
	ErrorMessage string

	// ExcludedPaths are paths that should not have timeout applied
	// (e.g., WebSocket endpoints, file uploads)
	ExcludedPaths []string
}

// DefaultTimeoutConfig returns default timeout configuration
func DefaultTimeoutConfig() TimeoutConfig {
	return TimeoutConfig{
		Timeout:      30 * time.Second,
		ErrorMessage: "Request timeout",
		ExcludedPaths: []string{
			"/api/v1/ws/",    // WebSocket endpoints
			"/api/v1/upload", // File uploads
		},
	}
}

// Timeout middleware enforces a timeout on requests to prevent slow loris attacks
// and ensure resources are freed in a timely manner
func Timeout(config TimeoutConfig) gin.HandlerFunc {
	// Build exclusion map for fast lookup
	excluded := make(map[string]bool)
	for _, path := range config.ExcludedPaths {
		excluded[path] = true
	}

	return func(c *gin.Context) {
		// Check if path should be excluded
		path := c.Request.URL.Path
		for excludedPath := range excluded {
			if len(path) >= len(excludedPath) && path[:len(excludedPath)] == excludedPath {
				c.Next()
				return
			}
		}

		// Create context with timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), config.Timeout)
		defer cancel()

		// Replace request context
		c.Request = c.Request.WithContext(ctx)

		// Channel to signal completion
		finished := make(chan struct{})

		go func() {
			c.Next()
			close(finished)
		}()

		select {
		case <-finished:
			// Request completed successfully
			return
		case <-ctx.Done():
			// Timeout occurred
			c.AbortWithStatusJSON(http.StatusRequestTimeout, gin.H{
				"error":   config.ErrorMessage,
				"message": "The request took too long to process",
				"timeout": config.Timeout.String(),
			})
			return
		}
	}
}

// TimeoutWithDuration creates a timeout middleware with specified duration
func TimeoutWithDuration(timeout time.Duration) gin.HandlerFunc {
	config := DefaultTimeoutConfig()
	config.Timeout = timeout
	return Timeout(config)
}
