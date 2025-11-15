package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// StructuredLogger provides structured logging for all requests
// Logs include request ID, method, path, status, duration, and client IP
func StructuredLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate request duration
		duration := time.Since(start)

		// Get request ID (if RequestID middleware is used)
		requestID := GetRequestID(c)

		// Get status code
		status := c.Writer.Status()

		// Build log entry with structured fields
		logEntry := map[string]interface{}{
			"request_id": requestID,
			"method":     c.Request.Method,
			"path":       path,
			"query":      raw,
			"status":     status,
			"duration":   duration.String(),
			"duration_ms": duration.Milliseconds(),
			"client_ip":  c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		}

		// Add user info if authenticated
		if userID, exists := c.Get("userID"); exists {
			logEntry["user_id"] = userID
		}
		if username, exists := c.Get("username"); exists {
			logEntry["username"] = username
		}

		// Add error if present
		if len(c.Errors) > 0 {
			logEntry["errors"] = c.Errors.String()
		}

		// Determine log level based on status code
		if status >= 500 {
			log.Printf("ERROR %v", logEntry)
		} else if status >= 400 {
			log.Printf("WARN %v", logEntry)
		} else {
			log.Printf("INFO %v", logEntry)
		}
	}
}

// StructuredLoggerWithConfig allows customization of structured logging
type StructuredLoggerConfig struct {
	// SkipPaths is a list of paths to skip logging (e.g., health checks)
	SkipPaths []string

	// SkipHealthCheck if true, skips logging for /health endpoint
	SkipHealthCheck bool

	// LogQuery if false, skips logging query parameters (for privacy)
	LogQuery bool

	// LogUserAgent if false, skips logging user agent
	LogUserAgent bool
}

// DefaultStructuredLoggerConfig returns default configuration
func DefaultStructuredLoggerConfig() StructuredLoggerConfig {
	return StructuredLoggerConfig{
		SkipPaths:       []string{},
		SkipHealthCheck: true,
		LogQuery:        true,
		LogUserAgent:    true,
	}
}

// StructuredLoggerWithConfigFunc creates a structured logger with custom config
func StructuredLoggerWithConfigFunc(config StructuredLoggerConfig) gin.HandlerFunc {
	// Build skip map for fast lookup
	skipMap := make(map[string]bool)
	for _, path := range config.SkipPaths {
		skipMap[path] = true
	}
	if config.SkipHealthCheck {
		skipMap["/health"] = true
		skipMap["/api/v1/health"] = true
	}

	return func(c *gin.Context) {
		// Skip logging for certain paths
		path := c.Request.URL.Path
		if skipMap[path] {
			c.Next()
			return
		}

		// Start timer
		start := time.Now()
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate request duration
		duration := time.Since(start)

		// Get request ID
		requestID := GetRequestID(c)

		// Get status code
		status := c.Writer.Status()

		// Build log entry
		logEntry := map[string]interface{}{
			"request_id": requestID,
			"method":     c.Request.Method,
			"path":       path,
			"status":     status,
			"duration":   duration.String(),
			"duration_ms": duration.Milliseconds(),
			"client_ip":  c.ClientIP(),
		}

		// Conditionally add query
		if config.LogQuery && raw != "" {
			logEntry["query"] = raw
		}

		// Conditionally add user agent
		if config.LogUserAgent {
			logEntry["user_agent"] = c.Request.UserAgent()
		}

		// Add user info if authenticated
		if userID, exists := c.Get("userID"); exists {
			logEntry["user_id"] = userID
		}
		if username, exists := c.Get("username"); exists {
			logEntry["username"] = username
		}

		// Add error if present
		if len(c.Errors) > 0 {
			logEntry["errors"] = c.Errors.String()
		}

		// Log based on status code
		if status >= 500 {
			log.Printf("ERROR %v", logEntry)
		} else if status >= 400 {
			log.Printf("WARN %v", logEntry)
		} else {
			log.Printf("INFO %v", logEntry)
		}
	}
}
