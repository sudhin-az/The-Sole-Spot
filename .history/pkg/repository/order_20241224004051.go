package repository

import (
	"ecommerce_clean_architecture/pkg/repository/interfaces"

	"gorm.io/gorm"
)

type OrderRepository struct {
	DB *gorm.DB
}

func NewOrderRepository(db *gorm.DB) interfaces.OrderRepository {
	return &OrderRepository{
		DB: db,
	}
}
