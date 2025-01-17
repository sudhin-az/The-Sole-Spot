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

func ExtractUserIDFromToken(tokenString string) (int, string, error) {
	// Parse the token and validate the claims
	token, err := jwt.ParseWithClaims(tokenString, &authCustomClaimsUsers{}, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return []byte("132457689"), nil // Replace with your actual secret key
	})

	// Check if there's an error parsing the token
	if err != nil {
		return 0, "", fmt.Errorf("error parsing token: %v", err)
	}

	// Extract claims from the parsed token
	claims, ok := token.Claims.(*authCustomClaimsUsers)
	if !ok {
		return 0, "", fmt.Errorf("invalid token claims")
	}

	// Return user ID and email from claims
	return claims.Id, claims.Email, nil
}
