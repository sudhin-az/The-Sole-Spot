package helper

import (
	"ecommerce_clean_architecture/pkg/utils/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type CustomClaims struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	UserID int    `json:"user_id"`
	jwt.StandardClaims
}

func GenerateTokenAdmin(admin models.AdminDetailsResponse) (string, error) {
	claims := &CustomClaims{
		Name:  admin.Name,
		Email: admin.Email,
		Role:  "admin",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte("12345678"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
