package application

import (
	"errors"
	"testing"
	"threat-intel-backend/domain"
	"threat-intel-backend/infrastructure/jwt"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Save(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(id uuid.UUID) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func TestAuthService_Login(t *testing.T) {
	mockRepo := new(MockUserRepository)
	jwtService := jwt.NewService("test-secret")
	authService := NewAuthService(mockRepo, jwtService)

	user, _ := domain.NewUser("test@example.com", "password123", domain.RoleViewer)
	user.ID = uuid.New()

	t.Run("successful login", func(t *testing.T) {
		mockRepo.On("FindByEmail", "test@example.com").Return(user, nil).Once()

		req := LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		resp, err := authService.Login(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.AccessToken)
		assert.NotEmpty(t, resp.RefreshToken)
		assert.Equal(t, user, resp.User)
		mockRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo.On("FindByEmail", "notfound@example.com").Return(nil, errors.New("not found")).Once()

		req := LoginRequest{
			Email:    "notfound@example.com",
			Password: "password123",
		}

		resp, err := authService.Login(req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "invalid credentials", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("inactive user", func(t *testing.T) {
		inactiveUser := *user
		inactiveUser.IsActive = false
		mockRepo.On("FindByEmail", "test@example.com").Return(&inactiveUser, nil).Once()

		req := LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		resp, err := authService.Login(req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "account is inactive", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid password", func(t *testing.T) {
		mockRepo.On("FindByEmail", "test@example.com").Return(user, nil).Once()

		req := LoginRequest{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}

		resp, err := authService.Login(req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "invalid credentials", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestAuthService_Register(t *testing.T) {
	mockRepo := new(MockUserRepository)
	jwtService := jwt.NewService("test-secret")
	authService := NewAuthService(mockRepo, jwtService)

	t.Run("successful registration", func(t *testing.T) {
		mockRepo.On("FindByEmail", "new@example.com").Return(nil, errors.New("not found")).Once()
		mockRepo.On("Save", mock.AnythingOfType("*domain.User")).Return(nil).Once()

		req := RegisterRequest{
			Email:    "new@example.com",
			Password: "password123",
			Role:     domain.RoleViewer,
		}

		resp, err := authService.Register(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.AccessToken)
		assert.NotEmpty(t, resp.RefreshToken)
		assert.Equal(t, "new@example.com", resp.User.Email)
		assert.Equal(t, domain.RoleViewer, resp.User.Role)
		mockRepo.AssertExpectations(t)
	})

	t.Run("email already exists", func(t *testing.T) {
		existingUser, _ := domain.NewUser("existing@example.com", "password123", domain.RoleViewer)
		mockRepo.On("FindByEmail", "existing@example.com").Return(existingUser, nil).Once()

		req := RegisterRequest{
			Email:    "existing@example.com",
			Password: "password123",
			Role:     domain.RoleViewer,
		}

		resp, err := authService.Register(req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "email already exists", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("save user fails", func(t *testing.T) {
		mockRepo.On("FindByEmail", "new@example.com").Return(nil, errors.New("not found")).Once()
		mockRepo.On("Save", mock.AnythingOfType("*domain.User")).Return(errors.New("save failed")).Once()

		req := RegisterRequest{
			Email:    "new@example.com",
			Password: "password123",
			Role:     domain.RoleViewer,
		}

		resp, err := authService.Register(req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "save failed", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestAuthService_RefreshToken(t *testing.T) {
	mockRepo := new(MockUserRepository)
	jwtService := jwt.NewService("test-secret")
	authService := NewAuthService(mockRepo, jwtService)

	user, _ := domain.NewUser("test@example.com", "password123", domain.RoleViewer)
	user.ID = uuid.New()

	t.Run("successful token refresh", func(t *testing.T) {
		refreshToken, _ := jwtService.GenerateRefreshToken(user.ID)
		mockRepo.On("FindByID", user.ID).Return(user, nil).Once()

		resp, err := authService.RefreshToken(refreshToken)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.AccessToken)
		assert.NotEmpty(t, resp.RefreshToken)
		assert.Equal(t, user, resp.User)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid refresh token", func(t *testing.T) {
		resp, err := authService.RefreshToken("invalid-token")

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "invalid refresh token", err.Error())
	})

	t.Run("user not found", func(t *testing.T) {
		refreshToken, _ := jwtService.GenerateRefreshToken(user.ID)
		mockRepo.On("FindByID", user.ID).Return(nil, errors.New("not found")).Once()

		resp, err := authService.RefreshToken(refreshToken)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "user not found", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestNewAuthService(t *testing.T) {
	mockRepo := new(MockUserRepository)
	jwtService := jwt.NewService("test-secret")

	authService := NewAuthService(mockRepo, jwtService)

	assert.NotNil(t, authService)
	assert.Equal(t, mockRepo, authService.userRepo)
	assert.Equal(t, jwtService, authService.jwtService)
}