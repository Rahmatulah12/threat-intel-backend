package configs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	t.Run("load with default values", func(t *testing.T) {
		config := Load()
		t.Setenv("JWT_SECRET", "your-secret-key-change-in-production")

		assert.Equal(t, "8080", config.Server.Port)
		assert.Equal(t, "0.0.0.0", config.Server.Host)
		assert.Equal(t, "localhost", config.Database.Host)
		assert.Equal(t, "5432", config.Database.Port)
		assert.Equal(t, "postgres", config.Database.User)
		assert.Equal(t, "password", config.Database.Password)
		assert.Equal(t, "threat_intel", config.Database.DBName)
		assert.Equal(t, "disable", config.Database.SSLMode)
		assert.Equal(t, "localhost:6379", config.Redis.Addr)
		assert.Equal(t, "", config.Redis.Password)
		assert.Equal(t, 0, config.Redis.DB)
		assert.Equal(t, "your-secret-key-change-in-production", config.JWT.SecretKey)
		assert.Equal(t, "", config.NewRelic.LicenseKey)
		assert.Equal(t, "zentara-threat-intel-api", config.NewRelic.AppName)
	})

	t.Run("load with environment variables", func(t *testing.T) {
		t.Setenv("SERVER_PORT", "9000")
		t.Setenv("DB_HOST", "testdb")
		t.Setenv("REDIS_DB", "5")
		t.Setenv("JWT_SECRET", "test-secret")

		config := Load()

		assert.Equal(t, "9000", config.Server.Port)
		assert.Equal(t, "testdb", config.Database.Host)
		assert.Equal(t, 5, config.Redis.DB)
		assert.Equal(t, "test-secret", config.JWT.SecretKey)
	})
}

func TestGetEnv(t *testing.T) {
	t.Run("returns environment value when set", func(t *testing.T) {
		t.Setenv("TEST_KEY", "test_value")

		result := getEnv("TEST_KEY", "default")

		assert.Equal(t, "test_value", result)
	})

	t.Run("returns default value when not set", func(t *testing.T) {
		result := getEnv("NON_EXISTENT_KEY", "default")

		assert.Equal(t, "default", result)
	})

	t.Run("returns default value when empty", func(t *testing.T) {
		t.Setenv("EMPTY_KEY", "")

		result := getEnv("EMPTY_KEY", "default")

		assert.Equal(t, "default", result)
	})
}

func TestGetEnvAsInt(t *testing.T) {
	t.Run("returns parsed int when valid", func(t *testing.T) {
		t.Setenv("INT_KEY", "42")

		result := getEnvAsInt("INT_KEY", 10)

		assert.Equal(t, 42, result)
	})

	t.Run("returns default when not set", func(t *testing.T) {
		result := getEnvAsInt("NON_EXISTENT_INT", 10)

		assert.Equal(t, 10, result)
	})

	t.Run("returns default when invalid int", func(t *testing.T) {
		t.Setenv("INVALID_INT", "not_a_number")

		result := getEnvAsInt("INVALID_INT", 10)

		assert.Equal(t, 10, result)
	})

	t.Run("returns default when empty", func(t *testing.T) {
		t.Setenv("EMPTY_INT", "")

		result := getEnvAsInt("EMPTY_INT", 10)

		assert.Equal(t, 10, result)
	})
}
