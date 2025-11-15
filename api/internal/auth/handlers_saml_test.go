package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/streamspace/streamspace/api/internal/db"
	"github.com/streamspace/streamspace/api/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSAMLAuthenticator mocks the SAML authenticator
type MockSAMLAuthenticator struct {
	mock.Mock
}

func (m *MockSAMLAuthenticator) GetMiddleware() interface{} {
	args := m.Called()
	return args.Get(0)
}

func (m *MockSAMLAuthenticator) GetServiceProvider() interface{} {
	args := m.Called()
	return args.Get(0)
}

func (m *MockSAMLAuthenticator) ExtractUserFromAssertion(assertion interface{}) UserAttributes {
	args := m.Called(assertion)
	return args.Get(0).(UserAttributes)
}

// MockUserDB mocks the user database
type MockUserDB struct {
	mock.Mock
}

func (m *MockUserDB) GetUserByEmail(ctx interface{}, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserDB) CreateUser(ctx interface{}, req *models.CreateUserRequest) (*models.User, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserDB) UpdateUser(ctx interface{}, userID string, req *models.UpdateUserRequest) error {
	args := m.Called(ctx, userID, req)
	return args.Error(0)
}

func (m *MockUserDB) GetUserGroups(ctx interface{}, userID string) ([]string, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return []string{}, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

// MockJWTManager mocks the JWT manager
type MockJWTManager struct {
	mock.Mock
}

func (m *MockJWTManager) GenerateToken(userID, username, email, role string, groups []string) (string, error) {
	args := m.Called(userID, username, email, role, groups)
	return args.String(0), args.Error(1)
}

func TestSAMLLogin_NotConfigured(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUserDB := new(MockUserDB)
	mockJWT := new(MockJWTManager)

	// Create handler without SAML (nil)
	handler := NewAuthHandler(mockUserDB, mockJWT, nil)

	// Create test context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/auth/saml/login", nil)

	// Call handler
	handler.SAMLLogin(c)

	// Assert
	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "not configured")
}

func TestSAMLLogin_WithConfiguration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUserDB := new(MockUserDB)
	mockJWT := new(MockJWTManager)
	mockSAML := new(MockSAMLAuthenticator)

	// Mock middleware
	mockMiddleware := &struct{}{}
	mockSAML.On("GetMiddleware").Return(mockMiddleware)

	handler := NewAuthHandler(mockUserDB, mockJWT, mockSAML)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/auth/saml/login?return_url=/dashboard", nil)

	// Note: This test verifies that SAML is called, but we can't test the redirect
	// without a full SAML middleware implementation
	handler.SAMLLogin(c)

	// Cookie should be set
	cookies := w.Result().Cookies()
	var returnURLCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "saml_return_url" {
			returnURLCookie = cookie
			break
		}
	}

	assert.NotNil(t, returnURLCookie)
	assert.Equal(t, "/dashboard", returnURLCookie.Value)
	assert.True(t, returnURLCookie.HttpOnly)
}

func TestSAMLCallback_NotConfigured(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUserDB := new(MockUserDB)
	mockJWT := new(MockJWTManager)

	handler := NewAuthHandler(mockUserDB, mockJWT, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/saml/acs", nil)

	handler.SAMLCallback(c)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "not configured")
}

func TestSAMLCallback_NoAssertion(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUserDB := new(MockUserDB)
	mockJWT := new(MockJWTManager)
	mockSAML := new(MockSAMLAuthenticator)

	handler := NewAuthHandler(mockUserDB, mockJWT, mockSAML)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/saml/acs", nil)
	// No assertion set in context

	handler.SAMLCallback(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "No SAML assertion")
}

func TestSAMLCallback_MissingEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUserDB := new(MockUserDB)
	mockJWT := new(MockJWTManager)
	mockSAML := new(MockSAMLAuthenticator)

	// Mock user attributes with empty email
	mockSAML.On("ExtractUserFromAssertion", mock.Anything).Return(UserAttributes{
		Email:    "", // Missing email
		FullName: "Test User",
		Groups:   []string{},
	})

	handler := NewAuthHandler(mockUserDB, mockJWT, mockSAML)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/saml/acs", nil)
	c.Set("saml_assertion", map[string]interface{}{})

	handler.SAMLCallback(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "missing required email")
}

func TestSAMLCallback_CreateNewUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUserDB := new(MockUserDB)
	mockJWT := new(MockJWTManager)
	mockSAML := new(MockSAMLAuthenticator)

	// Mock user attributes
	mockSAML.On("ExtractUserFromAssertion", mock.Anything).Return(UserAttributes{
		Email:    "test@example.com",
		FullName: "Test User",
		Groups:   []string{"group1"},
	})

	// User doesn't exist
	mockUserDB.On("GetUserByEmail", mock.Anything, "test@example.com").Return(nil, db.ErrUserNotFound)

	// Create new user
	newUser := &models.User{
		ID:       "user123",
		Username: "test@example.com",
		Email:    "test@example.com",
		FullName: "Test User",
		Provider: "saml",
		Role:     "user",
		Active:   true,
	}
	mockUserDB.On("CreateUser", mock.Anything, mock.AnythingOfType("*models.CreateUserRequest")).Return(newUser, nil)

	// Get user groups
	mockUserDB.On("GetUserGroups", mock.Anything, "user123").Return([]string{"group1"}, nil)

	// Generate JWT token
	mockJWT.On("GenerateToken", "user123", "test@example.com", "test@example.com", "user", []string{"group1"}).Return("jwt-token-123", nil)

	handler := NewAuthHandler(mockUserDB, mockJWT, mockSAML)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/saml/acs", nil)
	c.Set("saml_assertion", map[string]interface{}{})

	handler.SAMLCallback(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "jwt-token-123", response["token"])
	assert.Equal(t, "/", response["returnUrl"]) // Default return URL

	mockUserDB.AssertExpectations(t)
	mockJWT.AssertExpectations(t)
	mockSAML.AssertExpectations(t)
}

