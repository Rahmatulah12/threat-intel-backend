package domain

import (
	"time"
	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	ID        uuid.UUID   `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID   `json:"user_id" gorm:"type:uuid;not null"`
	ItemID    string      `json:"item_id" gorm:"not null"`
	Quantity  int         `json:"quantity" gorm:"not null;default:1"`
	Status    OrderStatus `json:"status" gorm:"not null;default:'pending'"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	User      User        `json:"user" gorm:"foreignKey:UserID"`
}

type OrderAggregate struct {
	Order *Order
}

func NewOrder(userID uuid.UUID, itemID string, quantity int) *OrderAggregate {
	order := &Order{
		ID:        uuid.New(),
		UserID:    userID,
		ItemID:    itemID,
		Quantity:  quantity,
		Status:    OrderStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return &OrderAggregate{Order: order}
}

func (oa *OrderAggregate) Confirm() {
	oa.Order.Status = OrderStatusConfirmed
	oa.Order.UpdatedAt = time.Now()
}

func (oa *OrderAggregate) Complete() {
	oa.Order.Status = OrderStatusCompleted
	oa.Order.UpdatedAt = time.Now()
}

func (oa *OrderAggregate) Cancel() {
	oa.Order.Status = OrderStatusCancelled
	oa.Order.UpdatedAt = time.Now()
}

type OrderRepository interface {
	Save(order *Order) error
	FindByID(id uuid.UUID) (*Order, error)
	FindByUserID(userID uuid.UUID) ([]*Order, error)
}

type UserRepository interface {
	Save(user *User) error
	FindByID(id uuid.UUID) (*User, error)
	FindByEmail(email string) (*User, error)
}