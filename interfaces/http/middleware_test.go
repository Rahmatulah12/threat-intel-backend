package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"threat-intel-backend/domain"
	"threat-intel-backend/infrastructure/jwt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) ValidateAccessToken(token string) (*jwt.Claims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwt.Claims), args.Error(1)
}

func (m *MockJWTService) GenerateAccessToken(userID uuid.UUID, role domain.UserRole) (string, error) {
	args := m.Called(userID, role)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) GenerateRefreshToken(userID uuid.UUID) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) ValidateRefreshToken(token string) (uuid.UUID, error) {
	args := m.Called(token)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func setupMiddleware() (*Middleware, *MockJWTService) {
	mockJWT := &MockJWTService{}
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	middleware := NewMiddleware(mockJWT, logger)
	return middleware, mockJWT
}

func TestCORS(t *testing.T) {
	middleware, _ := setupMiddleware()

	t.Run("sets CORS headers", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)

		middleware.CORS()(c)

		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
		assert.Contains(t, w.Header().Get("Access-Control-Allow-Headers"), "Authorization")
		assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "POST")
	})

	t.Run("handles OPTIONS request", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("OPTIONS", "/test", nil)

		middleware.CORS()(c)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.True(t, c.IsAborted())
	})
}

func TestLogger(t *testing.T) {
	middleware, _ := setupMiddleware()

	t.Run("logs request", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, engine := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test?param=value", nil)

		called := false
		engine.Use(middleware.Logger())
		engine.GET("/test", func(c *gin.Context) {
			called = true
			c.Status(200)
		})

		engine.ServeHTTP(w, c.Request)

		assert.True(t, called)
	})
}

func TestAuth(t *testing.T) {
	middleware, mockJWT := setupMiddleware()

	t.Run("missing authorization header", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)

		middleware.Auth()(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.True(t, c.IsAborted())
	})

	t.Run("invalid bearer format", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)
		c.Request.Header.Set("Authorization", "InvalidFormat token")

		middleware.Auth()(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.True(t, c.IsAborted())
	})

	t.Run("invalid token", func(t *testing.T) {
		mockJWT.On("ValidateAccessToken", "invalid_token").Return(nil, errors.New("invalid token"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)
		c.Request.Header.Set("Authorization", "Bearer invalid_token")

		middleware.Auth()(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.True(t, c.IsAborted())
		mockJWT.AssertExpectations(t)
	})

	t.Run("valid token", func(t *testing.T) {
		userID := uuid.New()
		claims := &jwt.Claims{UserID: userID, Role: domain.RoleViewer}
		mockJWT.On("ValidateAccessToken", "valid_token").Return(claims, nil)

		w := httptest.NewRecorder()
		c, engine := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)
		c.Request.Header.Set("Authorization", "Bearer valid_token")

		called := false
		engine.Use(middleware.Auth())
		engine.GET("/test", func(c *gin.Context) {
			called = true
			c.Status(200)
		})

		engine.ServeHTTP(w, c.Request)

		assert.True(t, called)
		assert.Equal(t, http.StatusOK, w.Code)

		mockJWT.AssertExpectations(t)
	})
}

func TestRequireRole(t *testing.T) {
	middleware, _ := setupMiddleware()

	t.Run("missing user role", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)

		middleware.RequireRole(domain.RoleAdmin)(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.True(t, c.IsAborted())
	})

	t.Run("insufficient permissions", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)
		c.Set("user_role", domain.RoleViewer)

		middleware.RequireRole(domain.RoleAdmin)(c)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.True(t, c.IsAborted())
	})

	t.Run("sufficient permissions - same role", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, engine := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)
		c.Set("user_role", domain.RoleAnalyst)

		engine.Use(middleware.RequireRole(domain.RoleAnalyst))
		engine.GET("/test", func(c *gin.Context) {
			c.Status(200)
		})

		engine.ServeHTTP(w, c.Request)
	})

	t.Run("sufficient permissions - higher role", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, engine := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)
		c.Set("user_role", domain.RoleAdmin)

		engine.Use(middleware.RequireRole(domain.RoleViewer))
		engine.GET("/test", func(c *gin.Context) {
			c.Status(200)
		})

		engine.ServeHTTP(w, c.Request)
	})
}

func TestRateLimit(t *testing.T) {
	middleware, _ := setupMiddleware()

	t.Run("allows requests within limit", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, engine := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)

		engine.Use(middleware.RateLimit())
		engine.GET("/test", func(c *gin.Context) {
			c.Status(200)
		})

		engine.ServeHTTP(w, c.Request)
	})
}
