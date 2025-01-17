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
func (o *OrderRepository) FetchOrderWithItems(orderID int64) (models.OrderWithItem, error) {
	var order models.OrderWithItem
	err := o.DB.Table("orders").
		Select("orders.*, order_items.product_id, order_items.quantity, order_items.total_price").
		Joins("JOIN order_items ON orders.order_id = order_items.order_id").
		Where("orders.order_id = ?", orderID).
		Scan(&order).Error
	return order, err
}
func (o *OrderRepository) GetOrderWithItems(orderID int64) ([]models.OrderWithItem, error) {
	var results []models.OrderWithItem
	err := o.DB.Table("order_with_items").
		Where("order_id = ?", orderID).
		Find(&results).Error
	return results, err
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
func (o *OrderRepository) GetBriefOrderDetails(orderID int) (models.Order, error) {
	var orderSuccessResponse models.Order

	err := o.DB.Raw("select * from orders where order_id = ?", orderID).Scan(&orderSuccessResponse).Error
	if err != nil {
		return models.Order{}, err
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

func (o *OrderRepository) UserOrderRelationship(orderID string, userID int) (int, error) {
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

func (o *OrderRepository) GetOrderDetails(userID int) ([]models.FullOrderDetails, error) {
	var orderDetails []models.OrderDetails

	if err := o.DB.Raw("SELECT order_id, final_price, shipment_status, payment_status FROM orders where user_id = ?", userID).Scan(&orderDetails).Error; err != nil {
		return nil, err
	}
	var fullOrderDetails []models.FullOrderDetails
	for _, op := range orderDetails {
		var orderProductDetails []models.OrderProductDetails

		if err := o.DB.Raw(`SELECT order_items.product_id, products.product_name, order_items.quantity, order_items.total_price FROM order_items
		INNER JOIN products ON order_items.product_id = products.id 
		WHERE order_items.order_id = ?`, op.OrderId).Scan(&orderProductDetails).Error; err != nil {
			return nil, err
		}
		fullOrderDetails = append(fullOrderDetails, models.FullOrderDetails{OrderDetails: op, OrderProductDetails: orderProductDetails})
	}
	return fullOrderDetails, nil
}

func (o *OrderRepository) CancelOrders(orderID string) error {
	shipmentStatus := "cancelled"
	err := o.DB.Exec("update orders set shipment_status = ? where order_id = ?", shipmentStatus, orderID).Error
	if err != nil {
		return err
	}
	var paymentMethod int
	err = o.DB.Raw("select payment_method_id from orders where order_id = ?", orderID).Scan(&paymentMethod).Error
	if err != nil {
		return err
	}
	if paymentMethod == 1 || paymentMethod == 3 {
		err = o.DB.Exec("update orders set payment_status = 'refunded' where order_id = ?", orderID).Error
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

func (o *OrderRepository) GetOrderDetailsBrief(page int) ([]models.CombinedOrderDetails, error) {
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * 2
	var orderDetails []models.CombinedOrderDetails

	err := o.DB.Raw("select orders.order_id, orders.final_price, orders.shipment_status, orders.payment_status, users.name,users.email, users.phone,addresses.house_name,addresses.state, addresses.pin, addresses.street,addresses.city from orders inner join users on orders.user_id = users.id inner join addresses on users.id = addresses.user_id limit ? offset ?", 2, offset).Scan(&orderDetails).Error
	if err != nil {
		return []models.CombinedOrderDetails{}, err
	}
	return orderDetails, nil
}
func (o *OrderRepository) GetPaymentStatus(orderID string) (string, error) {
	var paymentStatus string
	err := o.DB.Raw("select payment_status from orders where order_id = ?", orderID).Scan(&paymentStatus).Error
	if err != nil {
		return "", err
	}
	return paymentStatus, nil
}
func (o *OrderRepository) GetPriceOftheProduct(orderID string) (float64, error) {
	var price float64
	err := o.DB.Raw("select grand_total from orders where order_id = ?", orderID).Scan(&price).Error
	if err != nil {
		return 0.0, err
	}
	return price, nil
}
func (o *OrderRepository) CheckOrderID(orderID string) (bool, error) {
	var count int
	err := o.DB.Raw("select count(*) from orders where order_id = ?", orderID).Scan(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (o *OrderRepository) ApproveOrder(orderID string) error {

	err := o.DB.Exec("update orders set shipment_status = 'order placed',approval=true where order_id = ?", orderID).Error
	if err != nil {
		return err
	}
	return nil
}

func (o *OrderRepository) GetOrderDetailsOfProduct(orderID string) (models.OrderDetails, error) {

	var orderDetails models.OrderDetails
	err := o.DB.Raw("select final_price, shipment_status, payment_status where order_id = ?", orderID).Scan(&orderDetails).Error
	if err != nil {
		return models.OrderDetails{}, err
	}
	return orderDetails, nil
}

func (o *OrderRepository) GetAddressDetailsFromID(orderID string) (models.Address, error) {

	var address_id int
	var addresses models.Address

	err := o.DB.Raw("select address_id from orders where order_id = ?", orderID).Scan(&addresses).Error
	if err != nil {
		return models.Address{}, err
	}
	err = o.DB.Raw("select * from addresses where id = ?", address_id).Scan(&addresses).Error
	if err != nil {
		return models.Address{}, err
	}
	return models.Address{}, nil
}
func (o *OrderRepository) GetOrderDetailsByID(orderID string) (models.CombinedOrderDetails, error) {
	var orders models.CombinedOrderDetails

	query := `select orders.user_id,users.name,users.email,users.phone,addresses.house_name,addresses.street,addresses.city,addresses.state,addresses.pin,orders.address_id,orders.payment_method_id,payment_methods.payment_name,orders.final_price from orders inner join users on orders.user_id=users.id inner join addresses on orders.address_id=addresses.id inner join payment_methods on orders.payment_method_id=payment_methods.id where orders.order_id=?  `

	err := o.DB.Raw(query, orderID).Scan(&orders).Error
	if err != nil {
		return models.CombinedOrderDetails{}, err
	}
	return orders, nil
}
func (o *OrderRepository) GetItemsByOrderId(orderID string) ([]models.ProductDetails, error) {
	var items []models.ProductDetails

	query := "select products.name,order_items.quantity,products.price,order_items.total_price from order_items  join products on order_items.product_id=products.id where order_items.order_id=?"

	if err := o.DB.Raw(query, orderID).Scan(&items).Error; err != nil {
		return []models.ProductDetails{}, err
	}
	return items, nil

}
