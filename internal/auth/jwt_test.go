package auth

import (
	"testing"
	"time"

	"github.com/lucas-de-lima/go-auth-system/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestJWTService_GenerateAndValidateToken(t *testing.T) {
	jwtService := NewJWTService("test-secret", 1, "test-refresh", 1)
	user := &domain.User{ID: "123", Email: "test@example.com", Roles: []string{"admin"}}
	token, err := jwtService.GenerateToken(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := jwtService.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, claims.UserID)
	assert.Equal(t, user.Email, claims.Email)
	assert.Equal(t, user.Roles, claims.Roles)
}

func TestJWTService_ValidateToken_InvalidToken(t *testing.T) {
	jwtService := NewJWTService("test-secret", 1, "test-refresh", 1)
	_, err := jwtService.ValidateToken("tokeninvalido")
	assert.Error(t, err)
}

func TestJWTService_ValidateToken_Expired(t *testing.T) {
	jwtService := NewJWTService("test-secret", 0, "test-refresh", 1)
	user := &domain.User{ID: "123", Email: "test@example.com", Roles: []string{"admin"}}
	token, err := jwtService.GenerateToken(user)
	assert.NoError(t, err)
	time.Sleep(2 * time.Second)
	_, err = jwtService.ValidateToken(token)
	assert.Error(t, err)
}

func TestJWTService_GenerateAndValidateRefreshToken(t *testing.T) {
	jwtService := NewJWTService("test-secret", 1, "test-refresh", 1)
	token, err := jwtService.GenerateRefreshToken("123")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := jwtService.ValidateRefreshToken(token)
	assert.NoError(t, err)
	assert.Equal(t, "123", claims.Subject)
}

func TestJWTService_ValidateRefreshToken_InvalidToken(t *testing.T) {
	jwtService := NewJWTService("test-secret", 1, "test-refresh", 1)
	_, err := jwtService.ValidateRefreshToken("tokeninvalido")
	assert.Error(t, err)
}

func TestJWTService_ValidateRefreshToken_Expired(t *testing.T) {
	jwtService := NewJWTService("test-secret", 1, "test-refresh", 0)
	token, err := jwtService.GenerateRefreshToken("123")
	assert.NoError(t, err)
	time.Sleep(2 * time.Second)
	_, err = jwtService.ValidateRefreshToken(token)
	assert.Error(t, err)
}

func TestJWTService_Getters(t *testing.T) {
	jwtService := NewJWTService("test-secret", 1, "test-refresh", 1)
	assert.Equal(t, "test-secret", jwtService.GetSecretKey())
	assert.Equal(t, "test-refresh", jwtService.GetRefreshKey())
}
