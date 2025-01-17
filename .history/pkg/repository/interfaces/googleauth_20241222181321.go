package interfaces

import "ecommerce_clean_architecture/pkg/domain"

type AuthRepository interface {
    GetUserByEmail(email string) (domain.Users, error)
    CreateUser(user domain.Users) error
}
