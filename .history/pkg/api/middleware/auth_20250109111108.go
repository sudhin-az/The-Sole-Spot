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
		c.Set("id", ID)
		c.Set("id", userID)
		c.Set("email", email)
		c.Next()
	}
}
