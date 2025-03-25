package interfaces

import "ecommerce_clean_arch/pkg/utils/models"

type PaymentRepository interface {
	AddRazorPayDetails(orderID string, razorPayOrderID string) error
	GetOrderDetailsByOrderId(orderID string) (models.CombinedOrderDetails, error)
	CheckPaymentStatus(orderID string) (string, error)
	UpdateOnlinePaymentSucess(orderID string) (*[]models.CombinedOrderDetails, error)
	UpdatePaymentDetails(orderID string, paymentID string) error
}
