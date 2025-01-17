package interfaces

import "ecommerce_clean_architecture/pkg/utils/models"

type OrderRepository interface {
	DoesCartExist(userID int) (bool, error)
	AddressExist(orderBody models.OrderIncoming) (bool, error)
}
