package repository

import (
	"ecommerce_clean_architecture/pkg/domain"
	"ecommerce_clean_architecture/pkg/utils/models"
	"fmt"

	"gorm.io/gorm"
)

type AdminRepository struct {
	DB *gorm.DB
}

func NewAdminRepository(DB *gorm.DB) *AdminRepository {
	return &AdminRepository{
		DB: DB,
	}
}

func (ad *AdminRepository) CheckAdminAvailability(admin models.AdminSignUp) bool {
	var count int

	if err := ad.DB.Raw("select count(*) from admin_details where email = ?", admin.Email).Scan(&count).Error; err != nil {
		return false
	}

	return count > 0
}

func (ad *AdminRepository) SignUpHandler(admin models.AdminSignUp) (models.AdminDetailsResponse, error) {
	var adminDetails models.AdminDetailsResponse

	if err := ad.DB.Raw("insert into admin_details(name, email, password) values(?, ?, ?) returning id, name, email", admin.Name, admin.Email, admin.Password).Scan(&adminDetails).Error; err != nil {
		return models.AdminDetailsResponse{}, err
	}
	return adminDetails, nil
}

func (ad *AdminRepository) LoginHandler(admin models.AdminLogin) (domain.AdminDetails, error) {

	var adminCompareDetails domain.AdminDetails
	if err := ad.DB.Raw("select * from admin_details where email = ?", admin.Email).Scan(&adminCompareDetails).Error; err != nil {
		return domain.AdminDetails{}, err
	}
	return adminCompareDetails, nil
}

func (ad *AdminRepository) GetUsers() ([]models.User, error) {
	var listofusers []models.User
	err := ad.DB.Raw("SELECT * FROM users").Scan(&listofusers).Error
	if err != nil {
		return nil, err
	}
	return listofusers, nil
}

func (ad *AdminRepository) GetUserByID(userID int) (models.User, error) {
	querry := fmt.Sprintf("SELECT * FROM users WHERE id = '%d'", userID)
	var userDetails models.User
	if err := ad.DB.Raw(querry).Scan(&userDetails).Error; err != nil {
		return models.User{}, err
	}
	return userDetails, nil
}

func (ad *AdminRepository) UpdateBlockUserByID(user models.User) error {
	err := ad.DB.Exec("UPDATE users SET blocked = ? WHERE id = ?", user.Blocked, user.ID).Error
	if err != nil {
		return err
	}
	return nil
}

func (o *OrderRepository) BeginnTransaction() (*gorm.DB, error) {
	tx := o.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (o *OrderRepository) CommittTransaction(tx *gorm.DB) error {
	return tx.Commit().Error
}

func (o *OrderRepository) RollbackkTransaction(tx *gorm.DB) error {
	return tx.Rollback().Error
}
func (o *OrderRepository) GetProductStockk(productID int) (int, error) {
	var stock int
	err := o.DB.Raw("select stock from products where id = ?", productID).Scan(&stock).Error
	if err != nil {
		return 0, err
	}
	return stock, nil
}
func (o *OrderRepository) UpdateProductStockk(tx *gorm.DB, productID int, newStock int) error {
	return tx.Model(&domain.Products{}).Where("id = ?", productID).Update("stock", newStock).Error
}
func (o *OrderRepository) AdminOrderRelationship(orderID int, userID int) (int, error) {
	var testUserID int
	err := o.DB.Raw("SELECT user_id FROM orders WHERE order_id = ?", orderID).Scan(&testUserID).Error
	if err != nil {
		return -1, err
	}
	return testUserID, nil
}
func (o *OrderRepository) GetProductDetailFromOrders(orderID int) ([]models.OrderProducts, error) {
	var orderProductDetails []models.OrderProducts
	err := o.DB.Raw("SELECT product_id, quantity FROM order_items WHERE order_id = ?", orderID).Scan(&orderProductDetails).Error
	if err != nil {
		return []models.OrderProducts{}, err
	}
	return orderProductDetails, nil
}
func (o *OrderRepository) Getshipmentstatus(orderID int) (string, error) {
	var shipmentStatus string
	err := o.DB.Raw("SELECT shipment_status FROM orders WHERE order_id = ?", orderID).Scan(&shipmentStatus).Error
	if err != nil {
		return "", err
	}
	return shipmentStatus, nil
}

func (o *OrderRepository) Getorderdetails(userID int) ([]models.FullOrderDetails, error) {
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
func (o *OrderRepository) Cancelorders(orderID int) error {
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
