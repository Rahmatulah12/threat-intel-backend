package http

import (
	"net/http"
	"threat-intel-backend/application"
	"threat-intel-backend/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type AuthServiceInterface interface {
	Login(req application.LoginRequest) (*application.AuthResponse, error)
	Register(req application.RegisterRequest) (*application.AuthResponse, error)
	RefreshToken(token string) (*application.AuthResponse, error)
}

type OrderServiceInterface interface {
	CreateOrder(userID uuid.UUID, req application.CreateOrderRequest) (*application.OrderResponse, error)
	GetOrder(orderID, userID uuid.UUID) (*domain.Order, error)
	GetUserOrders(userID uuid.UUID) ([]*domain.Order, error)
}

type Handler struct {
	authService  AuthServiceInterface
	orderService OrderServiceInterface
	logger       *logrus.Logger
}

func NewHandler(authService AuthServiceInterface, orderService OrderServiceInterface, logger *logrus.Logger) *Handler {
	return &Handler{
		authService:  authService,
		orderService: orderService,
		logger:       logger,
	}
}

// @Summary Health check
// @Description Check if the service is running
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"service": "threat-intel-backend",
	})
}

// @Summary User login
// @Description Authenticate user and return JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body application.LoginRequest true "Login credentials"
// @Success 200 {object} application.AuthResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req application.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authService.Login(req)
	if err != nil {
		h.logger.WithError(err).Error("Login failed")
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	h.logger.WithField("user_id", response.User.ID).Info("User logged in")
	c.JSON(http.StatusOK, response)
}

// @Summary User registration
// @Description Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body application.RegisterRequest true "Registration data"
// @Success 201 {object} application.AuthResponse
// @Failure 400 {object} map[string]string
// @Router /auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req application.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authService.Register(req)
	if err != nil {
		h.logger.WithError(err).Error("Registration failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.WithField("user_id", response.User.ID).Info("User registered")
	c.JSON(http.StatusCreated, response)
}

// @Summary Refresh token
// @Description Refresh access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Refresh token"
// @Success 200 {object} application.AuthResponse
// @Failure 400 {object} map[string]string
// @Router /auth/refresh [post]
func (h *Handler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Create order
// @Description Create a new order for threat intelligence data
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body application.CreateOrderRequest true "Order data"
// @Success 201 {object} application.OrderResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /orders [post]
func (h *Handler) CreateOrder(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	var req application.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.orderService.CreateOrder(userID.(uuid.UUID), req)
	if err != nil {
		h.logger.WithError(err).Error("Order creation failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"user_id": userID,
		"order_id": response.OrderID,
	}).Info("Order created")

	c.JSON(http.StatusCreated, response)
}

// @Summary Get order
// @Description Get order by ID
// @Tags orders
// @Produce json
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Success 200 {object} domain.Order
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /orders/{id} [get]
func (h *Handler) GetOrder(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	order, err := h.orderService.GetOrder(orderID, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// @Summary Get user orders
// @Description Get all orders for the authenticated user
// @Tags orders
// @Produce json
// @Security BearerAuth
// @Success 200 {array} domain.Order
// @Failure 401 {object} map[string]string
// @Router /orders [get]
func (h *Handler) GetUserOrders(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	orders, err := h.orderService.GetUserOrders(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}