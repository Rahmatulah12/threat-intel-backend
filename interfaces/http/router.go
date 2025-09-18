package http

import (
	"threat-intel-backend/domain"
	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/newrelic"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	handler    *Handler
	middleware *Middleware
}

func NewRouter(handler *Handler, middleware *Middleware) *Router {
	return &Router{
		handler:    handler,
		middleware: middleware,
	}
}

func (r *Router) Setup(app *newrelic.Application) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Global middleware
	router.Use(r.middleware.NewRelic(app))
	router.Use(r.middleware.CORS())
	router.Use(r.middleware.Logger())
	router.Use(r.middleware.RateLimit())
	router.Use(gin.Recovery())

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check
	router.GET("/health", r.handler.Health)

	// Auth routes
	auth := router.Group("/auth")
	{
		auth.POST("/login", r.handler.Login)
		auth.POST("/register", r.handler.Register)
		auth.POST("/refresh", r.handler.RefreshToken)
	}

	// Protected routes
	api := router.Group("/api/v1")
	api.Use(r.middleware.Auth())
	{
		// Order routes
		orders := api.Group("/orders")
		{
			orders.POST("", r.handler.CreateOrder)
			orders.GET("", r.handler.GetUserOrders)
			orders.GET("/:id", r.handler.GetOrder)
		}

		// Admin routes
		admin := api.Group("/admin")
		admin.Use(r.middleware.RequireRole(domain.RoleAdmin))
		{
			admin.GET("/users", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Admin endpoint - list users"})
			})
		}

		// Analyst routes
		analyst := api.Group("/analyst")
		analyst.Use(r.middleware.RequireRole(domain.RoleAnalyst))
		{
			analyst.GET("/reports", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Analyst endpoint - view reports"})
			})
		}
	}

	return router
}