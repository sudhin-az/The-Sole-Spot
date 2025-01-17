package models

import "time"

type Order struct {
	UserID          int `json:"user_id"`
	AddressID       int `json:"address_id"`
	PaymentMethodID int `json:"payment_method_id"`
	CouponID        int `json:"coupon_id"`
}
type OrderFromCart struct {
	PaymentID int `json:"payment_id" binding:"required"`
	AddressID int `json:"address_id" binding:"required"`
}
type OrderIncoming struct {
	UserID    int `json:"user_id"`
	PaymentID int `json:"payment_id"`
	AddressID int `json:"address_id"`
}
type OrderResponse struct {
	UserID         int       `json:"user_id"`
	OrderID        int       `json:"order_id"`
	Quantity       int       `json:"quantity"`
	DiscountAmount float64   `json:"discount_amount"`
	Total          float64   `json:"total"`
	Method         string    `json:"method"`
	Status         string    `json:"status"`
	PaymentStatus  string    `json:"payment_status"`
	OrderDate      time.Time `json:"order_date"`
}
type OrderDetails struct {
	OrderId    string
	FinalPrice float64
	Status     string
}
type OrderProductDetails struct {
	ProductID  int     `json:"product_id"`
	Name       string  `json:"name"`
	Quantity   int     `json:"quantity"`
	TotalPrice float64 `json:"total_price"`
}
type FullOrderDetails struct {
	OrderDetails        OrderDetails
	OrderProductDetails []OrderProductDetails
}
type OrderProducts struct {
	ProductId string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}
