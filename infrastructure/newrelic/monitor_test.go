package newrelic

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewMonitor(t *testing.T) {
	t.Run("handles invalid license key", func(t *testing.T) {
		monitor, err := NewMonitor("invalid-key", "test-app")

		assert.Error(t, err)
		assert.Nil(t, monitor)
	})

	t.Run("handles empty license key", func(t *testing.T) {
		monitor, err := NewMonitor("", "test-app")

		assert.Error(t, err)
		assert.Nil(t, monitor)
	})
}

func TestMonitor_GetApplication(t *testing.T) {
	validLicense := "1234567890123456789012345678901234567890"
	monitor, err := NewMonitor(validLicense, "test-app")
	if err != nil {
		t.Skip("Skipping test due to NewRelic configuration requirements")
	}

	app := monitor.GetApplication()

	assert.NotNil(t, app)
	assert.Equal(t, monitor.app, app)
}

func TestMonitor_RecordCustomEvent(t *testing.T) {
	validLicense := "1234567890123456789012345678901234567890"
	monitor, err := NewMonitor(validLicense, "test-app")
	if err != nil {
		t.Skip("Skipping test due to NewRelic configuration requirements")
	}

	params := map[string]interface{}{
		"user_id": "123",
		"action":  "login",
	}

	assert.NotPanics(t, func() {
		monitor.RecordCustomEvent("UserAction", params)
	})
}

func TestMonitor_RecordCustomMetric(t *testing.T) {
	validLicense := "1234567890123456789012345678901234567890"
	monitor, err := NewMonitor(validLicense, "test-app")
	if err != nil {
		t.Skip("Skipping test due to NewRelic configuration requirements")
	}

	assert.NotPanics(t, func() {
		monitor.RecordCustomMetric("custom.metric", 42.5)
	})
}

func TestMonitor_Shutdown(t *testing.T) {
	validLicense := "1234567890123456789012345678901234567890"
	monitor, err := NewMonitor(validLicense, "test-app")
	if err != nil {
		t.Skip("Skipping test due to NewRelic configuration requirements")
	}

	assert.NotPanics(t, func() {
		monitor.Shutdown()
	})
}

func TestSetupLogger(t *testing.T) {
	validLicense := "1234567890123456789012345678901234567890"
	monitor, err := NewMonitor(validLicense, "test-app")
	if err != nil {
		t.Skip("Skipping test due to NewRelic configuration requirements")
	}

	logger := SetupLogger(monitor.GetApplication())

	assert.NotNil(t, logger)
	assert.IsType(t, &logrus.Logger{}, logger)
	assert.IsType(t, &logrus.JSONFormatter{}, logger.Formatter)
}