package models

type PaymentMethod struct {
	ID           uint   `gorm:"primarykey"`
	Payment_Name string `json:"payment_name"`
}

type RazorPay struct {
	ID        int    `json:"id" gorm:"primarykey not null"`
	OrderID   string `json:"order_id"`
	RazorID   string `json:"razor_id"`
	PaymentID string `json:"payment_id"`
}

type OnlinePaymentVerification struct {
	PaymentID       string `json:"payment_id" validate:"required"`
	OrderID         int    `json:"order_id" validate:"required"`
	RazorPayOrderID string `json:"razorpay_order_id" validate:"required`
	Signature       string `json:"signature" validate:"required"`
}
