package interfaces

import (
	"ecommerce_clean_architecture/pkg/domain"
	"ecommerce_clean_architecture/pkg/utils/models"
)

type OrderRepository interface {
	DoesCartExist(userID int) (bool, error)
	AddressExist(orderBody models.OrderIncoming) (bool, error)
	CreateOrder(orderDetails domain.Order) error
	AddOrderItems(orderItemDetails domain.OrderItem, userID int, ProductID int, Quantity float64) error
	GetBriefOrderDetails(orderID string) (domain.OrderSuccessResponse, error)
}
