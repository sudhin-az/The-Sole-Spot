package repository

import (
	"ecommerce_clean_architecture/pkg/domain"
	"ecommerce_clean_architecture/pkg/utils/models"
	"fmt"

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
func (o *OrderRepository) BeginTransaction() (*gorm.DB, error) {
	tx := o.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (o *OrderRepository) CommitTransaction(tx *gorm.DB) error {
	return tx.Commit().Error
}

func (o *OrderRepository) RollbackTransaction(tx *gorm.DB) error {
	return tx.Rollback().Error
}

func (o *OrderRepository) DoesCartExist(userID int) (bool, error) {
	var exist bool
	err := o.DB.Raw("select exists(select * from carts where user_id = ?)", userID).Scan(&exist).Error
	fmt.Println("Exist", exist)
	if err != nil {
		return false, err
	}
	return exist, nil
}

func (o *OrderRepository) AddressExist(AddressID int) (bool, error) {
	var count int
	err := o.DB.Raw("select count(*)from addresses where id=?", AddressID).Scan(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (o *OrderRepository) GetProductStock(productID int) (int, error) {
	var stock int
	err := o.DB.Raw("select stock from products where id = ?", productID).Scan(&stock).Error
	if err != nil {
		return 0, err
	}
	return stock, nil
}
func (o *OrderRepository) UpdateProductStock(tx *gorm.DB, productID int, newStock int) error {
	return tx.Model(&domain.Products{}).Where("id = ?", productID).Update("stock", newStock).Error
}

func (o *OrderRepository) CreateOrder(tx *gorm.DB, orderDetails models.Order) (int, error) {
	result := tx.Create(&orderDetails)
	if result.Error != nil {
		return 0, result.Error
	}
	return orderDetails.OrderId, nil
}
func (o *OrderRepository) CreateOrderItems(tx *gorm.DB, orderItems []domain.OrderItem) error {
	result := tx.Create(&orderItems)
	return result.Error
}

func (o *OrderRepository) FetchCartItem(userID int) ([]domain.Cart, error) {
	var cartItems []domain.Cart
	err := o.DB.Raw("select * from carts where user_id = ? and deleted_at is null", userID).Scan(&cartItems).Error
	if err != nil {
		return []domain.Cart{}, err
	}
	return cartItems, nil
}
func (o *OrderRepository) GetBriefOrderDetails(orderID int) (models.Order, error) {
	var orderSuccessResponse models.Order

	err := o.DB.Raw("select * from orders where order_id = ?", orderID).Scan(&orderSuccessResponse).Error
	if err != nil {
		return models.Order{}, err
	}
	return orderSuccessResponse, nil
}

func (o *OrderRepository) UserOrderRelationship(orderID int, userID int) (int, error) {
	var testUserID int
	err := o.DB.Raw("SELECT user_id FROM orders WHERE order_id = ?", orderID).Scan(&testUserID).Error
	if err != nil {
		return -1, err
	}
	return testUserID, nil
}

func (o *OrderRepository) GetProductDetailsFromOrders(orderID int) ([]models.OrderProducts, error) {
	var orderProductDetails []models.OrderProducts
	err := o.DB.Raw("SELECT product_id, quantity FROM order_items WHERE order_id = ?", orderID).Scan(&orderProductDetails).Error
	if err != nil {
		return []models.OrderProducts{}, err
	}
	return orderProductDetails, nil
}

func (o *OrderRepository) GetShipmentStatus(orderID int) (string, error) {
	var shipmentStatus string
	err := o.DB.Raw("SELECT shipment_status FROM orders WHERE order_id = ?", orderID).Scan(&shipmentStatus).Error
	if err != nil {
		return "", err
	}
	return shipmentStatus, nil
}

func (o *OrderRepository) GetOrderDetails(userID int) ([]models.FullOrderDetails, error) {
	var orderDetails []models.OrderDetails
	o.DB.Raw("SELECT order_id, final_price, shipment_status, payment_status FROM orders WHERE user_id = ?", userID).Scan(&orderDetails)
	fmt.Println(orderDetails)

	var fullOrderDetails []models.FullOrderDetails
	for _, od := range orderDetails {
		var orderProductDetails []models.OrderProductDetails
		o.DB.Raw(`
			SELECT 
				order_items.product_id, 
				products.name AS product_name, 
				order_items.quantity, 
				order_items.total_price 
			FROM order_items 
			INNER JOIN products ON order_items.product_id = products.id 
			WHERE order_items.order_id = ?`, od.OrderId).Scan(&orderProductDetails)
		fullOrderDetails = append(fullOrderDetails, models.FullOrderDetails{
			OrderDetails:        od,
			OrderProductDetails: orderProductDetails,
		})
	}
	return fullOrderDetails, nil
}

func (o *OrderRepository) CancelOrders(orderID int) error {
	shipmentStatus := "cancelled"
	err := o.DB.Exec("UPDATE orders SET shipment_status = ? WHERE order_id = ?", shipmentStatus, orderID).Error
	if err != nil {
		return err
	}
	var paymentMethod int
	err = o.DB.Raw("SELECT payment_method_id FROM orders WHERE order_id = ?", orderID).Scan(&paymentMethod).Error
	if err != nil {
		return err
	}
	if paymentMethod == 1 || paymentMethod == 3 {
		err = o.DB.Exec("UPDATE orders SET payment_status = 'refunded' WHERE order_id = ?", orderID).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *OrderRepository) UpdateQuantityOfProduct(orderProducts []models.OrderProducts) error {
	for _, od := range orderProducts {
		var quantity int
		err := o.DB.Raw("select quantity from products where id = ?", od.ProductID).Scan(&quantity).Error
		if err != nil {
			return err
		}
		od.Quantity += quantity
		if err := o.DB.Exec("update products set quantity = ? where id = ?", od.Quantity, od.ProductID).Error; err != nil {
			return err
		}
	}
	return nil
}
