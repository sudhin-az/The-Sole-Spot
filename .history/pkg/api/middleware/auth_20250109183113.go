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
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header is missing"})
			c.Abort()
			return
		}

		// Check if the header starts with "Bearer "
		const bearerPrefix = "Bearer "
		if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header format is invalid"})
			c.Abort()
			return
		}

		// Extract the token part
		token := authHeader[len(bearerPrefix):]

		// Extract claims from the token
		name, email, role, err := helper.ExtractClaimsFromToken(token)
		if err != nil || role != "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized or invalid token", "error": err.Error()})
			c.Abort()
			return
		}

		// Set admin details in the context
		c.Set("user_id", userID) // Assuming userID is an integer obtained from the toke
		c.Set("admin_email", email)
		c.Next()
	}
}
