package handlers

import (
	"net/http"

	"itv-task/internal/models"
	"itv-task/internal/services"
	"itv-task/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	AuthService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{AuthService: authService}
}

// Login godoc
// @Summary Login user
// @Description Authenticate user with hardcoded credentials
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var request models.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid request", "Failed to parse request body")
		return
	}

	response, err := h.AuthService.Login(request)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "Invalid credentials", err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Generate a new access token using the refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var request models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid request", "Failed to parse request body")
		return
	}

	response, err := h.AuthService.RefreshToken(request)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "Invalid or expired refresh token", err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}
