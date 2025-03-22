package helper

import (
	"log"

	"github.com/razorpay/razorpay-go/utils"
)

func VerifyPayment(RazorPayOrderID string, paymentID, signature, razopaySecret string) bool {

	params := map[string]interface{}{
		"razorpay_order_id":   RazorPayOrderID,
		"razorpay_payment_id": paymentID,
	}

	result := utils.VerifyPaymentSignature(params, signature, razopaySecret)
	log.Println("*****", result)
	return result
}
