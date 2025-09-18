package jwt

import (
	"testing"
	"threat-intel-backend/domain"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	service := NewService("test-secret")

	assert.NotNil(t, service)
	assert.Equal(t, []byte("test-secret"), service.secretKey)
	assert.Equal(t, 15*time.Minute, service.accessTokenTTL)
	assert.Equal(t, 7*24*time.Hour, service.refreshTokenTTL)
}

func TestService_GenerateAccessToken(t *testing.T) {
	service := NewService("test-secret")
	userID := uuid.New()

	t.Run("generates valid access token", func(t *testing.T) {
		token, err := service.GenerateAccessToken(userID, domain.RoleViewer)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("generates tokens with different roles", func(t *testing.T) {
		token1, _ := service.GenerateAccessToken(userID, domain.RoleViewer)
		token2, _ := service.GenerateAccessToken(userID, domain.RoleAdmin)

		assert.NotEqual(t, token1, token2)
	})
}

func TestService_GenerateRefreshToken(t *testing.T) {
	service := NewService("test-secret")
	userID := uuid.New()

	t.Run("generates valid refresh token", func(t *testing.T) {
		token, err := service.GenerateRefreshToken(userID)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("generates tokens for different users", func(t *testing.T) {
		userID2 := uuid.New()
		token1, _ := service.GenerateRefreshToken(userID)
		token2, _ := service.GenerateRefreshToken(userID2)

		assert.NotEqual(t, token1, token2)
	})
}

func TestService_ValidateAccessToken(t *testing.T) {
	service := NewService("test-secret")
	userID := uuid.New()

	t.Run("validates valid access token", func(t *testing.T) {
		token, _ := service.GenerateAccessToken(userID, domain.RoleAdmin)

		claims, err := service.ValidateAccessToken(token)

		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, domain.RoleAdmin, claims.Role)
	})

	t.Run("rejects invalid token", func(t *testing.T) {
		claims, err := service.ValidateAccessToken("invalid-token")

		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("rejects token with wrong secret", func(t *testing.T) {
		wrongService := NewService("wrong-secret")
		token, _ := service.GenerateAccessToken(userID, domain.RoleViewer)

		claims, err := wrongService.ValidateAccessToken(token)

		assert.Error(t, err)
		assert.Nil(t, claims)
	})
}

func TestService_ValidateRefreshToken(t *testing.T) {
	service := NewService("test-secret")
	userID := uuid.New()

	t.Run("validates valid refresh token", func(t *testing.T) {
		token, _ := service.GenerateRefreshToken(userID)

		extractedUserID, err := service.ValidateRefreshToken(token)

		assert.NoError(t, err)
		assert.Equal(t, userID, extractedUserID)
	})

	t.Run("rejects invalid token", func(t *testing.T) {
		extractedUserID, err := service.ValidateRefreshToken("invalid-token")

		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, extractedUserID)
	})

	t.Run("rejects token with wrong secret", func(t *testing.T) {
		wrongService := NewService("wrong-secret")
		token, _ := service.GenerateRefreshToken(userID)

		extractedUserID, err := wrongService.ValidateRefreshToken(token)

		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, extractedUserID)
	})
}