package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"threat-intel-backend/application"
	"threat-intel-backend/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(req application.LoginRequest) (*application.AuthResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*application.AuthResponse), args.Error(1)
}

func (m *MockAuthService) Register(req application.RegisterRequest) (*application.AuthResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*application.AuthResponse), args.Error(1)
}

func (m *MockAuthService) RefreshToken(token string) (*application.AuthResponse, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*application.AuthResponse), args.Error(1)
}

type MockOrderService struct {
	mock.Mock
}

func (m *MockOrderService) CreateOrder(userID uuid.UUID, req application.CreateOrderRequest) (*application.OrderResponse, error) {
	args := m.Called(userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*application.OrderResponse), args.Error(1)
}

func (m *MockOrderService) GetOrder(orderID, userID uuid.UUID) (*domain.Order, error) {
	args := m.Called(orderID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Order), args.Error(1)
}

func (m *MockOrderService) GetUserOrders(userID uuid.UUID) ([]*domain.Order, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Order), args.Error(1)
}

func setupHandler() (*Handler, *MockAuthService, *MockOrderService) {
	mockAuth := &MockAuthService{}
	mockOrder := &MockOrderService{}
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	handler := NewHandler(mockAuth, mockOrder, logger)
	return handler, mockAuth, mockOrder
}

func TestHealth(t *testing.T) {
	handler, _, _ := setupHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/health", nil)

	handler.Health(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "healthy", response["status"])
	assert.Equal(t, "threat-intel-backend", response["service"])
}

func TestLogin(t *testing.T) {
	handler, mockAuth, _ := setupHandler()

	t.Run("successful login", func(t *testing.T) {
		req := application.LoginRequest{Email: "test@example.com", Password: "password123"}
		user := &domain.User{ID: uuid.New(), Email: "test@example.com"}
		response := &application.AuthResponse{AccessToken: "token", User: user}

		mockAuth.On("Login", req).Return(response, nil)

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Login(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockAuth.AssertExpectations(t)
	})

	t.Run("invalid request", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer([]byte("invalid")))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Login(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("login failure", func(t *testing.T) {
		req := application.LoginRequest{Email: "test@example.com", Password: "wrong"}
		mockAuth.On("Login", req).Return(nil, errors.New("invalid credentials"))

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Login(c)
	})
}

func TestRegister(t *testing.T) {
	handler, mockAuth, _ := setupHandler()

	t.Run("successful registration", func(t *testing.T) {
		req := application.RegisterRequest{Email: "new@example.com", Password: "password123", Role: domain.RoleViewer}
		user := &domain.User{ID: uuid.New(), Email: "new@example.com"}
		response := &application.AuthResponse{AccessToken: "token", User: user}

		mockAuth.On("Register", req).Return(response, nil)

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Register(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockAuth.AssertExpectations(t)
	})

	t.Run("registration failure", func(t *testing.T) {
		req := application.RegisterRequest{Email: "existing@example.com", Password: "password123", Role: domain.RoleViewer}
		mockAuth.On("Register", req).Return(nil, errors.New("email already exists"))

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Register(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockAuth.AssertExpectations(t)
	})
}

func TestRefreshToken(t *testing.T) {
	handler, mockAuth, _ := setupHandler()

	t.Run("successful refresh", func(t *testing.T) {
		token := "refresh_token"
		user := &domain.User{ID: uuid.New()}
		response := &application.AuthResponse{AccessToken: "new_token", User: user}

		mockAuth.On("RefreshToken", token).Return(response, nil)

		body, _ := json.Marshal(map[string]string{"refresh_token": token})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.RefreshToken(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockAuth.AssertExpectations(t)
	})

	t.Run("invalid token", func(t *testing.T) {
		token := "invalid_token"
		mockAuth.On("RefreshToken", token).Return(nil, errors.New("invalid token"))

		body, _ := json.Marshal(map[string]string{"refresh_token": token})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.RefreshToken(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		mockAuth.AssertExpectations(t)
	})
}

func TestCreateOrder(t *testing.T) {
	handler, _, mockOrder := setupHandler()

	t.Run("successful order creation", func(t *testing.T) {
		userID := uuid.New()
		req := application.CreateOrderRequest{ItemID: "intel-basic", Quantity: 1}
		response := &application.OrderResponse{OrderID: uuid.New().String(), Status: domain.OrderStatusConfirmed}

		mockOrder.On("CreateOrder", userID, req).Return(response, nil)

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/orders", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("user_id", userID)

		handler.CreateOrder(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockOrder.AssertExpectations(t)
	})

	t.Run("missing user_id", func(t *testing.T) {
		req := application.CreateOrderRequest{ItemID: "intel-basic", Quantity: 1}
		body, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/orders", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateOrder(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("order creation failure", func(t *testing.T) {
		userID := uuid.New()
		req := application.CreateOrderRequest{ItemID: "invalid", Quantity: 1}

		mockOrder.On("CreateOrder", userID, req).Return(nil, errors.New("invalid item_id"))

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/orders", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("user_id", userID)

		handler.CreateOrder(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockOrder.AssertExpectations(t)
	})
}

func TestGetOrder(t *testing.T) {
	handler, _, mockOrder := setupHandler()

	t.Run("successful get order", func(t *testing.T) {
		userID := uuid.New()
		orderID := uuid.New()
		order := &domain.Order{ID: orderID, UserID: userID}

		mockOrder.On("GetOrder", orderID, userID).Return(order, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/orders/"+orderID.String(), nil)
		c.Params = gin.Params{{Key: "id", Value: orderID.String()}}
		c.Set("user_id", userID)

		handler.GetOrder(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockOrder.AssertExpectations(t)
	})

	t.Run("invalid order ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/orders/invalid", nil)
		c.Params = gin.Params{{Key: "id", Value: "invalid"}}

		handler.GetOrder(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("order not found", func(t *testing.T) {
		userID := uuid.New()
		orderID := uuid.New()

		mockOrder.On("GetOrder", orderID, userID).Return(nil, errors.New("order not found"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/orders/"+orderID.String(), nil)
		c.Params = gin.Params{{Key: "id", Value: orderID.String()}}
		c.Set("user_id", userID)

		handler.GetOrder(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockOrder.AssertExpectations(t)
	})
}

func TestGetUserOrders(t *testing.T) {
	handler, _, mockOrder := setupHandler()

	t.Run("successful get user orders", func(t *testing.T) {
		userID := uuid.New()
		orders := []*domain.Order{{ID: uuid.New(), UserID: userID}}

		mockOrder.On("GetUserOrders", userID).Return(orders, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/orders", nil)
		c.Set("user_id", userID)

		handler.GetUserOrders(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockOrder.AssertExpectations(t)
	})

	t.Run("missing user_id", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/orders", nil)

		handler.GetUserOrders(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		userID := uuid.New()

		mockOrder.On("GetUserOrders", userID).Return(nil, errors.New("database error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/orders", nil)
		c.Set("user_id", userID)

		handler.GetUserOrders(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockOrder.AssertExpectations(t)
	})
}
