package domain

import (
	"time"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RoleAnalyst UserRole = "analyst"
	RoleViewer  UserRole = "viewer"
)

type User struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email        string    `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string    `json:"-" gorm:"not null"`
	Role         UserRole  `json:"role" gorm:"not null;default:'viewer'"`
	IsActive     bool      `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func NewUser(email, password string, role UserRole) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         role,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func (u *User) HasPermission(requiredRole UserRole) bool {
	roleHierarchy := map[UserRole]int{
		RoleViewer:  1,
		RoleAnalyst: 2,
		RoleAdmin:   3,
	}
	return roleHierarchy[u.Role] >= roleHierarchy[requiredRole]
}