package application

import (
	"errors"
	"testing"
	"threat-intel-backend/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Save(order *domain.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockOrderRepository) FindByID(id uuid.UUID) (*domain.Order, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Order), args.Error(1)
}

func (m *MockOrderRepository) FindByUserID(userID uuid.UUID) ([]*domain.Order, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Order), args.Error(1)
}

func TestOrderService_CreateOrder(t *testing.T) {
	mockOrderRepo := new(MockOrderRepository)
	mockUserRepo := new(MockUserRepository)
	orderService := NewOrderService(mockOrderRepo, mockUserRepo)

	userID := uuid.New()
	user, _ := domain.NewUser("test@example.com", "password123", domain.RoleViewer)
	user.ID = userID

	t.Run("successful order creation", func(t *testing.T) {
		mockUserRepo.On("FindByID", userID).Return(user, nil).Once()
		mockOrderRepo.On("Save", mock.AnythingOfType("*domain.Order")).Return(nil).Once()

		req := CreateOrderRequest{
			ItemID:   "intel-basic",
			Quantity: 1,
		}

		resp, err := orderService.CreateOrder(userID, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.OrderID)
		assert.Equal(t, domain.OrderStatusConfirmed, resp.Status)
		mockUserRepo.AssertExpectations(t)
		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("invalid item_id", func(t *testing.T) {
		req := CreateOrderRequest{
			ItemID:   "invalid-item",
			Quantity: 1,
		}

		resp, err := orderService.CreateOrder(userID, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "invalid item_id", err.Error())
	})

	t.Run("user not found", func(t *testing.T) {
		mockUserRepo.On("FindByID", userID).Return(nil, errors.New("not found")).Once()

		req := CreateOrderRequest{
			ItemID:   "intel-basic",
			Quantity: 1,
		}

		resp, err := orderService.CreateOrder(userID, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "user not found", err.Error())
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("save order fails", func(t *testing.T) {
		mockUserRepo.On("FindByID", userID).Return(user, nil).Once()
		mockOrderRepo.On("Save", mock.AnythingOfType("*domain.Order")).Return(errors.New("save failed")).Once()

		req := CreateOrderRequest{
			ItemID:   "intel-basic",
			Quantity: 1,
		}

		resp, err := orderService.CreateOrder(userID, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "save failed", err.Error())
		mockUserRepo.AssertExpectations(t)
		mockOrderRepo.AssertExpectations(t)
	})
}

func TestOrderService_GetOrder(t *testing.T) {
	mockOrderRepo := new(MockOrderRepository)
	mockUserRepo := new(MockUserRepository)
	orderService := NewOrderService(mockOrderRepo, mockUserRepo)

	userID := uuid.New()
	orderID := uuid.New()
	user, _ := domain.NewUser("test@example.com", "password123", domain.RoleViewer)
	user.ID = userID

	order := &domain.Order{
		ID:     orderID,
		UserID: userID,
		ItemID: "intel-basic",
		Status: domain.OrderStatusConfirmed,
	}

	t.Run("successful get order by owner", func(t *testing.T) {
		mockOrderRepo.On("FindByID", orderID).Return(order, nil).Once()
		mockUserRepo.On("FindByID", userID).Return(user, nil).Once()

		result, err := orderService.GetOrder(orderID, userID)

		assert.NoError(t, err)
		assert.Equal(t, order, result)
		mockOrderRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("successful get order by analyst", func(t *testing.T) {
		analystID := uuid.New()
		analyst, _ := domain.NewUser("analyst@example.com", "password123", domain.RoleAnalyst)
		analyst.ID = analystID

		mockOrderRepo.On("FindByID", orderID).Return(order, nil).Once()
		mockUserRepo.On("FindByID", analystID).Return(analyst, nil).Once()

		result, err := orderService.GetOrder(orderID, analystID)

		assert.NoError(t, err)
		assert.Equal(t, order, result)
		mockOrderRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("order not found", func(t *testing.T) {
		mockOrderRepo.On("FindByID", orderID).Return(nil, errors.New("not found")).Once()

		result, err := orderService.GetOrder(orderID, userID)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockOrderRepo.On("FindByID", orderID).Return(order, nil).Once()
		mockUserRepo.On("FindByID", userID).Return(nil, errors.New("not found")).Once()

		result, err := orderService.GetOrder(orderID, userID)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockOrderRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("access denied", func(t *testing.T) {
		otherUserID := uuid.New()
		otherUser, _ := domain.NewUser("other@example.com", "password123", domain.RoleViewer)
		otherUser.ID = otherUserID

		mockOrderRepo.On("FindByID", orderID).Return(order, nil).Once()
		mockUserRepo.On("FindByID", otherUserID).Return(otherUser, nil).Once()

		result, err := orderService.GetOrder(orderID, otherUserID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "access denied", err.Error())
		mockOrderRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestOrderService_GetUserOrders(t *testing.T) {
	mockOrderRepo := new(MockOrderRepository)
	mockUserRepo := new(MockUserRepository)
	orderService := NewOrderService(mockOrderRepo, mockUserRepo)

	userID := uuid.New()
	orders := []*domain.Order{
		{ID: uuid.New(), UserID: userID, ItemID: "intel-basic"},
		{ID: uuid.New(), UserID: userID, ItemID: "intel-premium"},
	}

	t.Run("successful get user orders", func(t *testing.T) {
		mockOrderRepo.On("FindByUserID", userID).Return(orders, nil).Once()

		result, err := orderService.GetUserOrders(userID)

		assert.NoError(t, err)
		assert.Equal(t, orders, result)
		assert.Len(t, result, 2)
		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockOrderRepo.On("FindByUserID", userID).Return(nil, errors.New("db error")).Once()

		result, err := orderService.GetUserOrders(userID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "db error", err.Error())
		mockOrderRepo.AssertExpectations(t)
	})
}

func TestNewOrderService(t *testing.T) {
	mockOrderRepo := new(MockOrderRepository)
	mockUserRepo := new(MockUserRepository)

	orderService := NewOrderService(mockOrderRepo, mockUserRepo)

	assert.NotNil(t, orderService)
	assert.Equal(t, mockOrderRepo, orderService.orderRepo)
	assert.Equal(t, mockUserRepo, orderService.userRepo)
}