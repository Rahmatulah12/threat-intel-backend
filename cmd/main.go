package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"threat-intel-backend/application"
	"threat-intel-backend/configs"
	"threat-intel-backend/infrastructure/jwt"
	"threat-intel-backend/infrastructure/newrelic"
	"threat-intel-backend/infrastructure/postgres"
	"threat-intel-backend/infrastructure/redis"
	httpInterface "threat-intel-backend/interfaces/http"
	newrelicAgent "github.com/newrelic/go-agent/v3/newrelic"
)

// @title Zentara Threat Intelligence API
// @version 1.0
// @description A secure threat intelligence backend system
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	config := configs.Load()

	// Initialize New Relic monitoring
	var monitor *newrelic.Monitor
	var err error
	if config.NewRelic.LicenseKey != "" {
		monitor, err = newrelic.NewMonitor(config.NewRelic.LicenseKey, config.NewRelic.AppName)
		if err != nil {
			log.Printf("Failed to initialize New Relic: %v", err)
		}
	}

	// Setup logger
	var logger = newrelic.SetupLogger(nil)
	if monitor != nil {
		logger = newrelic.SetupLogger(monitor.GetApplication())
	}

	// Initialize database
	db, err := postgres.NewConnection(postgres.Config{
		Host:     config.Database.Host,
		Port:     config.Database.Port,
		User:     config.Database.User,
		Password: config.Database.Password,
		DBName:   config.Database.DBName,
		SSLMode:  config.Database.SSLMode,
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	if err := postgres.Migrate(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize Redis
	redisClient := redis.NewClient(config.Redis.Addr, config.Redis.Password, config.Redis.DB)
	if err := redisClient.Ping(context.Background()); err != nil {
		log.Printf("Redis connection failed: %v", err)
	}

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	orderRepo := postgres.NewOrderRepository(db)

	// Initialize services
	jwtService := jwt.NewService(config.JWT.SecretKey)
	authService := application.NewAuthService(userRepo, jwtService)
	orderService := application.NewOrderService(orderRepo, userRepo)

	// Initialize HTTP layer
	middleware := httpInterface.NewMiddleware(jwtService, logger)
	handler := httpInterface.NewHandler(authService, orderService, logger)
	router := httpInterface.NewRouter(handler, middleware)

	// Setup router with New Relic
	var app *newrelicAgent.Application
	if monitor != nil {
		app = monitor.GetApplication()
	}
	r := router.Setup(app)

	// Start server
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port),
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		logger.Infof("Server starting on %s:%s", config.Server.Host, config.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	// Shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	// Cleanup resources
	if monitor != nil {
		monitor.Shutdown()
	}
	redisClient.Close()

	logger.Info("Server exited")
}