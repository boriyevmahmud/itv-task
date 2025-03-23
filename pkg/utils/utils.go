package utils

import (
	"errors"
	"time"

	"itv-task/config"
	"itv-task/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateTokens(username string, cfg *config.Config) (string, string, error) {
	accessTokenClaims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(config.AccessTokenTTL).Unix(),
		"iat":      time.Now().Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString([]byte(cfg.JWTAccessSecret))
	if err != nil {
		return "", "", err
	}

	refreshTokenClaims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(config.AccessTokenTTL).Unix(),
		"iat":      time.Now().Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(cfg.JWTRefreshSecret))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func ValidateToken(tokenString string, isRefresh bool, cfg *config.Config) (jwt.MapClaims, error) {
	// Choose correct secret
	secret := []byte(cfg.JWTAccessSecret)
	if isRefresh {
		secret = []byte(cfg.JWTRefreshSecret)
	}

	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token expired")
		}
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

func SendErrorResponse(c *gin.Context, code int, message string, detail string) {
	c.JSON(code, models.NewErrorResponse(code, message, detail))
}
