package domain

import "time"

type Order struct {
	OrderID         int        `json:"order_id"`
	UserID          int        `json:"user_id"`
	AddressID       int        `json:"address_id"`
	CouponID        *int       `json:"coupon_id"`
	Discount        float64    `json:"discount"`
	GrandTotal      float64    `json:"grand_total"`
	Status          string     `json:"status"`
	Method          *string    `json:"method"`
	PaymentStatus   string     `json:"payment_status"`
	PaymentMethodID int        `json:"payment_method_id"`
	OrderDate       time.Time  `json:"order_date"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at"`
	DeliveryTime    *time.Time `json:"delivery_time"`
	ShipmentStatus  string     `json:"shipment_status"`
	Approval        bool       `json:"approval"`
	FinalPrice      float64    `json:"final_price"`
	PaymentMethod   string     `json:"payment_method"`
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
