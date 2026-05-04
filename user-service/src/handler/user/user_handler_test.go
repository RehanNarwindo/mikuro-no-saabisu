package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestPublicHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	PublicHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "message")
}

func TestGetUserProfileHandler_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	GetUserProfileHandler(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetUserProfileHandler_Authorized(t *testing.T) {
	t.Skip("Skipping - requires database connection. Will test with mock later.")

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	claims := jwt.MapClaims{
		"user_id": "123e4567-e89b-12d3-a456-426614174000",
		"email":   "test@example.com",
		"role":    "user",
	}
	c.Set("user", claims)

	GetUserProfileHandler(c)

	assert.NotNil(t, w)
}

func TestGetUserByIdHandler_InvalidUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = gin.Params{
		{Key: "id", Value: "invalid-uuid"},
	}

	claims := jwt.MapClaims{
		"user_id": "123e4567-e89b-12d3-a456-426614174000",
		"role":    "admin",
	}
	c.Set("user", claims)

	GetUserByIdHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}
	assert.Contains(t, response["message"], "Invalid user ID")
}

func TestGetUserByIdHandler_ValidUUID(t *testing.T) {
	t.Skip("Skipping - requires database connection. Will test with mock later.")

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	validUUID := "123e4567-e89b-12d3-a456-426614174000"
	c.Params = gin.Params{
		{Key: "id", Value: validUUID},
	}

	claims := jwt.MapClaims{
		"user_id": "123e4567-e89b-12d3-a456-426614174000",
		"role":    "admin",
	}
	c.Set("user", claims)

	GetUserByIdHandler(c)

	assert.NotEqual(t, http.StatusBadRequest, w.Code)
}

func TestGetAllUserHandler_InvalidLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("GET", "/users?limit=-1", nil)

	claims := jwt.MapClaims{
		"user_id": "123e4567-e89b-12d3-a456-426614174000",
		"role":    "admin",
	}
	c.Set("user", claims)

	GetAllUserHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetAllUserHandler_InvalidSortDir(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("GET", "/users?sort_dir=INVALID", nil)

	claims := jwt.MapClaims{
		"user_id": "123e4567-e89b-12d3-a456-426614174000",
		"role":    "admin",
	}
	c.Set("user", claims)

	GetAllUserHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetAllUserHandler_InvalidSortBy(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("GET", "/users?sort_by=invalid_field", nil)

	claims := jwt.MapClaims{
		"user_id": "123e4567-e89b-12d3-a456-426614174000",
		"role":    "admin",
	}
	c.Set("user", claims)

	GetAllUserHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteUserHandler_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	DeleteUserHandler(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetClaimsFromContext_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	expectedClaims := jwt.MapClaims{
		"user_id": "123",
		"email":   "test@example.com",
	}
	c.Set("user", expectedClaims)

	claims, err := getClaimsFromContext(c)

	assert.NoError(t, err)
	assert.Equal(t, expectedClaims["user_id"], claims["user_id"])
	assert.Equal(t, expectedClaims["email"], claims["email"])
}

func TestGetClaimsFromContext_NoUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	claims, err := getClaimsFromContext(c)

	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Equal(t, "unauthorized", err.Error())
}

func TestGetClaimsFromContext_InvalidClaimsType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set("user", "not-a-map")

	claims, err := getClaimsFromContext(c)

	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Equal(t, "invalid claims format", err.Error())
}
