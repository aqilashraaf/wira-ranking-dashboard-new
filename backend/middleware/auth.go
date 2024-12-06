package middleware

import (
	"net/http"
	"strings"
	"wira-dashboard/utils"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware verifies the JWT token in the Authorization header
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract the token from the Authorization header
		// Format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		claims, err := utils.ValidateJWT(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Store user information in the context
		c.Set("user_id", int((*claims)["user_id"].(float64)))
		c.Set("username", (*claims)["username"].(string))
		c.Next()
	}
}

// Optional2FA middleware checks if 2FA is required for the user
func Optional2FA() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement 2FA check when needed
		c.Next()
	}
}
