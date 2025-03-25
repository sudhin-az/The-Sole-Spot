package interfaces

import (
	"ecommerce_clean_arch/pkg/domain"
	"ecommerce_clean_arch/pkg/utils/models"
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

	GetDateRange(startDate, endDate, limit string) (string, string)
	TotalOrders(fromDate, toDate, PaymentStatus string) (models.OrderCount, models.AmountInformation, error)

	BestSellingProduct() ([]models.BestSellingProduct, error)
	BestSellingCategory() ([]models.BestSellingCategory, error)
}
