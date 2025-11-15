package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// RequestIDHeader is the header name for request ID
	RequestIDHeader = "X-Request-ID"

	// RequestIDKey is the context key for request ID
	RequestIDKey = "request_id"
)

// RequestID middleware generates or extracts a correlation ID for each request
// This enables request tracing across distributed systems and log correlation
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try to get request ID from header first (for distributed tracing)
		requestID := c.GetHeader(RequestIDHeader)

		// If not provided, generate a new UUID
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Store in context for use by handlers
		c.Set(RequestIDKey, requestID)

		// Set response header so client can reference this request
		c.Header(RequestIDHeader, requestID)

		c.Next()
	}
}

// GetRequestID retrieves the request ID from the Gin context
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get(RequestIDKey); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}
