package interfaces

import (
	"ecommerce_clean_arch/pkg/domain"
	"ecommerce_clean_arch/pkg/utils/models"

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
	GetAllOrderDetails() ([]models.FullOrderDetails, error)
	GetProductStockk(ProductID int) (int, error)
	UpdateProductStockk(tx *gorm.DB, productID int, newStock int) error
	AdminOrderRelationship(orderID string, userID int) (int, error)
	GetProductDetailFromOrders(orderID string) ([]models.OrderProducts, error)
	Cancelorders(orderID string) error
	Getshipmentstatus(orderID string) (string, error)
	UpdatequantityOfproduct(orderProducts []models.OrderProducts) error
	ChangeOrderStatus(orderID string, Status string) (models.Order, error)

	GetTotalOrders(fromDate, toDate, PaymentStatus string) (models.OrderCount, models.AmountInformation, error)

	BestSellingProduct() ([]models.BestSellingProduct, error)
	BestSellingCategory() ([]models.BestSellingCategory, error)
}
