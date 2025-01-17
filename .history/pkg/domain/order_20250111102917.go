package domain

import "time"

type Order struct {
	OrderId         int        `json:"order_id" gorm:"primaryKey;not null"`
	UserID          int        `json:"user_id" gorm:"not null"`
	AddressID       uint       `json:"address_id"`
	Address         Address    `json:"-" gorm:"foreignkey:AddressID"`
	CouponID        *int       `json:"coupon_id"`
	Discount        float64    `json:"discount"`
	GrandTotal      float64    `json:"grand_total"`
	Status          string     `json:"status"`
	Method          *string    `json:"method"`
	PaymentStatus   string     `json:"payment_status"`
	PaymentMethodID uint       `json:"paymentmethod_id"`
	OrderDate       time.Time  `json:"order_date"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at" gorm:"index"`
	DeliveryTime    *time.Time `json:"delivery_time"`
	ShipmentStatus  string     `json:"shipment_status"`
	Approval        bool       `json:"approval"`
	FinalPrice      float64    `json:"discount_price"`
	PaymentMethod   string     `json:"-" gorm:"foreignkey:PaymentMethodID"`
}
type OrderItem struct {
	ID         int      `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderID    int      `json:"order_id"`
	Orders     Order    `json:"-" gorm:"foreignkey:OrderID;constraint:OnDelete:CASCADE"`
	ProductID  int      `json:"product_id"`
	Products   Products `json:"-" gorm:"foreignkey:ProductID"`
	Quantity   int      `json:"quantity"`
	TotalPrice float64  `json:"total_price"`
}

type OrderSuccessResponse struct {
	OrderID        string `json:"order_id"`
	ShipmentStatus string `json:"order_status"`
}
