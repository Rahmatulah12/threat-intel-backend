package http

import (
	"net/http"
	"strings"
	"time"
	"threat-intel-backend/domain"
	"threat-intel-backend/infrastructure/jwt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

type JWTServiceInterface interface {
	ValidateAccessToken(token string) (*jwt.Claims, error)
	GenerateAccessToken(userID uuid.UUID, role domain.UserRole) (string, error)
	GenerateRefreshToken(userID uuid.UUID) (string, error)
	ValidateRefreshToken(token string) (uuid.UUID, error)
}

type Middleware struct {
	jwtService JWTServiceInterface
	logger     *logrus.Logger
}

func NewMiddleware(jwtService JWTServiceInterface, logger *logrus.Logger) *Middleware {
	return &Middleware{
		jwtService: jwtService,
		logger:     logger,
	}
}

func (m *Middleware) CORS() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
}

func (m *Middleware) Logger() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		m.logger.WithFields(logrus.Fields{
			"status_code": statusCode,
			"latency":     latency,
			"client_ip":   clientIP,
			"method":      method,
			"path":        path,
		}).Info("HTTP Request")
	})
}

func (m *Middleware) Auth() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			c.Abort()
			return
		}

		claims, err := m.jwtService.ValidateAccessToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)
		c.Next()
	})
}

func (m *Middleware) RequireRole(requiredRole domain.UserRole) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
			c.Abort()
			return
		}

		role := userRole.(domain.UserRole)
		roleHierarchy := map[domain.UserRole]int{
			domain.RoleViewer:  1,
			domain.RoleAnalyst: 2,
			domain.RoleAdmin:   3,
		}

		if roleHierarchy[role] < roleHierarchy[requiredRole] {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	})
}

func (m *Middleware) RateLimit() gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Every(time.Minute), 60)
	
	return gin.HandlerFunc(func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}
		c.Next()
	})
}

func (m *Middleware) NewRelic(app *newrelic.Application) gin.HandlerFunc {
	return nrgin.Middleware(app)
}