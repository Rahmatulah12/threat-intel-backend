package newrelic

import (
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
)

type Monitor struct {
	app *newrelic.Application
}

func NewMonitor(licenseKey, appName string) (*Monitor, error) {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(appName),
		newrelic.ConfigLicense(licenseKey),
		newrelic.ConfigDistributedTracerEnabled(true),
	)
	if err != nil {
		return nil, err
	}

	return &Monitor{app: app}, nil
}

func (m *Monitor) GetApplication() *newrelic.Application {
	return m.app
}

func (m *Monitor) RecordCustomEvent(eventType string, params map[string]interface{}) {
	m.app.RecordCustomEvent(eventType, params)
}

func (m *Monitor) RecordCustomMetric(name string, value float64) {
	m.app.RecordCustomMetric(name, value)
}

func (m *Monitor) Shutdown() {
	m.app.Shutdown(10 * time.Second)
}

func SetupLogger(app *newrelic.Application) *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	return logger
}
