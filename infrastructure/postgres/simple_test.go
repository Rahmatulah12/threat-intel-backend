package postgres

import (
	"testing"
)

func TestConfig(t *testing.T) {
	config := Config{
		Host:     "localhost",
		Port:     "5432",
		User:     "test",
		Password: "test",
		DBName:   "test",
		SSLMode:  "disable",
	}

	if config.Host != "localhost" {
		t.Errorf("Expected Host to be localhost, got %s", config.Host)
	}
	if config.Port != "5432" {
		t.Errorf("Expected Port to be 5432, got %s", config.Port)
	}
}

func TestRepositoryCreation(t *testing.T) {
	userRepo := NewUserRepository(nil)
	if userRepo == nil {
		t.Error("Expected UserRepository to be created")
	}

	orderRepo := NewOrderRepository(nil)
	if orderRepo == nil {
		t.Error("Expected OrderRepository to be created")
	}
}