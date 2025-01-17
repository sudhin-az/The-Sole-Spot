package middleware

import (
	"ecommerce_clean_architecture/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header is missing"})
			c.Abort()
			return
		}

		// Get the token from the Authorization header
		tokenString := helper.GetTokenFromHeader(authHeader)

		// Verify the token and extract the user ID
		userID, email, err := helper.ExtractUserIDFromToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token", "error": err.Error()})
			c.Abort()
			return
		}

		// Set the user ID in the context
		c.Set("id", userID)
		c.Set("email", email)
		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header is missing"})
			c.Abort()
			return
		}

		name, email, role, err := helper.ExtractClaimsFromToken(token)
		if err != nil || role != "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized or invalid token", "error": err.Error()})
			c.Abort()
			return
		}

		// Set admin details in the context
		c.Set("admin_name", name)
		c.Set("admin_email", email)
		c.Next()
	}
}
