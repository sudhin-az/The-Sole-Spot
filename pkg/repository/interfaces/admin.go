package interfaces

import (
	"ecommerce_clean_architecture/pkg/domain"
	"ecommerce_clean_architecture/pkg/utils/models"
)

type AdminRepository interface {
	CheckAdminAvailability(admin models.AdminSignUp) bool
	SignUpHandler(admin models.AdminSignUp) (models.AdminDetailsResponse, error)
	LoginHandler(admin models.AdminLogin) (domain.AdminDetails, error)
	GetUsers(listusers models.UserSignUp) (models.UserSignUp, error)
}
