package interfaces

import "ecommerce_clean_architecture/pkg/utils/models"

type PaymentUsecase interface {
	CreatePayment(orderID string, userID int) (models.CombinedOrderDetails, string, error)
	OnlinePaymentVerification(details models.OnlinePaymentVerification) (*[]models.CombinedOrderDetails, error)
}
