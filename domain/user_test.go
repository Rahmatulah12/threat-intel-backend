package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	t.Run("creates user successfully", func(t *testing.T) {
		user, err := NewUser("test@example.com", "password123", RoleViewer)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "test@example.com", user.Email)
		assert.Equal(t, RoleViewer, user.Role)
		assert.True(t, user.IsActive)
		assert.NotEmpty(t, user.ID)
		assert.NotEmpty(t, user.PasswordHash)
		assert.NotEqual(t, "password123", user.PasswordHash)
	})

	t.Run("hashes password correctly", func(t *testing.T) {
		user1, _ := NewUser("test1@example.com", "password123", RoleViewer)
		user2, _ := NewUser("test2@example.com", "password123", RoleViewer)

		assert.NotEqual(t, user1.PasswordHash, user2.PasswordHash)
	})
}

func TestUser_ValidatePassword(t *testing.T) {
	user, _ := NewUser("test@example.com", "password123", RoleViewer)

	t.Run("validates correct password", func(t *testing.T) {
		result := user.ValidatePassword("password123")
		assert.True(t, result)
	})

	t.Run("rejects incorrect password", func(t *testing.T) {
		result := user.ValidatePassword("wrongpassword")
		assert.False(t, result)
	})
}

func TestUser_HasPermission(t *testing.T) {
	tests := []struct {
		userRole     UserRole
		requiredRole UserRole
		expected     bool
	}{
		{RoleAdmin, RoleViewer, true},
		{RoleAdmin, RoleAnalyst, true},
		{RoleAdmin, RoleAdmin, true},
		{RoleAnalyst, RoleViewer, true},
		{RoleAnalyst, RoleAnalyst, true},
		{RoleAnalyst, RoleAdmin, false},
		{RoleViewer, RoleViewer, true},
		{RoleViewer, RoleAnalyst, false},
		{RoleViewer, RoleAdmin, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.userRole)+"_vs_"+string(tt.requiredRole), func(t *testing.T) {
			user, _ := NewUser("test@example.com", "password123", tt.userRole)
			result := user.HasPermission(tt.requiredRole)
			assert.Equal(t, tt.expected, result)
		})
	}
}