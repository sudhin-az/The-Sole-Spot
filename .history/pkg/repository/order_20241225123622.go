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

func (o *OrderRepository) GetOrderDetailsByOrderId(orderID string) (models.CombinedOrderDetails, error) {
	var orderDetails models.CombinedOrderDetails

	err := o.DB.Raw("select orders.order_id, orders.final_price, orders.shipment_status, orders.payment_status, users.name, users.email, users.phone, addresses.house_name, addresses.street, addresses.city, addresses.district, addresses.state, addresses.pin from orders join users on orders.user_id = users.id join addresses on orders.address_id = addresses.id where orders.order_id= ?", orderID).Scan(&orderDetails).Error
	if err != nil {
		return models.CombinedOrderDetails{}, err
	}
	return orderDetails, nil
}

func (o *OrderRepository) GetOrders(orderID string) (domain.Order, error) {
	var body domain.Order

	if err := o.DB.Raw("select * from orders where order_id = ?", orderID).Scan(&body).Error; err != nil {
		return domain.Order{}, err
	}
	return body, nil
}

func (o *OrderRepository) UserOrderRelationship(orderID string, userID string) (int, error) {
	var testUserID int
	err := o.DB.Raw("select user_id from orders where order_id = ?", orderID).Scan(&testUserID).Error
	if err != nil {
		return -1, err
	}
	return testUserID, nil
}

func (o *OrderRepository) GetProductDetailsFromOrders(orderID string) ([]models.OrderProducts, error) {
	var orderProductDetails []models.OrderProducts
	err := o.DB.Raw("select product_id, quantity from order_items where order_id = ?", orderID).Scan(&orderProductDetails).Error
	if err != nil {
		return []models.OrderProducts{}, err
	}
	return orderProductDetails, nil
}

func (o *OrderRepository) GetShipmentStatus(ordeID string) (string, error) {
	var shipmentStatus string
	err := o.DB.Raw("select shipment_status from orders where order_id = ?", ordeID).Scan(&shipmentStatus).Error
	if err != nil {
		return "", err
	}
	return shipmentStatus, nil
}
