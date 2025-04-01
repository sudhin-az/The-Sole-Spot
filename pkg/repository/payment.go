package repository

import (
	"ecommerce_clean_arch/pkg/utils/models"
	"errors"

	"gorm.io/gorm"
)

type PaymentRepository struct {
	DB *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{DB: db}
}
func (pay *PaymentRepository) AddRazorPayDetails(orderID string, razorPayOrderID string) error {

	razorPay := models.RazorPay{OrderID: orderID, RazorID: razorPayOrderID}
	err := pay.DB.Create(&razorPay).Error
	if err != nil {
		return err
	}
	return nil
}
func (pay *PaymentRepository) GetOrderDetailsByOrderId(orderID string) (models.CombinedOrderDetails, error) {
	var orderDetails models.CombinedOrderDetails
	err := pay.DB.Raw(`
	SELECT 
		orders.order_id, orders.final_price, orders.order_status, 
		orders.payment_status, users.first_name, users.email, users.phone,
		addresses.house_name, addresses.street, addresses.city, 
		addresses.district, addresses.state, addresses.pin 
	FROM orders 
	INNER JOIN users ON orders.user_id = users.id 
	INNER JOIN addresses ON orders.address_id = addresses.id 
	WHERE orders.order_id = ? 
`, orderID).Scan(&orderDetails).Error
	if err != nil {
		return models.CombinedOrderDetails{}, err
	}
	if orderDetails.OrderId == "" {
		return models.CombinedOrderDetails{}, errors.New("order not found for this user")
	}
	return orderDetails, nil
}
func (pay *PaymentRepository) CheckPaymentStatus(orderID int) (string, error) {
	var paymentStatus string
	err := pay.DB.Raw("select payment_status from orders where order_id = ?", orderID).Scan(&paymentStatus).Error
	if err != nil {
		return "", err
	}
	return paymentStatus, nil
}

func (pay *PaymentRepository) UpdateOnlinePaymentSucess(orderID int) (*[]models.CombinedOrderDetails, error) {
	var orders []models.CombinedOrderDetails
	err := pay.DB.Raw("UPDATE orders set payment_status = 'paid',order_status = 'success' where order_id = ?", orderID).Scan(&orders).Error
	if err != nil {
		return nil, err
	}
	return &orders, nil
}
func (pay *PaymentRepository) UpdatePaymentDetails(orderID int, paymentID string) error {

	err := pay.DB.Model(&models.RazorPay{}).Where("order_id = ?", orderID).Update("payment_id", paymentID).Error
	if err != nil {
		return err
	}
	return nil

}
