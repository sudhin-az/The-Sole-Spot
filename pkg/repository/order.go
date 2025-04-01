package repository

import (
	"ecommerce_clean_arch/pkg/domain"
	"ecommerce_clean_arch/pkg/utils/models"
	"errors"
	"log"

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
	log.Println("Exist", exist)
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

func (o *OrderRepository) GetProductStock(tx *gorm.DB, productID int) (int, error) {
	var stock int
	err := tx.Raw("select stock from products where id = ?", productID).Scan(&stock).Error
	if err != nil {
		return 0, err
	}
	return stock, nil
}
func (o *OrderRepository) UpdateProductStock(tx *gorm.DB, productID int, newStock int) error {
	return tx.Model(&domain.Products{}).Where("id = ?", productID).Update("stock", newStock).Error
}

func (o *OrderRepository) CreateOrder(tx *gorm.DB, orderDetails models.Order) (int, error) {
	if orderDetails.CouponID != nil {
		var count int64
		if err := tx.Table("coupons").Where("id = ?", orderDetails.CouponID).Count(&count).Error; err != nil {
			return 0, err
		}
		if count == 0 {
			return 0, errors.New("invalid coupon id")
		}
	}
	result := tx.Omit("OrderId").Create(&orderDetails)
	if result.Error != nil {
		return 0, result.Error
	}
	return orderDetails.OrderId, nil
}
func (o *OrderRepository) CreateOrderItems(tx *gorm.DB, orderItems []domain.OrderItem) error {
	for i := range orderItems {
		orderItems[i].ID = 0
	}
	result := tx.Omit("ID").Create(&orderItems)
	return result.Error
}

func (o *OrderRepository) FetchCartItem(tx *gorm.DB, userID int) ([]domain.Cart, error) {
	var cartItems []domain.Cart
	err := tx.Raw("select * from carts where user_id = ? and deleted_at is null", userID).Scan(&cartItems).Error
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

func (o *OrderRepository) GetProductDetailsFromOrders(tx *gorm.DB, orderID int) ([]models.OrderProducts, error) {
	var orderProductDetails []models.OrderProducts
	err := tx.Raw("select * from order_items where order_id = ?", orderID).Scan(&orderProductDetails).Error
	if err != nil {
		return nil, err
	}
	var FinalPrice []float64
	err = tx.Raw("select final_price from orders where order_id = ?", orderID).Scan(&FinalPrice).Error
	if err != nil {
		return []models.OrderProducts{}, err
	}
	for i, val := range orderProductDetails {
		orderProductDetails[i].FinalPrice = FinalPrice[i]
		log.Println("HASHIM________________________", val.FinalPrice, FinalPrice[i])
	}
	log.Println("final____________-price", FinalPrice, orderProductDetails)
	return orderProductDetails, nil
}

func (o *OrderRepository) GetOrderItemPrice(tx *gorm.DB, orderItemID int) (float64, error) {
	var price float64
	err := tx.Raw("select total_price from order_items where id = ?", orderItemID).Scan(&price).Error
	if err != nil {
		return 0.0, err
	}
	return price, nil
}

func (o *OrderRepository) GetOrderItemDetails(tx *gorm.DB, orderItemID int) (int, int, error) {
	var orderItem struct {
		productID int
		quantity  int
	}
	err := tx.Raw("SELECT product_id, quantity FROM order_items WHERE id = ?", orderItemID).Scan(&orderItem).Error
	if err != nil {
		return 0, 0, err
	}
	return orderItem.productID, orderItem.quantity, nil
}

func (o *OrderRepository) GetOrderStatus(tx *gorm.DB, orderID int) (string, error) {
	var OrderStatus string
	err := tx.Raw("SELECT order_status FROM orders WHERE order_id = ?", orderID).Scan(&OrderStatus).Error
	if err != nil {
		return "", err
	}
	return OrderStatus, nil
}
func (o *OrderRepository) GetPaymentStatus(tx *gorm.DB, orderID string) (string, error) {
	var paymentstatus string
	err := tx.Raw("select payment_status from orders where order_id=?", orderID).Scan(&paymentstatus).Error
	if err != nil {
		return "", err
	}
	return paymentstatus, nil
}
func (o *OrderRepository) UpdatePaymentStatus(tx *gorm.DB, orderID int, paymentStatus string) error {
	return tx.Exec("UPDATE orders SET payment_status = ? WHERE order_id = ?", paymentStatus, orderID).Error
}

func (o *OrderRepository) GetPriceoftheproduct(tx *gorm.DB, orderID string) (float64, error) {
	var a float64
	err := tx.Raw("select final_price from orders where order_id=?", orderID).Scan(&a).Error
	if err != nil {
		return 0.0, err
	}
	return a, nil
}
func (o *OrderRepository) FetchOrderDetailsFromDB(orderID string) (models.OrdersDetails, error) {
	var order models.Order
	if err := o.DB.Where("order_id = ?", orderID).First(&order).Error; err != nil {
		return models.OrdersDetails{}, err
	}

	var user models.User
	if err := o.DB.Where("id = ?", order.UserID).First(&user).Error; err != nil {
		return models.OrdersDetails{}, err
	}

	var address models.Address
	if err := o.DB.Where("id = ?", order.AddressID).First(&address).Error; err != nil {
		return models.OrdersDetails{}, err
	}

	var orderItems []domain.OrderItem
	if err := o.DB.Where("order_id = ?", orderID).Find(&orderItems).Error; err != nil {
		return models.OrdersDetails{}, err
	}
	var products []domain.Products
	var RawTotal float64
	for _, item := range orderItems {
		var product domain.Products
		if err := o.DB.Model(&product).Where("id = ?", item.ProductID).First(&product).Error; err != nil {
			return models.OrdersDetails{}, nil
		}
		products = append(products, product)

		for _, product := range products {
			if item.Product.ID == product.ID {
				RawTotal += product.Price * float64(item.Quantity)
			}
		}
	}
	log.Println("rawTotal", RawTotal)
	var items []models.InvoiceItem
	for _, orderItem := range orderItems {
		var product domain.Products
		if err := o.DB.Model(&product).Where("id = ?", orderItem.ProductID).First(&product).Error; err != nil {
			return models.OrdersDetails{}, err
		}
		items = append(items, models.InvoiceItem{
			Name:     product.Name,
			Quantity: uint(orderItem.Quantity),
			Price:    product.Price,
		})
	}

	orderDetails := models.OrdersDetails{
		CustomerName:        user.FirstName,
		CustomerPhoneNumber: user.Phone,
		CustomerAddress:     address,
		CustomerCity:        address.City,
		OrderDate:           order.OrderDate,
		Items:               items,
		OrderStatus:         order.OrderStatus,
		GrandTotal:          order.GrandTotal,
		CategoryDiscount:    order.CategoryDiscount,
		RawAmount:           order.RawTotal,
		FinalPrice:          order.FinalPrice,
		Discount:            order.DiscountAmount,
		DeliveryCharge:      order.DeliveryCharge,
	}
	log.Println("orderDetails", orderDetails)
	return orderDetails, nil
}
func (o *OrderRepository) GetOrderDetails(userID int) ([]models.FullOrderDetails, error) {
	var orderDetails []models.OrderDetails
	o.DB.Raw("SELECT order_id, discount_amount, category_discount, grand_total, final_price, order_status, payment_status FROM orders WHERE user_id = ?", userID).Scan(&orderDetails)
	log.Println(orderDetails)

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

func (o *OrderRepository) GetWalletAmount(tx *gorm.DB, userID int) (float64, error) {
	var walletAvailable float64
	err := tx.Raw("select balance from wallets where user_id = ?", userID).Scan(&walletAvailable).Error
	if err != nil {
		return 0.0, err
	}
	return walletAvailable, nil
}

func (o *OrderRepository) UpdateWalletAmount(tx *gorm.DB, walletAmount float64, UserID int) error {
	return tx.Exec("UPDATE wallets SET balance = ? WHERE user_id = ?", walletAmount, UserID).Error
}

func (o *OrderRepository) CancelOrders(tx *gorm.DB, orderID int) error {
	OrderStatus := "cancelled"
	err := tx.Exec("UPDATE orders SET order_status = ? WHERE order_id = ?", OrderStatus, orderID).Error
	if err != nil {
		return err
	}
	var paymentMethod int
	err = tx.Raw("SELECT payment_method_id FROM orders WHERE order_id = ?", orderID).Scan(&paymentMethod).Error
	if err != nil {
		return err
	}
	if paymentMethod == 1 || paymentMethod == 3 {
		err = tx.Exec("UPDATE orders SET payment_status = 'refunded' WHERE order_id = ?", orderID).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *OrderRepository) UpdateQuantityOfProduct(tx *gorm.DB, orderProducts []models.OrderProducts) error {
	for _, od := range orderProducts {
		var quantity int
		err := tx.Raw("select quantity from products where id = ?", od.ProductID).Scan(&quantity).Error
		if err != nil {
			return err
		}
		od.Quantity += quantity
		if err := tx.Exec("update products set quantity = ? where id = ?", od.Quantity, od.ProductID).Error; err != nil {
			return err
		}
	}
	return nil
}

func (o *OrderRepository) CancelOrderItem(tx *gorm.DB, orderItemID int) error {
	result := tx.Exec("DELETE FROM order_items WHERE id = ?", orderItemID)
	if result.Error != nil {
		return errors.New("error cancelling order item")
	}
	return nil
}

func (o *OrderRepository) UpdateUserOrderReturn(tx *gorm.DB, orderID int, userID int) error {
	query := "UPDATE orders SET order_status = 'returned', payment_status = 'refunded' WHERE order_id = ? AND user_id = ?"
	result := tx.Exec(query, orderID, userID)
	if result.Error != nil {
		return errors.New("error updating order return status")
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows updated, check if order ID and user ID are correct")
	}
	return nil
}

func (o *OrderRepository) GetCouponDetails(couponCode string) (models.Coupon, error) {
	var coupon models.Coupon
	query := "SELECT id, coupon_code, discount, minimum_required, maximum_allowed, maximum_usage, expire_date FROM coupons WHERE coupon_code = $1"
	err := o.DB.Raw(query, couponCode).Scan(&coupon)
	if err != nil {
		return models.Coupon{}, errors.New("coupon does not exist")
	}
	return coupon, nil
}

func (o *OrderRepository) CheckCouponUsage(userID uint, couponCode string) (int, error) {
	var usageCount int
	query := "SELECT COUNT(*) FROM orders WHERE user_id = $1 AND coupon = $2"
	err := o.DB.Raw(query, userID, couponCode).Scan(&usageCount).Error
	if err != nil {
		return 0, err
	}
	return usageCount, nil
}

// func (o *OrderRepository) RecordCouponUsage(tx *gorm.DB, userID int, couponCode string) error {
// 	type CouponUsage struct {
// 		ID         uint      `gorm:"primarykey;autoIncrement"`
// 		UserID     int       `gorm:"not null"`
// 		CouponCode string    `gorm:"not null"`
// 		UsedAt     time.Time `gorm:"not null"`
// 	}
// 	usage := CouponUsage{
// 		UserID:     userID,
// 		CouponCode: couponCode,
// 		UsedAt:     time.Now(),
// 	}
// 	if err := tx.Create(&usage).Error; err != nil {
// 		return fmt.Errorf("failed to record coupon usage: %w", err)
// 	}
// 	return nil
// }

func (o *OrderRepository) CheckCouponAppliedOrNot(tx *gorm.DB, userID int, couponID string) uint {
	var exist uint
	query := "SELECT COUNT(*) FROM orders WHERE user_id = ? AND coupon_code = ?"
	tx.Raw(query, userID, couponID).Scan(&exist)
	return exist
}
