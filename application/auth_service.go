package application

import (
	"errors"
	"threat-intel-backend/domain"
	"threat-intel-backend/infrastructure/jwt"
)

type AuthService struct {
	userRepo   domain.UserRepository
	jwtService *jwt.Service
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
	Email    string          `json:"email" binding:"required,email"`
	Password string          `json:"password" binding:"required,min=6"`
	Role     domain.UserRole `json:"role" binding:"required"`
}

type AuthResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         *domain.User `json:"user"`
}

func NewAuthService(userRepo domain.UserRepository, jwtService *jwt.Service) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (s *AuthService) Login(req LoginRequest) (*AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.IsActive {
		return nil, errors.New("account is inactive")
	}

	if !user.ValidatePassword(req.Password) {
		return nil, errors.New("invalid credentials")
	}

	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

func (s *AuthService) Register(req RegisterRequest) (*AuthResponse, error) {
	existingUser, _ := s.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	user, err := domain.NewUser(req.Email, req.Password, req.Role)
	if err != nil {
		return nil, err
	}

	if err := s.userRepo.Save(user); err != nil {
		return nil, err
	}

	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (*AuthResponse, error) {
	userID, err := s.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		User:         user,
	}, nil
}
