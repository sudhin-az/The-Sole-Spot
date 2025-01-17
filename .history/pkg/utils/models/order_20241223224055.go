package models

import "time"

type Order struct {
	OrderID         int       `json:"order_id" gorm:"column:order_id;primaryKey;autoIncrement"`
	UserID          int       `json:"user_id" gorm:"column:user_id"`
	AddressID       int       `json:"address_id" gorm:"column:address_id"`
	CouponID        *int      `json:"coupon_id,omitempty" gorm:"column:coupon_id"`
	Discount        float64   `json:"discount" gorm:"column:discount"`
	Quantity        int       `json:"quantity" gorm:"column:quantity"`
	Status          string    `gorm:"check(status IN('Pending', 'Shipped', 'Delivered', 'Canceled','Failed','Returned'))"`
	Method          string    `gorm:"check(method IN('Credit Card', 'PayPal', 'Bank Transfer'))"`
	PaymentStatus   string    `gorm:"check(payment_status IN('Processing', 'Success', 'Failed'))"`
	OrderDate       time.Time `json:"order_date" gorm:"column:order_date"`
	PaymentMethodID *int      `json:"payment_method_id,omitempty" gorm:"column:payment_method_id"`
	TotalPrice      float64   `json:"total_price" gorm:"column:total_price"`
	CreatedAt       time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"column:updated_at"`
}
