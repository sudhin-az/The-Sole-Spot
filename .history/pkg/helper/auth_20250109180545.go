package helper

import (
	"fmt"

	"github.com/golang-jwt/jwt"
)

func GetTokenFromHeader(header string) string {

	if len(header) > 7 && header[:7] == "Bearer " {
		return header[7:]
	}

	return header
	// return ""
}

func ExtractClaimsFromToken(tokenString string) (string, string, string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("12345678"), nil // Replace with your secret key
	})
	if err != nil {
		return "", "", "", err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return "", "", "", fmt.Errorf("invalid token")
	}

	return claims.Name, claims.Email, claims.Role, nil
}
