package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *Router {
	mockAuth := &MockAuthService{}
	mockOrder := &MockOrderService{}
	mockJWT := &MockJWTService{}
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	handler := NewHandler(mockAuth, mockOrder, logger)
	middleware := NewMiddleware(mockJWT, logger)

	return NewRouter(handler, middleware)
}

func TestNewRouter(t *testing.T) {
	router := setupRouter()

	assert.NotNil(t, router)
	assert.NotNil(t, router.handler)
	assert.NotNil(t, router.middleware)
}

func TestRouterSetup(t *testing.T) {
	router := setupRouter()
	engine := router.Setup(nil)

	assert.NotNil(t, engine)
}

func TestHealthRoute(t *testing.T) {
	router := setupRouter()
	engine := router.Setup(nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/health", nil)

	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSwaggerRoute(t *testing.T) {
	router := setupRouter()
	engine := router.Setup(nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/swagger/index.html", nil)

	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthRoutes(t *testing.T) {
	router := setupRouter()
	engine := router.Setup(nil)

	routes := []struct {
		method string
		path   string
	}{
		{"POST", "/auth/login"},
		{"POST", "/auth/register"},
		{"POST", "/auth/refresh"},
	}

	for _, route := range routes {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(route.method, route.path, nil)

		engine.ServeHTTP(w, req)

		assert.NotEqual(t, http.StatusNotFound, w.Code)
	}
}

func TestProtectedRoutes(t *testing.T) {
	router := setupRouter()
	engine := router.Setup(nil)

	routes := []struct {
		method string
		path   string
	}{
		{"POST", "/api/v1/orders"},
		{"GET", "/api/v1/orders"},
		{"GET", "/api/v1/orders/123"},
	}

	for _, route := range routes {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(route.method, route.path, nil)

		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	}
}

func TestAdminRoutes(t *testing.T) {
	router := setupRouter()
	engine := router.Setup(nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/admin/users", nil)

	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAnalystRoutes(t *testing.T) {
	router := setupRouter()
	engine := router.Setup(nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/analyst/reports", nil)

	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
