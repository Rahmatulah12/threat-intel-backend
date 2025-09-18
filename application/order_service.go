package application

import (
	"errors"
	"threat-intel-backend/domain"
	"github.com/google/uuid"
)

type OrderService struct {
	orderRepo domain.OrderRepository
	userRepo  domain.UserRepository
}

type CreateOrderRequest struct {
	ItemID   string `json:"item_id" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,min=1"`
}

type OrderResponse struct {
	OrderID string             `json:"order_id"`
	Status  domain.OrderStatus `json:"status"`
}

var validItems = map[string]bool{
	"intel-basic":    true,
	"intel-premium":  true,
	"intel-enterprise": true,
}

func NewOrderService(orderRepo domain.OrderRepository, userRepo domain.UserRepository) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
		userRepo:  userRepo,
	}
}

func (s *OrderService) CreateOrder(userID uuid.UUID, req CreateOrderRequest) (*OrderResponse, error) {
	if !validItems[req.ItemID] {
		return nil, errors.New("invalid item_id")
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !user.HasPermission(domain.RoleViewer) {
		return nil, errors.New("insufficient permissions")
	}

	orderAggregate := domain.NewOrder(userID, req.ItemID, req.Quantity)
	orderAggregate.Confirm()

	if err := s.orderRepo.Save(orderAggregate.Order); err != nil {
		return nil, err
	}

	return &OrderResponse{
		OrderID: orderAggregate.Order.ID.String(),
		Status:  orderAggregate.Order.Status,
	}, nil
}

func (s *OrderService) GetOrder(orderID uuid.UUID, userID uuid.UUID) (*domain.Order, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if order.UserID != userID && !user.HasPermission(domain.RoleAnalyst) {
		return nil, errors.New("access denied")
	}

	return order, nil
}

func (s *OrderService) GetUserOrders(userID uuid.UUID) ([]*domain.Order, error) {
	return s.orderRepo.FindByUserID(userID)
}