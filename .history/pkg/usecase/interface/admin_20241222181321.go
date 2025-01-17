package interfaces

import (
	"ecommerce_clean_architecture/pkg/domain"
	"ecommerce_clean_architecture/pkg/utils/models"
)

type AdminUseCase interface {
	SignUpHandler(admin models.AdminSignUp) (domain.TokenAdmin, error)
	LoginHandler(admin models.AdminLogin) (domain.TokenAdmin, error)
	GetUsers(listusers models.UserSignUp) (models.UserSignUp, error)
}
