package repository

import (
	"gorm.io/gorm"
)

type OrderRepository struct {
	DB *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{
		DB: db,
	}
}

func (o *OrderRepository) DoesCartExist(userID int) (bool, error) {
	var exist bool
	err := o.DB.Raw("select exists(select 1 from carts where user_id = ?)", userID).Scan(&exist).Error
	if err != nil {
		return false, err
	}
	return exist, nil
}
