package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware_NoToken(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	mockJWT := &JWTManager{secret: []byte("test-secret")}
	mockUserDB := nil // Would be a mock in real tests

	middleware := AuthMiddleware(mockJWT, mockUserDB)

	// Create test router
	router := gin.New()
	router.Use(middleware)
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Execute request without token
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	mockJWT := &JWTManager{secret: []byte("test-secret")}
	mockUserDB := nil

	middleware := AuthMiddleware(mockJWT, mockUserDB)

	// Create test router
	router := gin.New()
	router.Use(middleware)
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Execute request with invalid token
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestOptionalAuthMiddleware_NoToken(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	mockJWT := &JWTManager{secret: []byte("test-secret")}
	mockUserDB := nil

	middleware := OptionalAuthMiddleware(mockJWT, mockUserDB)

	// Create test router
	router := gin.New()
	router.Use(middleware)
	router.GET("/test", func(c *gin.Context) {
		// Check if userID is set
		_, exists := c.Get("userID")
		c.JSON(http.StatusOK, gin.H{"authenticated": exists})
	})

	// Execute request without token
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// Assert - should allow request but not set user
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoleMiddleware_RequiredRole(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	middleware := RoleMiddleware("admin")

	// Create test router
	router := gin.New()
	router.Use(func(c *gin.Context) {
		// Simulate authenticated user with 'user' role
		c.Set("userRole", "user")
		c.Next()
	})
	router.Use(middleware)
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Execute request
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// Assert - should deny access (user role < admin role)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestRoleMiddleware_SufficientRole(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	middleware := RoleMiddleware("user")

	// Create test router
	router := gin.New()
	router.Use(func(c *gin.Context) {
		// Simulate authenticated admin user
		c.Set("userRole", "admin")
		c.Next()
	})
	router.Use(middleware)
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Execute request
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// Assert - should allow access (admin role >= user role)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoleMiddleware_NoRoleSet(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	middleware := RoleMiddleware("user")

	// Create test router
	router := gin.New()
	router.Use(middleware)
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Execute request without role
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// Assert - should deny access
	assert.Equal(t, http.StatusForbidden, w.Code)
}

// Benchmark tests
func BenchmarkAuthMiddleware(b *testing.B) {
	gin.SetMode(gin.TestMode)

	mockJWT := &JWTManager{secret: []byte("test-secret")}
	middleware := AuthMiddleware(mockJWT, nil)

	router := gin.New()
	router.Use(middleware)
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)
	}
}
