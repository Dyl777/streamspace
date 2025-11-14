package errors

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandler is a middleware that handles errors consistently
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			// Check if it's an AppError
			if appErr, ok := err.Err.(*AppError); ok {
				// Log the error with details
				if appErr.StatusCode >= 500 {
					log.Printf("[ERROR] %s - %s (Details: %s)", appErr.Code, appErr.Message, appErr.Details)
				} else {
					log.Printf("[WARN] %s - %s", appErr.Code, appErr.Message)
				}

				// Send the error response
				c.JSON(appErr.StatusCode, appErr.ToResponse())
				return
			}

			// Handle generic errors
			log.Printf("[ERROR] Unhandled error: %v", err.Err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   ErrCodeInternalServer,
				Message: "An unexpected error occurred",
				Code:    ErrCodeInternalServer,
			})
		}
	}
}

// Recovery is a middleware that recovers from panics
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[PANIC] Recovered from panic: %v", err)

				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Error:   ErrCodeInternalServer,
					Message: "An unexpected error occurred",
					Code:    ErrCodeInternalServer,
				})

				c.Abort()
			}
		}()

		c.Next()
	}
}

// HandleError is a helper function to handle errors in handlers
func HandleError(c *gin.Context, err error) {
	if appErr, ok := err.(*AppError); ok {
		c.Error(appErr)
		c.JSON(appErr.StatusCode, appErr.ToResponse())
	} else {
		internalErr := InternalServer(err.Error())
		c.Error(internalErr)
		c.JSON(internalErr.StatusCode, internalErr.ToResponse())
	}
}

// AbortWithError is a helper to abort request with error
func AbortWithError(c *gin.Context, err *AppError) {
	c.Error(err)
	c.AbortWithStatusJSON(err.StatusCode, err.ToResponse())
}
