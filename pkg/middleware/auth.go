package utils

import (
	"itv-task/config"
	"itv-task/pkg/utils"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware checks the validity of the access token
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			utils.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Missing Authorization header")
			c.Abort()
			return
		}

		// Extract Bearer token
		parts := strings.Split(token, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Invalid token format")
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(parts[1], false, cfg) // Change `nil` if you need to pass config
		if err != nil {
			utils.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Invalid or expired token")
			c.Abort()
			return
		}

		c.Set("user", claims)
		c.Next()
	}
}

// AuthLogger logs requests for debugging purposes
func AuthLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("📢 Incoming request: %s %s", c.Request.Method, c.Request.URL.Path)
		c.Next()
	}
}
