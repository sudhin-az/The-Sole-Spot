package domain

import (
	"time"

	"gorm.io/gorm"
)

type Cart struct {
	ID         int            `json:"id"`
	UserID     int            `json:"user_id"`
	ProductID  int            `json:"product_id"`
	Quantity   int            `json:"quantity"`
	Price      float64        `json:"price"`
	OfferPrice float64        `json:"offer_price"`
	TotalPrice float64        `json:"total_price"`
	CreatedAt  time.Time      `json:"created_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
