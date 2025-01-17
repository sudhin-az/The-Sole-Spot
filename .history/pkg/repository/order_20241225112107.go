package repository

import (
	"ecommerce_clean_architecture/pkg/utils/models"

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

func (o *OrderRepository) AddressExist(orderBody models.OrderIncoming) (bool, error) {
	var count int
	err := o.DB.Raw("select count(*) from addresses where user_id = ? and id = ?", orderBody.UserID, orderBody.AddressID).Scan(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (o *OrderRepository) CreateOrder(orderDetails domain.) {
	
}
