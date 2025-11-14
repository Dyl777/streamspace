package errors

import (
	"fmt"
	"net/http"
)

// AppError represents a standardized application error
type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	Details    string `json:"details,omitempty"`
	StatusCode int    `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s - %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// ErrorResponse represents the JSON error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

// Error codes
const (
	// Client errors (4xx)
	ErrCodeBadRequest          = "BAD_REQUEST"
	ErrCodeUnauthorized        = "UNAUTHORIZED"
	ErrCodeForbidden           = "FORBIDDEN"
	ErrCodeNotFound            = "NOT_FOUND"
	ErrCodeConflict            = "CONFLICT"
	ErrCodeValidationFailed    = "VALIDATION_FAILED"
	ErrCodeQuotaExceeded       = "QUOTA_EXCEEDED"
	ErrCodeRateLimitExceeded   = "RATE_LIMIT_EXCEEDED"
	ErrCodeSessionNotRunning   = "SESSION_NOT_RUNNING"
	ErrCodeSessionNotFound     = "SESSION_NOT_FOUND"
	ErrCodeTemplateNotFound    = "TEMPLATE_NOT_FOUND"
	ErrCodeUserNotFound        = "USER_NOT_FOUND"
	ErrCodeGroupNotFound       = "GROUP_NOT_FOUND"
	ErrCodeInvalidCredentials  = "INVALID_CREDENTIALS"
	ErrCodeTokenExpired        = "TOKEN_EXPIRED"
	ErrCodeTokenInvalid        = "TOKEN_INVALID"

	// Server errors (5xx)
	ErrCodeInternalServer      = "INTERNAL_SERVER_ERROR"
	ErrCodeDatabaseError       = "DATABASE_ERROR"
	ErrCodeKubernetesError     = "KUBERNETES_ERROR"
	ErrCodeServiceUnavailable  = "SERVICE_UNAVAILABLE"
)

// New creates a new AppError
func New(code string, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: getStatusCodeForErrorCode(code),
	}
}

// NewWithDetails creates a new AppError with details
func NewWithDetails(code string, message string, details string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		Details:    details,
		StatusCode: getStatusCodeForErrorCode(code),
	}
}

// Wrap wraps an existing error with an AppError
func Wrap(code string, message string, err error) *AppError {
	details := ""
	if err != nil {
		details = err.Error()
	}
	return NewWithDetails(code, message, details)
}

// getStatusCodeForErrorCode returns the HTTP status code for an error code
func getStatusCodeForErrorCode(code string) int {
	switch code {
	case ErrCodeBadRequest, ErrCodeValidationFailed:
		return http.StatusBadRequest
	case ErrCodeUnauthorized, ErrCodeInvalidCredentials, ErrCodeTokenExpired, ErrCodeTokenInvalid:
		return http.StatusUnauthorized
	case ErrCodeForbidden, ErrCodeQuotaExceeded:
		return http.StatusForbidden
	case ErrCodeNotFound, ErrCodeSessionNotFound, ErrCodeTemplateNotFound, ErrCodeUserNotFound, ErrCodeGroupNotFound:
		return http.StatusNotFound
	case ErrCodeConflict:
		return http.StatusConflict
	case ErrCodeRateLimitExceeded:
		return http.StatusTooManyRequests
	case ErrCodeServiceUnavailable:
		return http.StatusServiceUnavailable
	case ErrCodeInternalServer, ErrCodeDatabaseError, ErrCodeKubernetesError:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// ToResponse converts AppError to ErrorResponse
func (e *AppError) ToResponse() ErrorResponse {
	return ErrorResponse{
		Error:   e.Code,
		Message: e.Message,
		Code:    e.Code,
		Details: e.Details,
	}
}

// Common error constructors for convenience

func BadRequest(message string) *AppError {
	return New(ErrCodeBadRequest, message)
}

func Unauthorized(message string) *AppError {
	return New(ErrCodeUnauthorized, message)
}

func Forbidden(message string) *AppError {
	return New(ErrCodeForbidden, message)
}

func NotFound(resource string) *AppError {
	return New(ErrCodeNotFound, fmt.Sprintf("%s not found", resource))
}

func Conflict(message string) *AppError {
	return New(ErrCodeConflict, message)
}

func ValidationFailed(message string) *AppError {
	return New(ErrCodeValidationFailed, message)
}

func QuotaExceeded(message string) *AppError {
	return New(ErrCodeQuotaExceeded, message)
}

func SessionNotRunning(sessionID string) *AppError {
	return New(ErrCodeSessionNotRunning, fmt.Sprintf("Session %s is not running", sessionID))
}

func SessionNotFound(sessionID string) *AppError {
	return New(ErrCodeSessionNotFound, fmt.Sprintf("Session %s not found", sessionID))
}

func TemplateNotFound(templateName string) *AppError {
	return New(ErrCodeTemplateNotFound, fmt.Sprintf("Template %s not found", templateName))
}

func UserNotFound(username string) *AppError {
	return New(ErrCodeUserNotFound, fmt.Sprintf("User %s not found", username))
}

func GroupNotFound(groupName string) *AppError {
	return New(ErrCodeGroupNotFound, fmt.Sprintf("Group %s not found", groupName))
}

func InvalidCredentials() *AppError {
	return New(ErrCodeInvalidCredentials, "Invalid username or password")
}

func TokenExpired() *AppError {
	return New(ErrCodeTokenExpired, "Authentication token has expired")
}

func TokenInvalid() *AppError {
	return New(ErrCodeTokenInvalid, "Invalid authentication token")
}

func InternalServer(message string) *AppError {
	return New(ErrCodeInternalServer, message)
}

func DatabaseError(err error) *AppError {
	return Wrap(ErrCodeDatabaseError, "Database operation failed", err)
}

func KubernetesError(err error) *AppError {
	return Wrap(ErrCodeKubernetesError, "Kubernetes operation failed", err)
}

func ServiceUnavailable(service string) *AppError {
	return New(ErrCodeServiceUnavailable, fmt.Sprintf("%s is currently unavailable", service))
}
