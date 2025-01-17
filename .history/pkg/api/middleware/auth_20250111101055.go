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

<<<<<<< HEAD
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token format"})
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := helper.VerifyAccessToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token"})
=======
		// Get the token from the Authorization header
		tokenString := helper.GetTokenFromHeader(authHeader)

		// Verify the token and extract the user ID
		userID, email, err := helper.ExtractUserIDFromToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token", "error": err.Error()})
>>>>>>> 9c98fd5 (order management)
			c.Abort()
			return
		}

<<<<<<< HEAD
		id, ok := claims["id"].(float64) // Use float64 for numeric claims
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token payload"})
			c.Abort()
			return
		}
		c.Set("id", int(id)) // Save ID as an int
=======
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
		name, email, role, userID, err := helper.ExtractClaimsFromToken(token) // Extract userID along with name, email, role
		if err != nil || role != "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized or invalid token", "error": err.Error()})
			c.Abort()
			return
		}

		// Set admin details and user_id in the context
		c.Set("user_id", userID)
		c.Set("admin_name", name)
		c.Set("admin_email", email)
>>>>>>> 9c98fd5 (order management)
		c.Next()
	}
}
