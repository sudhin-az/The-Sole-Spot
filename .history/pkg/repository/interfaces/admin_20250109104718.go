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
	BeginnTransaction() (*gorm.DB, error)
	CommittTransaction(tx *gorm.DB) error
	RollbackkTransaction(tx *gorm.DB) error
	GetOrderDetails(userID int) ([]models.FullOrderDetails, error)
	GetProductStock(ProductID int) (int, error)
	UpdateProductStock(tx *gorm.DB, productID int, newStock int) error
	AdminOrderRelationship(orderID string, userID int) (int, error)
	GetProductDetailsFromOrders(orderID string) ([]models.OrderProducts, error)
	CancelOrders(orderID string) error
	GetShipmentStatus(orderID string) (string, error)
	UpdateQuantityOfProduct(orderProducts []models.OrderProducts) error
}
