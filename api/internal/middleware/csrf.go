package middleware

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// CSRF Constants
const (
	// CSRFTokenLength is the length of CSRF tokens in bytes
	CSRFTokenLength = 32

	// CSRFTokenHeader is the HTTP header for CSRF tokens
	CSRFTokenHeader = "X-CSRF-Token"

	// CSRFCookieName is the name of the CSRF cookie
	CSRFCookieName = "csrf_token"

	// CSRFTokenExpiry is how long CSRF tokens are valid
	CSRFTokenExpiry = 24 * time.Hour
)

// CSRFStore stores CSRF tokens with expiration
type CSRFStore struct {
	tokens map[string]time.Time
	mu     sync.RWMutex
}

var (
	globalCSRFStore = &CSRFStore{
		tokens: make(map[string]time.Time),
	}
	csrfCleanupOnce sync.Once
)

// generateCSRFToken generates a random CSRF token
func generateCSRFToken() (string, error) {
	bytes := make([]byte, CSRFTokenLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// addToken adds a token to the store with expiration
func (cs *CSRFStore) addToken(token string) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.tokens[token] = time.Now().Add(CSRFTokenExpiry)
}

// validateToken checks if a token is valid and not expired
func (cs *CSRFStore) validateToken(token string) bool {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	
	expiry, exists := cs.tokens[token]
	if !exists {
		return false
	}
	
	// Check if expired
	if time.Now().After(expiry) {
		return false
	}
	
	return true
}

// removeToken removes a token from the store
func (cs *CSRFStore) removeToken(token string) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	delete(cs.tokens, token)
}

// cleanup removes expired tokens
func (cs *CSRFStore) cleanup() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		cs.mu.Lock()
		now := time.Now()
		for token, expiry := range cs.tokens {
			if now.After(expiry) {
				delete(cs.tokens, token)
			}
		}
		cs.mu.Unlock()
	}
}

// CSRFProtection middleware validates CSRF tokens for state-changing requests
func CSRFProtection() gin.HandlerFunc {
	// Start cleanup goroutine once
	csrfCleanupOnce.Do(func() {
		go globalCSRFStore.cleanup()
	})

	return func(c *gin.Context) {
		// Skip CSRF for safe methods (GET, HEAD, OPTIONS)
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
			// For GET requests, generate and set a CSRF token
			token, err := generateCSRFToken()
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to generate CSRF token",
				})
				return
			}

			// Store token
			globalCSRFStore.addToken(token)

			// Set token in response header
			c.Header(CSRFTokenHeader, token)

			// Set token in cookie (HttpOnly for security)
			c.SetCookie(
				CSRFCookieName,
				token,
				int(CSRFTokenExpiry.Seconds()),
				"/",
				"",
				true,  // Secure (HTTPS only in production)
				true,  // HttpOnly
			)

			c.Next()
			return
		}

		// For state-changing methods (POST, PUT, DELETE, PATCH), validate CSRF token
		// Get token from header
		headerToken := c.GetHeader(CSRFTokenHeader)
		
		// Get token from cookie
		cookieToken, err := c.Cookie(CSRFCookieName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   "CSRF token missing",
				"message": "CSRF cookie not found",
			})
			return
		}

		// Tokens must match
		if subtle.ConstantTimeCompare([]byte(headerToken), []byte(cookieToken)) != 1 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   "CSRF token mismatch",
				"message": "CSRF tokens do not match",
			})
			return
		}

		// Validate token exists and is not expired
		if !globalCSRFStore.validateToken(cookieToken) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   "CSRF token invalid",
				"message": "CSRF token has expired or is invalid",
			})
			return
		}

		c.Next()
	}
}

// GetCSRFToken returns the current CSRF token for the request
// Useful for rendering in HTML forms or passing to frontend
func GetCSRFToken(c *gin.Context) string {
	return c.GetHeader(CSRFTokenHeader)
}
