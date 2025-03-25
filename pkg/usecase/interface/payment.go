package interfaces

import "ecommerce_clean_arch/pkg/utils/models"

type PaymentUsecase interface {
	CreatePayment(orderID string, userID int) (models.CombinedOrderDetails, string, error)
	OnlinePaymentVerification(details models.OnlinePaymentVerification) (*[]models.CombinedOrderDetails, error)
}
