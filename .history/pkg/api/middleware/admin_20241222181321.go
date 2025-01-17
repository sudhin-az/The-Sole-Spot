package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AuthorizationMiddleware(c *gin.Context) {

	//	"Bearer" is a type of token authentication.
	//It is used to indicate that the token being sent is a bearer token,
	//meaning that whoever possesses the token can access the associated resources without further authentication
	//geting the token from the authorization header
	s := c.Request.Header.Get("Authorization")

	var token string

	//check the authorization is startinfg with "bearer"
	if s[:7] == "Bearer " {
		// If it does, extract the token part after "Bearer "
		token = strings.TrimPrefix(s, "Bearer ")
	} else {
		token = s
	}

	if err := validateToken(token); err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

}

func validateToken(token string) error {
	// Parse the JWT token
	_, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		// Return the key used for validation
		return []byte("12345678"), nil
	})

	return err

}
