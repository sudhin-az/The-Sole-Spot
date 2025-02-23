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
	GetWalletAmount(userID int) (float64, error)
	UpdateWalletAmount(walletAmount float64, UserID int) error
	CancelOrders(orderID string) error
	GetOrderStatus(orderID string) (string, error)
	GetPaymentStatus(orderID string) (string, error)
	GetPriceoftheproduct(orderID string) (float64, error)
	GetProductDetailsFromOrders(orderID string) ([]models.OrderProducts, error)
	UpdateQuantityOfProduct(orderProducts []models.OrderProducts) error
	CancelOrderItem(orderItemID string, userID int) (domain.OrderItem, error)
	ReturnUserOrder(orderID string, userID int) error
	GetOrderItemPrice(tx *gorm.DB, orderItemID int) (float64, error)
	GetOrderItemDetails(tx *gorm.DB, orderItemID int) (int, int, error)

	GetCouponDetails(couponCode string) (models.Coupon, error)
	CheckCouponUsage(userID uint, couponCode string) (int, error)
	RecordCouponUsage(tx *gorm.DB, userID int, couponCode string) error
	CheckCouponAppliedOrNot(userID int, couponID string) uint
}
