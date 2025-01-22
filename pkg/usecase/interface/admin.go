package interfaces

import (
	"ecommerce_clean_architecture/pkg/domain"
	"ecommerce_clean_architecture/pkg/utils/models"
)

type AdminUseCase interface {
	SignUpHandler(admin models.AdminSignUp) (domain.TokenAdmin, error)
	LoginHandler(admin models.AdminLogin) (domain.TokenAdmin, error)
	GetUsers() ([]models.User, error)
	BlockUser(userID int) error
	UnBlockUsers(userID int) error
	GetAllOrderDetails() ([]models.FullOrderDetails, error)
	CancelOrders(orderID string, userID int) error
	ChangeOrderStatus(orderID string, Status string) (models.Order, error)
}
