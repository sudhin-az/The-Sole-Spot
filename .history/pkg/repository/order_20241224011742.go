package repository

import (
	"gorm.io/gorm"
)

type OrderRepository struct {
	DB *gorm.DB
}

func NewOrderRepository(db *gorm.DB) repository.OrderRepository {
	return &OrderRepository{
		DB: db,
	}
}
