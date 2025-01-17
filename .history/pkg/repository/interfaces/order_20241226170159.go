package interfaces

import (
	"ecommerce_clean_architecture/pkg/domain"
	"ecommerce_clean_architecture/pkg/utils/models"
)

type OrderRepository interface {
	DoesCartExist(userID int) (bool, error)
	AddressExist(orderBody models.OrderIncoming) (bool, error)
	GetProductStock(ProductID int) (int, error)
	UpdateProductStock(productID int, newStock int) error
	CreateOrder(orderDetails domain.Order) error
	AddOrderItems(orderItemDetails domain.OrderItem, UserID int, ProductID uint, Quantity float64) error
	GetBriefOrderDetails(orderID string) (domain.OrderSuccessResponse, error)
	GetOrderDetailsByOrderId(orderID string) (models.CombinedOrderDetails, error)
	GetOrders(orderID string) (domain.Order, error)
	UserOrderRelationship(orderID string, userID int) (int, error)
	GetOrderDetails(userID int) ([]models.FullOrderDetails, error)
	CancelOrders(orderID string) error
	GetShipmentStatus(orderID string) (string, error)
	GetProductDetailsFromOrders(orderID string) ([]models.OrderProducts, error)
	UpdateQuantityOfProduct(orderProducts []models.OrderProducts) error
	GetOrderDetailsBrief(page int) ([]models.CombinedOrderDetails, error)
	GetPaymentStatus(orderID string) (string, error)
	GetPriceoftheproduct(orderID string) (float64, error)
	CheckOrderID(orderID string) (bool, error)
	ApproveOrder(orderID string) error
	GetOrderDetailsofAproduct(orderID string) (models.OrderDetails, error)
	GetAddressDetailsFromID(orderID string) (models.Address, error)
	GetItemsByOrderId(orderID string) ([]models.ProductDetails, error)
	GetOrderDetailsByID(orderID string) (models.CombinedOrderDetails, error)
}
