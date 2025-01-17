package models

import "time"

type Order struct {
	OrderId         string     `json:"order_id" gorm:"primaryKey;not null"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at" gorm:"index"`
	DeliveryTime    time.Time  `json:"delivery_time"`
	UserID          int        `json:"user_id" gorm:"not null"`
	AddressID       uint       `json:"address_id"`
	Address         Address    `json:"-" gorm:"foreignkey:AddressID"`
	PaymentMethodID uint       `json:"paymentmethod_id"`
	PaymentMethod   string     `json:"-" gorm:"foreignkey:PaymentMethodID"`
	GrandTotal      float64    `json:"grand_total"`
	FinalPrice      float64    `json:"discount_price"`
	ShipmentStatus  string     `json:"status"`
	PaymentStatus   string     `json:"payment_status"`
	Approval        bool       `json:"approval"`
}
type OrderWithItem struct {
	OrderID         int64   `json:"order_id"`
	UserID          int64   `json:"user_id"`
	AddressID       int64   `json:"address_id"`
	CouponID        *int64  `json:"coupon_id"`
	Discount        float64 `json:"discount"`
	OrderTotalPrice float64 `json:"order_total_price"`
	GrandTotal      float64 `json:"grand_total"`
	OrderStatus     string  `json:"order_status"`
	PaymentStatus   string  `json:"payment_status"`
	PaymentMethod   string  `json:"payment_method"`
	OrderDate       string  `json:"order_date"`
	ShipmentStatus  string  `json:"shipment_status"`
	ProductID       int64   `json:"product_id"`
	ItemQuantity    int64   `json:"item_quantity"`
	ItemTotalPrice  float64 `json:"item_total_price"`
}
type OrderFromCart struct {
	PaymentID uint `json:"payment_id" binding:"required"`
	AddressID uint `json:"address_id" binding:"required"`
}

type OrderIncoming struct {
	UserID    uint `json:"user_id"`
	PaymentID uint `json:"payment_id"`
	AddressID uint `json:"address_id"`
}

type OrderProducts struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type OrderDetails struct {
	OrderId        string
	FinalPrice     float64
	ShipmentStatus string
	PaymentStatus  string
}

type OrderProductDetails struct {
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	TotalPrice  float64 `json:"total_price"`
}

type FullOrderDetails struct {
	OrderDetails        OrderDetails
	OrderProductDetails []OrderProductDetails
}

type CombinedOrderDetails struct {
	OrderId        string  `json:"order_id"`
	FinalPrice     float64 `json:"final_price"`
	ShipmentStatus string  `json:"shipment_status"`
	PaymentStatus  string  `json:"payment_status"`
	Name           string  `json:"firstname"`
	Email          string  `json:"email"`
	Phone          string  `json:"phone"`
	HouseName      string  `json:"house_name" validate:"required"`
	State          string  `json:"state" validate:"required"`
	Pin            string  `json:"pin" validate:"required"`
	Street         string  `json:"street"`
	City           string  `json:"city"`
}
