package domain

import "time"

type Order struct {
	OrderID         int     `json:"order_id" gorm:"primaryKey;not null"`
	UserID          int     `json:"user_id" gorm:"not null"`
	Users           Users   `json:"-" gorm:"foreignkey:UserID"`
	AddressID       int     `json:"address_id" gorm:"not null"`
	Address         Address `json:"-" gorm:"foreignkey:AddressID"`
	CouponID        int     `gorm:"index"`
	Discount        float64
	Quantity        int       `gorm:"default:0"`
	Status          string    `gorm:"check(status IN('Pending', 'Shipped', 'Delivered', 'Canceled','Failed','Returned'))"`
	Method          string    `gorm:"check(method IN('Credit Card', 'PayPal', 'Bank Transfer'))"`
	PaymentStatus   string    `gorm:"check(payment_status IN('Processing', 'Success', 'Failed'))"`
	OrderDate       time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	PaymentMethodID int
	PaymentMethod   PaymentMethod `gorm:"foreignkey:PaymentMethodID;references:ID"`
	TotalPrice      int           `json:"total_price"`
	CreatedAt       time.Time     `gorm:"autoCreateTime"`
	UpdatedAt       time.Time     `gorm:"autoUpdateTime"`
}

type UserOrderItem struct {
	ID         int      `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderID    string   `json:"order_id"`
	Order      Order    `json:"-" gorm:"foreignkey:OrderID;constraint:OnDelete:CASCADE"`
	Product_ID int      `json:"product_id"`
	Products   Products `json:"-" gorm:"foreignkey:ProductID"`
	Quantity   int      `json:"quantity"`
	TotalPrice int      `json:"total_price"`
}
type PaymentMethod struct {
	ID           int    `json:"id" gorm:"primarykey"`
	Payment_Name string `json:"payment_name"`
}
type OrderSuccessResponse struct {
	OrderID        string `json:"order_id"`
	ShipmentStatus string `json:"order_status"`
}
