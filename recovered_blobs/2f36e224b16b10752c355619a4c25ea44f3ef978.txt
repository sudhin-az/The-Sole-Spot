package helper

import (
	"ecommerce_clean_architecture/pkg/utils/models"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type authCustomClaimsUsers struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	jwt.StandardClaims
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func GenerateTokenUsers(userID uint, userEmail string, expirationTime time.Time) (string, error) {

	claims := &authCustomClaimsUsers{
		Id:    int(userID), // Convert to int for use in claims
		Email: userEmail,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("132457689"))
	fmt.Println("errrrrrrrrr", tokenString)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateAccessToken(user models.UserDetailsResponse) (string, error) {

	expirationTime := time.Now().Add(2 * time.Hour)
	tokenString, err := GenerateTokenUsers(uint(user.Id), user.Email, expirationTime)
	if err != nil {
		return "", err
	}
	return tokenString, nil

}

func GenerateRefreshToken(user models.UserDetailsResponse) (string, error) {

	expirationTime := time.Now().Add(24 * 90 * time.Hour)
	tokeString, err := GenerateTokenUsers(uint(user.Id), user.Email, expirationTime)
	if err != nil {
		return "", err
	}
	return tokeString, nil
}
func VerifyAccessToken(token string) (map[string]interface{}, error) {
	// Parse and validate the token (example using JWT)
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method and return the secret key
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("132457689"), nil // Match the signing key used during token generation
	})

	if err != nil || !parsedToken.Valid {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}
