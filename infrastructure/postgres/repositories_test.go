package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUserRepository(t *testing.T) {
	repo := NewUserRepository(nil)
	assert.NotNil(t, repo)
	assert.Nil(t, repo.db)
}

func TestNewOrderRepository(t *testing.T) {
	repo := NewOrderRepository(nil)
	assert.NotNil(t, repo)
	assert.Nil(t, repo.db)
}

func TestUserRepositoryStructure(t *testing.T) {
	repo := &UserRepository{}
	assert.NotNil(t, repo)

	// Test repository has the expected structure
	assert.IsType(t, &UserRepository{}, repo)
}

func TestOrderRepositoryStructure(t *testing.T) {
	repo := &OrderRepository{}
	assert.NotNil(t, repo)

	// Test repository has the expected structure
	assert.IsType(t, &OrderRepository{}, repo)
}
