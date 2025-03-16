package domain

import "time"

type Order struct {
	OrderId          int        `gorm:"primaryKey;autoIncrement" json:"order_id"`
	UserID           int        `json:"user_id" gorm:"not null"`
	Users            Users      `json:"-" gorm:"foreignkey:UserID"`
	AddressID        uint       `json:"address_id"`
	Address          Address    `json:"-" gorm:"foreignkey:AddressID"`
	CouponID         int        `json:"coupon_id"`
	Coupon           Coupons    `json:"-" gorm:"foreignkey:CouponID"`
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
	OrderStatus      string     `json:"order_status"`
	FinalPrice       float64    `json:"final_price"`
	PaymentMethod    string     `json:"-" gorm:"foreignkey:PaymentMethodID"`
}
type OrderItem struct {
	ID         int      `gorm:"primaryKey;autoIncrement"`
	OrderID    int      `json:"order_id" gorm:"column:order_id"`
	Order      Order    `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
	ProductID  int      `json:"product_id" gorm:"column:product_id"`
	Product    Products `gorm:"foreignKey:ProductID"`
	Quantity   int      `json:"quantity"`
	TotalPrice float64  `json:"total_price"`
}
type OrderSuccessResponse struct {
	OrderID     string `json:"order_id"`
	OrderStatus string `json:"order_status"`
}
