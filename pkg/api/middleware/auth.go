package middleware

import (
	"ecommerce_clean_architecture/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header is missing"})
			c.Abort()
			return
		}

		tokenString := helper.GetTokenFromHeader(authHeader)

		userID, email, err := helper.ExtractUserIDFromToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token", "error": err.Error()})
			c.Abort()
			return
		}

		c.Set("id", userID)
		c.Set("email", email)
		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header is missing"})
			c.Abort()
			return
		}

		const bearerPrefix = "Bearer "
		if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header format is invalid"})
			c.Abort()
			return
		}

		token := authHeader[len(bearerPrefix):]

		name, email, role, ID, err := helper.ExtractClaimsFromToken(token)
		if err != nil || role != "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized or invalid token", "error": err.Error()})
			c.Abort()
			return
		}

		c.Set("id", ID)
		c.Set("admin_name", name)
		c.Set("admin_email", email)
		c.Next()
	}
}
