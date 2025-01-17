package interfaces

import (
	"ecommerce_clean_architecture/pkg/domain"
	"ecommerce_clean_architecture/pkg/utils/models"

	"gorm.io/gorm"
)

type AdminRepository interface {
	CheckAdminAvailability(admin models.AdminSignUp) bool
	SignUpHandler(admin models.AdminSignUp) (models.AdminDetailsResponse, error)
	LoginHandler(admin models.AdminLogin) (domain.AdminDetails, error)
	GetUsers() ([]models.User, error)
	GetUserByID(userID int) (models.User, error)
	UpdateBlockUserByID(user models.User) error
	BeginTransaction() (*gorm.DB, error)
	CommitTransaction(tx *gorm.DB) error
	RollbackTransaction(tx *gorm.DB) error
}
