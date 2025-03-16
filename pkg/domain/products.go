package domain

import (
	"time"

	"gorm.io/gorm"
)

type Products struct {
	ID         int      `json:"id" gorm:"primaryKey;not null"`
	CategoryID int      `json:"category_id" gorm:"column:category_id"`             
	Category   Category `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE"` 
	Name       string   `json:"name" validate:"required"`
	Stock      int      `json:"stock"`
	Quantity   int      `json:"quantity"`
	Price      float64  `json:"price"`
	OfferPrice float64  `json:"offer_price"`
	CreatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}
