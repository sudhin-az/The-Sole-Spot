package domain

import "time"

type Order struct {
	OrderId         int        `gorm:"primaryKey;autoIncrement" json:"order_id"`
	UserID          int        `json:"user_id" gorm:"not null"`
	Users           Users      `json:"-" gorm:"foreignkey:UserID"`
	AddressID       uint       `json:"address_id"`
	Address         Address    `json:"-" gorm:"foreignkey:AddressID"`
	CouponID        int        `json:"coupon_id"`
	Coupon          Coupons    `json:"-" gorm:"foreignkey:CouponID"`
	CouponCode      string     `json:"coupon_code"`
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
type OrderItem struct {
	ID         int      `gorm:"primaryKey;autoIncrement"`
	OrderID    int      `json:"order_id"`
	Orders     Order    `json:"-" gorm:"foreignkey:OrderID;constraint:OnDelete:CASCADE"`
	ProductID  int      `json:"product_id"`
	Products   Products `json:"-" gorm:"foreignkey:ProductID"`
	Quantity   int      `json:"quantity"`
	TotalPrice float64  `json:"total_price"`
}

type OrderSuccessResponse struct {
	OrderID     string `json:"order_id"`
	OrderStatus string `json:"order_status"`
}
