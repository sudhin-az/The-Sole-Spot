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
}

func ExtractUserIDFromToken(tokenString string) (int, string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &authCustomClaimsUsers{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return []byte("132457689"), nil
	})

	if err != nil {
		return 0, "", fmt.Errorf("error parsing token: %v", err)
	}

	claims, ok := token.Claims.(*authCustomClaimsUsers)
	if !ok {
		return 0, "", fmt.Errorf("invalid token claims")
	}

	return claims.Id, claims.Email, nil
}

func ExtractClaimsFromToken(tokenString string) (name string, email string, role string, userID int, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("12345678"), nil
	})

	if err != nil {
		return "", "", "", 0, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return "", "", "", 0, fmt.Errorf("invalid token")
	}

	return claims.Name, claims.Email, claims.Role, claims.UserID, nil
}
