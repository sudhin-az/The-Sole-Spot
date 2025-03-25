package interfaces

import (
	"ecommerce_clean_arch/pkg/domain"
	"ecommerce_clean_arch/pkg/utils/models"

	"github.com/jung-kurt/gofpdf"
)

type OrderUseCase interface {
	OrderItemsFromCart(orderFromCart models.OrderFromCart, userID int) (domain.OrderSuccessResponse, error)
	GetOrderDetails(userID int, page int, count int) ([]models.FullOrderDetails, error)
	CancelOrders(orderID string, userID int) error
	CancelOrderItem(orderItemID string, userID int) (domain.OrderItem, error)
	ReturnUserOrder(orderID string, userID int) error
	GenerateInvoice(orderID string, userID int) (*gofpdf.Fpdf, error)
}
