package models

import "time"

type Order struct {
	OrderId          int        `gorm:"primaryKey;autoIncrement" json:"order_id"`
	UserID           int        `json:"user_id" gorm:"not null"`
	AddressID        uint       `json:"address_id"`
	Address          Address    `json:"-" gorm:"foreignkey:AddressID"`
	CouponID         *uint      `json:"coupon_id,omitempty"`
	CouponCode       string     `json:"coupon_code"`
	RawTotal         float64    `json:"raw_total"`
	Discount         float64    `json:"discount"`
	DiscountAmount   float64    `json:"discount_amount"`
	CategoryDiscount float64    `json:"category_discount"`
	GrandTotal       float64    `json:"grand_total"`
	DeliveryCharge   float64    `json:"delivery_charge"`
	PaymentStatus    string     `json:"payment_status"`
	PaymentMethodID  uint       `json:"paymentmethod_id"`
	OrderDate        time.Time  `json:"order_date"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at" gorm:"index"`
	OrderStatus      string     `gorm:"column:order_status"`
	FinalPrice       float64    `json:"final_price"`
	PaymentMethod    string     `json:"-" gorm:"foreignkey:PaymentMethodID"`
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
	ProductID  int     `json:"product_id"`
	Quantity   int     `json:"quantity"`
	FinalPrice float64 `json:"total_price"`
}

type OrdersDetails struct {
	CustomerName        string
	CustomerPhoneNumber string
	CustomerAddress     Address
	CustomerCity        string
	OrderDate           time.Time
	Items               []InvoiceItem
	OrderStatus         string
	CategoryDiscount    float64
	GrandTotal          float64
	RawAmount           float64
	FinalPrice          float64
	Discount            float64
	DeliveryCharge      float64
}

type OrderProductDetails struct {
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	TotalPrice  float64 `json:"total_price"`
}

type OrderDetails struct {
	OrderId          int     `json:"order_id"`
	DiscountAmount   float64 `json:"discount_amount"`
	CategoryDiscount float64 `json:"category_discount"`
	GrandTotal       float64 `json:"grand_total"`
	FinalPrice       float64 `json:"final_price"`
	OrderStatus      string  `json:"order_status"`
	PaymentStatus    string  `json:"payment_status"`
}

type FullOrderDetails struct {
	OrderDetails        OrderDetails          `json:"OrderDetails"`
	OrderProductDetails []OrderProductDetails `json:"OrderProductDetails"`
}

type CombinedOrderDetails struct {
	OrderId       string  `json:"order_id"`
	FinalPrice    float64 `json:"final_price"`
	OrderStatus   string  `json:"order_status"`
	PaymentStatus string  `json:"payment_status"`
	Name          string  `json:"first_name"`
	Email         string  `json:"email"`
	Phone         string  `json:"phone"`
	HouseName     string  `json:"house_name" validate:"required"`
	State         string  `json:"state" validate:"required"`
	District      string  `json:"district" validate:"required"`
	Pin           string  `json:"pin" validate:"required"`
	Street        string  `json:"street"`
	City          string  `json:"city"`
}

type OrderCount struct {
	TotalOrder     uint `json:"total_order"`
	TotalPending   uint `json:"total_pending"`
	TotalConfirmed uint `json:"total_confirmed"`
	TotalShipped   uint `json:"total_shipped"`
	TotalDelivered uint `json:"total_delivered"`
	TotalCancelled uint `json:"total_cancelled"`
	TotalReturned  uint `json:"total_returned"`
}

type AmountInformation struct {
	TotalAmountBeforeDeduction float64 `json:"total_amount_before_deduction"`
	TotalCouponDeduction       float64 `json:"total_coupon_deduction"`
	TotalProuctOfferDeduction  float64 `json:"total_product_offer_deduction"`
	TotalAmountAfterDeduction  float64 `json:"total_amount_after_deduction"`
}

type InvoiceItem struct {
	Name     string  `json:"name"`
	Quantity uint    `json:"quantity"`
	Price    float64 `json:"price"`
}
