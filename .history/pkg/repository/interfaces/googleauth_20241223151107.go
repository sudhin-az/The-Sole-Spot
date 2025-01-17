package interfaces

import "ecommerce_clean_architecture/pkg/utils/models"

type AuthRepository interface {
	GetUserByEmail(email string) (models.User, error)
	CreateUser(newuser models.User) error
}
