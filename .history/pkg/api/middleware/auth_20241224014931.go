package middleware

import (
	"ecommerce_clean_architecture/pkg/helper"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

package middleware

import (
	"ecommerce_clean_architecture/pkg/helper"
	"net/http"
	"strings"

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

		// Split and validate the Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token format"})
			c.Abort()
			return
		}
		token := parts[1]

		// Verify the token
		claims, err := helper.VerifyAccessToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Extract user ID from claims and set it in context
		userID, ok := claims["id"].(float64) // JWT unmarshals numbers as float64
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token payload"})
			c.Abort()
			return
		}
		c.Set("id", int(userID)) // Convert float64 to int
		c.Next()
	}
}

