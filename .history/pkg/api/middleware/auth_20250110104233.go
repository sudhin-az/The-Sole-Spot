package middleware

import (
	"ecommerce_clean_architecture/pkg/helper"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	ContextUserIDKey  = "user_id"
	ContextEmailKey   = "email"
	ContextAdminName  = "admin_name"
	ContextAdminEmail = "admin_email"
)

func extractBearerToken(authHeader string) (string, error) {
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("invalid Authorization header format")
	}
	return authHeader[len("Bearer "):], nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header is missing"})
			c.Abort()
			return
		}

		token, err := extractBearerToken(authHeader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			c.Abort()
			return
		}

		userID, email, err := helper.ExtractUserIDFromToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token", "error": err.Error()})
			c.Abort()
			return
		}

		c.Set(ContextUserIDKey, userID)
		c.Set(ContextEmailKey, email)
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

		token, err := extractBearerToken(authHeader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			c.Abort()
			return
		}

		name, email, role, userID, err := helper.ExtractClaimsFromToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized or invalid token", "error": err.Error()})
			c.Abort()
			return
		}

		if role != "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized - only admins are allowed"})
			c.Abort()
			return
		}

		c.Set(ContextUserIDKey, userID)
		c.Set(ContextAdminName, name)
		c.Set(ContextAdminEmail, email)
		c.Next()
	}
}
