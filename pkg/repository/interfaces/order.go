package interfaces

import (
	"ecommerce_clean_architecture/pkg/domain"
	"ecommerce_clean_architecture/pkg/utils/models"

	"gorm.io/gorm"
)

type OrderRepository interface {
	BeginTransaction() (*gorm.DB, error)
	CommitTransaction(tx *gorm.DB) error
	RollbackTransaction(tx *gorm.DB) error
	DoesCartExist(orderBody models.OrderFromCart, userID int) (bool, error)
	AddressExist(orderBody models.OrderIncoming) (bool, error)
	GetProductStock(ProductID int) (int, error)
	UpdateProductStock(tx *gorm.DB, productID int, newStock int) error
	CreateOrder(tx *gorm.DB, orderDetails models.OrderFromCart) error
	CreateOrderItems(tx *gorm.DB, orderItems []domain.OrderItem) error
	GetBriefOrderDetails(orderID string) (domain.OrderSuccessResponse, error)
	UserOrderRelationship(orderID string, userID int) (int, error)
	GetOrderDetails(userID int) ([]models.FullOrderDetails, error)
	CancelOrders(orderID string) error
	GetOrderStatus(orderID string) (string, error)
	GetProductDetailsFromOrders(orderID string) ([]models.OrderProducts, error)
	UpdateQuantityOfProduct(orderProducts []models.OrderProducts) error
}
