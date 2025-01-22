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

func (ad *AdminRepository) BeginnTransaction() (*gorm.DB, error) {
	tx := ad.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (ad *AdminRepository) CommittTransaction(tx *gorm.DB) error {
	return tx.Commit().Error
}

func (ad *AdminRepository) RollbackkTransaction(tx *gorm.DB) error {
	return tx.Rollback().Error
}
func (ad *AdminRepository) GetProductStockk(productID int) (int, error) {
	var stock int
	err := ad.DB.Raw("select stock from products where id = ?", productID).Scan(&stock).Error
	if err != nil {
		return 0, err
	}
	return stock, nil
}
func (ar *AdminRepository) UpdateProductStockk(tx *gorm.DB, productID, quantity int) error {
	query := "UPDATE products SET stock = stock + ? WHERE id = ?"
	err := tx.Exec(query, quantity, productID).Error
	if err != nil {
		return fmt.Errorf("failed to update product stock: %w", err)
	}
	return nil
}

func (ar *AdminRepository) AdminOrderRelationship(orderID, userID int) (int, error) {
	var associatedUserID int
	err := ar.DB.Raw("SELECT user_id FROM orders WHERE order_id = ?", orderID).Scan(&associatedUserID).Error
	if err != nil {
		return 0, fmt.Errorf("failed to fetch user ID for order ID %d: %w", orderID, err)
	}
	return associatedUserID, nil
}

func (ad *AdminRepository) GetProductDetailFromOrders(orderID int) ([]models.OrderProducts, error) {
	var orderProductDetails []models.OrderProducts
	err := ad.DB.Raw("SELECT product_id, quantity FROM order_items WHERE order_id = ?", orderID).Scan(&orderProductDetails).Error
	if err != nil {
		return []models.OrderProducts{}, err
	}
	return orderProductDetails, nil
}
func (ad *AdminRepository) Getorderstatus(orderID int) (string, error) {
	var orderStatus string
	err := ad.DB.Raw("SELECT order_status FROM orders WHERE order_id = ?", orderID).Scan(&orderStatus).Error
	if err != nil {
		return "", err
	}
	return orderStatus, nil
}
func (ad *AdminRepository) GetOrderDetails(orderID string) (models.Order, error) {
	var order models.Order
	err := ad.DB.Raw("SELECT * FROM orders WHERE order_id = ?", orderID).Scan(&order).Error
	if err != nil {
		return models.Order{}, err
	}
	return order, nil
}
func (ad *AdminRepository) GetAllOrderDetails() ([]models.FullOrderDetails, error) {
	var orderDetails []models.OrderDetails

	err := ad.DB.Raw("SELECT order_id, final_price, order_status, payment_status FROM orders").Scan(&orderDetails).Error
	if err != nil {
		fmt.Println("Error fetching order details:", err)
		return nil, err
	}

	if len(orderDetails) == 0 {
		fmt.Println("No orders found")
		return nil, nil
	}

	var fullOrderDetails []models.FullOrderDetails
	for _, od := range orderDetails {
		var orderProductDetails []models.OrderProductDetails

		err := ad.DB.Raw(`
			SELECT 
				order_items.product_id, 
				products.name AS product_name, 
				order_items.quantity, 
				order_items.total_price 
			FROM order_items 
			INNER JOIN products ON order_items.product_id = products.id 
			WHERE order_items.order_id = ?`, od.OrderId).Scan(&orderProductDetails).Error

		if err != nil {
			fmt.Println("Error fetching product details for order_id:", od.OrderId, "Error:", err)
			return nil, err
		}

		fullOrderDetails = append(fullOrderDetails, models.FullOrderDetails{
			OrderDetails:        od,
			OrderProductDetails: orderProductDetails,
		})
	}

	return fullOrderDetails, nil
}

func (ad *AdminRepository) Cancelorders(orderID int) error {
	orderStatus := "cancelled"
	err := ad.DB.Exec("UPDATE orders SET order_status = ? WHERE order_id = ?", orderStatus, orderID).Error
	if err != nil {
		return err
	}
	var paymentMethod int
	err = ad.DB.Raw("SELECT payment_method_id FROM orders WHERE order_id = ?", orderID).Scan(&paymentMethod).Error
	if err != nil {
		return err
	}
	if paymentMethod == 1 || paymentMethod == 3 {
		err = ad.DB.Exec("UPDATE orders SET payment_status = 'refunded' WHERE order_id = ?", orderID).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (ad *AdminRepository) UpdatequantityOfproduct(orderProducts []models.OrderProducts) error {
	for _, od := range orderProducts {
		var quantity int
		err := ad.DB.Raw("select quantity from products where id = ?", od.ProductID).Scan(&quantity).Error
		if err != nil {
			return err
		}
		od.Quantity += quantity
		if err := ad.DB.Exec("update products set quantity = ? where id = ?", od.Quantity, od.ProductID).Error; err != nil {
			return err
		}
	}
	return nil
}
func (ad *AdminRepository) ChangeOrderStatus(orderID string, Status string) (models.Order, error) {
	var changeorderstatus models.Order
	err := ad.DB.Model(&models.Order{}).Where("order_id = ?", orderID).Update("order_status", Status).Error
	if err != nil {
		return models.Order{}, err
	}
	err = ad.DB.Where("order_id = ?", orderID).First(&changeorderstatus).Error
	if err != nil {
		return models.Order{}, err
	}
	orderStatus := "delivered"
	err = ad.DB.Exec("UPDATE orders SET order_status = ? WHERE order_id = ?", orderStatus, orderID).Error
	if err != nil {
		return models.Order{}, err
	}
	var paymentMethod int
	err = ad.DB.Raw("SELECT payment_method_id FROM orders WHERE order_id = ?", orderID).Scan(&paymentMethod).Error
	if err != nil {
		return models.Order{}, err
	}
	if paymentMethod == 1 || paymentMethod == 3 {
		err = ad.DB.Exec("UPDATE orders SET payment_status = 'paid' WHERE order_id = ?", orderID).Error
		if err != nil {
			return models.Order{}, err
		}
	}
	return changeorderstatus, nil
}
