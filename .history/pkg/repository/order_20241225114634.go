package repository

import (
	"ecommerce_clean_architecture/pkg/domain"
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

func (o *OrderRepository) CreateOrder(orderDetails domain.Order) error {
	err := o.DB.Create(&orderDetails).Error
	if err != nil {
		return err
	}
	return nil
}

func (o *OrderRepository) AddOrderItems(orderItemDetails domain.OrderItem, userID int, ProductID int, Quantity float64) error {

	// after creating the order delete all cart items and also update the quantity of the product
	err := o.DB.Omit("id").Create(&orderItemDetails).Error
	if err != nil {
		return err
	}

	err = o.DB.Exec("delete from carts where user_id = ? and product_id = ?", userID, ProductID).Error
	if err != nil {
		return err
	}

	err = o.DB.Exec("update products set quantity = quantity - ? where id = ?", Quantity, ProductID).Error
	if err != nil {
		return err
	}
	return nil
}
func (o *OrderRepository) GetBriefOrderDetails(orderID string) (domain.OrderSuccessResponse, error) {
	var orderSuccessResponse domain.OrderSuccessResponse

	err := o.DB.Raw("select order_id, shipment_status from orders where order_id = ?", orderID).Scan(&orderSuccessResponse).Error
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}
	return orderSuccessResponse, nil
}

func (o *OrderRepository) GetOrderDetailsByOrderId(orderID string) (domain.) {
	
}
