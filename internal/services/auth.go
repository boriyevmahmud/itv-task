package services

import (
	"errors"
	"itv-task/config"
	"itv-task/internal/models"
	"itv-task/pkg/utils"

	"go.uber.org/fx"
)

type AuthService struct {
	Config *config.Config
}

func NewAuthService(cfg *config.Config) *AuthService {
	return &AuthService{Config: cfg}
}

// Login authenticates a user and generates JWT tokens.
func (s *AuthService) Login(request models.LoginRequest) (models.LoginResponse, error) {
	const HardcodedUsername = "admin"
	const HardcodedPassword = "password123"

	if request.Username != HardcodedUsername || request.Password != HardcodedPassword {
		return models.LoginResponse{}, errors.New("invalid credentials")
	}

	accessToken, refreshToken, err := utils.GenerateTokens(request.Username, s.Config)
	if err != nil {
		return models.LoginResponse{}, err
	}

	return models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// RefreshToken generates a new access token using a refresh token.
func (s *AuthService) RefreshToken(request models.RefreshTokenRequest) (models.LoginResponse, error) {
	claims, err := utils.ValidateToken(request.RefreshToken, true, s.Config)
	if err != nil {
		return models.LoginResponse{}, err
	}

	username, ok := claims["username"].(string)
	if !ok {
		return models.LoginResponse{}, errors.New("invalid token payload")
	}

	accessToken, refreshToken, err := utils.GenerateTokens(username, s.Config)
	if err != nil {
		return models.LoginResponse{}, err
	}

	return models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// ProvideAuthService is for fx dependency injection.
var ProvideAuthService = fx.Provide(NewAuthService)
