package models

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
	TotalPrice float64        `json:"total_price"`
	CreatedAt  time.Time      `json:"created_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

type CartResponse struct {
	TotalPrice float64 `json:"total_price"`
	Cart       []Cart  `json:"cart"`
}

type CartTotal struct {
	UserName   string  `json:"user_name"`
	TotalPrice float64 `json:"total_price"`
	FinalPrice float64 `json:"final_price"`
}
