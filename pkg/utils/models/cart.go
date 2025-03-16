package models

import (
	"time"

	"gorm.io/gorm"
)

type Cart struct {
	ID               int            `json:"id"`
	UserID           int            `json:"user_id"`
	ProductID        int            `json:"product_id"`
	Quantity         int            `json:"quantity"`
	Price            int            `json:"price"`
	OfferPrice       int            `json:"offer_price"`
	CategoryDiscount float64        `json:"category_discount"`
	TotalPrice       float64        `json:"total_price"`
	CreatedAt        time.Time      `json:"created_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
type CartResponse struct {
	TotalPrice float64 `json:"total_price"`
	Cart       []Cart  `json:"cart"`
}
