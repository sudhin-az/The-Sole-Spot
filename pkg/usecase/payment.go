package usecase

import (
	"ecommerce_clean_arch/pkg/helper"
	"ecommerce_clean_arch/pkg/repository"
	"ecommerce_clean_arch/pkg/utils/models"
	"errors"
	"fmt"
	"log"

	"github.com/razorpay/razorpay-go"
)

type PaymentUsecase struct {
	PaymentRepo repository.PaymentRepository
}

func NewPaymentUsecase(paymentRepo repository.PaymentRepository) *PaymentUsecase {
	return &PaymentUsecase{PaymentRepo: paymentRepo}
}

func (pay *PaymentUsecase) CreatePayment(orderID string, userID int) (models.CombinedOrderDetails, string, error) {

	combinedOrderDetails, err := pay.PaymentRepo.GetOrderDetailsByOrderId(orderID)
	if err != nil {
		return models.CombinedOrderDetails{}, "", fmt.Errorf("failed to fetch order details: %v", err)
	}

	client := razorpay.NewClient("rzp_test_X9t6tBXI0YHLgE", "YleZMSEez6iOgZYm2RSluiit")
	if client == nil {
		return models.CombinedOrderDetails{}, "", errors.New("failed to create Razorpay client")
	}

	data := map[string]interface{}{
		"amount":   int(combinedOrderDetails.FinalPrice) * 100,
		"currency": "INR",
		"receipt":  "some_receipt_id",
	}

	body, err := client.Order.Create(data, nil)
	if err != nil {
		return models.CombinedOrderDetails{}, "", fmt.Errorf("razorpay order creation failed: %v", err)
	}

	razorPayOrderID, ok := body["id"].(string)
	if !ok {
		return models.CombinedOrderDetails{}, "", errors.New("failed to retrieve Razorpay order ID")
	}

	err = pay.PaymentRepo.AddRazorPayDetails(orderID, razorPayOrderID)
	if err != nil {
		return models.CombinedOrderDetails{}, "", fmt.Errorf("failed to store Razorpay details: %v", err)
	}

	return combinedOrderDetails, razorPayOrderID, nil
}

func (pay *PaymentUsecase) OnlinePaymentVerification(details models.OnlinePaymentVerification) (*[]models.CombinedOrderDetails, error) {
	status, err := pay.PaymentRepo.CheckPaymentStatus(details.OrderID)
	if err != nil {
		return nil, err
	}
	if status == "failed" {
		return nil, errors.New("payment failed in razorpay")
	}
	if status == "paid" {
		return nil, errors.New("Already paid")
	}
	if status == "not paid" {
		err := pay.PaymentRepo.UpdatePaymentDetails(details.OrderID, details.PaymentID)
		if err != nil {
			return nil, err
		}
	}
	result := helper.VerifyPayment(details.RazorPayOrderID, details.PaymentID, details.Signature, "YleZMSEez6iOgZYm2RSluiit")
	if !result {
		return nil, errors.New("payment is unsuccessful")
	}
	orders, err := pay.PaymentRepo.UpdateOnlinePaymentSucess(details.OrderID)
	log.Println("OrderID is recieved: ", details.OrderID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}
