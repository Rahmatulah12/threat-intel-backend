package postgres

import (
	"threat-intel-backend/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

type OrderRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *UserRepository) Save(user *domain.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) FindByID(id uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *OrderRepository) Save(order *domain.Order) error {
	return r.db.Save(order).Error
}

func (r *OrderRepository) FindByID(id uuid.UUID) (*domain.Order, error) {
	var order domain.Order
	err := r.db.Preload("User").Where("id = ?", id).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) FindByUserID(userID uuid.UUID) ([]*domain.Order, error) {
	var orders []*domain.Order
	err := r.db.Where("user_id = ?", userID).Find(&orders).Error
	return orders, err
}