package models

import "time"

type Order struct {
	OrderId         int        `json:"order_id" gorm:"primaryKey;not null"`
	UserID          int        `json:"user_id" gorm:"not null"`
	AddressID       uint       `json:"address_id"`
	Address         Address    `json:"-" gorm:"foreignkey:AddressID"`
	CouponID        *int       `json:"coupon_id"`
	Discount        float64    `json:"discount"`
	GrandTotal      float64    `json:"grand_total"`
	Method          string     `json:"method"`
	PaymentStatus   string     `json:"payment_status"`
	PaymentMethodID uint       `json:"paymentmethod_id"`
	OrderDate       time.Time  `json:"order_date"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at" gorm:"index"`
	DeliveryTime    *time.Time `json:"delivery_time"`
	OrderStatus     string     `json:"order_status"`
	Approval        bool       `json:"approval"`
	FinalPrice      float64    `json:"final_price"`
	PaymentMethod   string     `json:"-" gorm:"foreignkey:PaymentMethodID"`
}
type OrderFromCart struct {
	PaymentID uint    `json:"payment_id" binding:"required"`
	AddressID uint    `json:"address_id" binding:"required"`
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
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
	OrderId        int
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