func TestSAMLCallback_UpdateExistingUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUserDB := new(MockUserDB)
	mockJWT := new(MockJWTManager)
	mockSAML := new(MockSAMLAuthenticator)

	// Mock user attributes
	mockSAML.On("ExtractUserFromAssertion", mock.Anything).Return(UserAttributes{
		Email:    "existing@example.com",
		FullName: "Updated Name",
		Groups:   []string{},
	})

	// User already exists
	existingUser := &models.User{
		ID:       "user456",
		Username: "existing@example.com",
		Email:    "existing@example.com",
		FullName: "Old Name",
		Provider: "saml",
		Role:     "user",
		Active:   true,
	}
	mockUserDB.On("GetUserByEmail", mock.Anything, "existing@example.com").Return(existingUser, nil)

	// Update user
	mockUserDB.On("UpdateUser", mock.Anything, "user456", mock.AnythingOfType("*models.UpdateUserRequest")).Return(nil)

	// Get user groups
	mockUserDB.On("GetUserGroups", mock.Anything, "user456").Return([]string{}, nil)

	// Generate JWT token
	mockJWT.On("GenerateToken", "user456", "existing@example.com", "existing@example.com", "user", []string{}).Return("jwt-token-456", nil)

	handler := NewAuthHandler(mockUserDB, mockJWT, mockSAML)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/saml/acs", nil)
	c.Set("saml_assertion", map[string]interface{}{})

	handler.SAMLCallback(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "jwt-token-456", response["token"])

	mockUserDB.AssertExpectations(t)
	mockJWT.AssertExpectations(t)
	mockSAML.AssertExpectations(t)
}

func TestSAMLCallback_InactiveUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUserDB := new(MockUserDB)
	mockJWT := new(MockJWTManager)
	mockSAML := new(MockSAMLAuthenticator)

	mockSAML.On("ExtractUserFromAssertion", mock.Anything).Return(UserAttributes{
		Email:    "inactive@example.com",
		FullName: "Inactive User",
		Groups:   []string{},
	})

	// User exists but is inactive
	inactiveUser := &models.User{
		ID:       "user789",
		Username: "inactive@example.com",
		Email:    "inactive@example.com",
		FullName: "Inactive User",
		Provider: "saml",
		Role:     "user",
		Active:   false, // Inactive!
	}
	mockUserDB.On("GetUserByEmail", mock.Anything, "inactive@example.com").Return(inactiveUser, nil)
	mockUserDB.On("UpdateUser", mock.Anything, "user789", mock.AnythingOfType("*models.UpdateUserRequest")).Return(nil)

	handler := NewAuthHandler(mockUserDB, mockJWT, mockSAML)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/saml/acs", nil)
	c.Set("saml_assertion", map[string]interface{}{})

	handler.SAMLCallback(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "disabled")
}

func TestSAMLMetadata_NotConfigured(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUserDB := new(MockUserDB)
	mockJWT := new(MockJWTManager)

	handler := NewAuthHandler(mockUserDB, mockJWT, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/auth/saml/metadata", nil)

	handler.SAMLMetadata(c)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "not configured")
}

func TestSAMLMetadata_NilServiceProvider(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUserDB := new(MockUserDB)
	mockJWT := new(MockJWTManager)
	mockSAML := new(MockSAMLAuthenticator)

	// SP is nil
	mockSAML.On("GetServiceProvider").Return(nil)

	handler := NewAuthHandler(mockUserDB, mockJWT, mockSAML)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/auth/saml/metadata", nil)

	handler.SAMLMetadata(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "not initialized")
}

// Benchmark tests
func BenchmarkSAMLLogin(b *testing.B) {
	gin.SetMode(gin.TestMode)

	mockUserDB := new(MockUserDB)
	mockJWT := new(MockJWTManager)

	handler := NewAuthHandler(mockUserDB, mockJWT, nil)

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/auth/saml/login", nil)
		handler.SAMLLogin(c)
	}
}

func BenchmarkSAMLCallback(b *testing.B) {
	gin.SetMode(gin.TestMode)

	mockUserDB := new(MockUserDB)
	mockJWT := new(MockJWTManager)

	handler := NewAuthHandler(mockUserDB, mockJWT, nil)

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/auth/saml/acs", nil)
		handler.SAMLCallback(c)
	}
}
