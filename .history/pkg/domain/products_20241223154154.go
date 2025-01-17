package domain

import (
	"time"

	"gorm.io/gorm"
)

type Products struct {
	ID         int      `json:"id" gorm:"primarykey;not null"`
	CategoryID int      `json:"category_id"`
	Category   Category `json:"-" gorm:"foreignkey:CategoryID;constraint:OnDelete:CASCADE"`
	Name       string   `json:"name" validate:"required"`
	Stock      int      `json:"stock"`
	Quantity   int      `json:"quantity"`
	Price      float64  `json:"price"`
	CreatedAt  time.Time
	DeletedAt  gorm.DeletedAt
}
