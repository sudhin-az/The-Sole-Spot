package models

import (
	"time"

	"gorm.io/gorm"
)

type AddProduct struct {
	ID         int     `gorm:"id"`
	CategoryID int     `json:"category_id"`
	Name       string  `json:"name"`
	Stock      int     `json:"stock"`
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`
	CreatedAt  time.Time
	DeletedAt  gorm.DeletedAt
}
type ProductResponse struct {
	ID          int     `json:"id" `
	Category_Id int     `json:"category_id"`
	Name        string  `json:"name" `
	Stock       int     `json:"stock"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	CreatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}
type SearchItems struct {
	Name string `json:"name" binding:"required"`
}
type ProductDetails struct {
	Name       string  `json:"name"`
	TotalPrice float64 `json:"total_price"`
	Price      float64 `json:"price" `
	Total      float64 `json:"total"`
	Quantity   int     `json:"quantity"`
}
