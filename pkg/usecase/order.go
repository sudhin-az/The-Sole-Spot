package usecase

import (
	"ecommerce_clean_architecture/pkg/domain"
	"ecommerce_clean_architecture/pkg/repository"
	"ecommerce_clean_architecture/pkg/utils/models"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"
)

type OrderUseCase struct {
	orderRepository repository.OrderRepository
	userRepository  repository.UserRepository
	cartRepository  repository.CartRepository
}

func NewOrderUseCase(orderRepository repository.OrderRepository, userRepository repository.UserRepository, cartRepository repository.CartRepository) *OrderUseCase {
	return &OrderUseCase{
		orderRepository: orderRepository,
		userRepository:  userRepository,
		cartRepository:  cartRepository,
	}
}
func (o *OrderUseCase) OrderItemsFromCart(order models.Order) (models.Order, error) {
	cartExist, err := o.orderRepository.DoesCartExist(order.UserID)
	if err != nil || !cartExist {
		return models.Order{}, errors.New("cart is empty; cannot place order")
	}

	addressExist, err := o.orderRepository.AddressExist(int(order.AddressID))
	if err != nil || !addressExist {
		return models.Order{}, errors.New("address does not exist")
	}

	tx, err := o.orderRepository.BeginTransaction()
	if err != nil {
		return models.Order{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if r := recover(); r != nil || err != nil {
			_ = o.orderRepository.RollbackTransaction(tx)
		}
	}()

	cartItems, err := o.orderRepository.FetchCartItem(order.UserID)
	if err != nil {
		return models.Order{}, err
	}

	var grandTotal float64
	for _, item := range cartItems {
		grandTotal += item.TotalPrice
	}
	order.GrandTotal = grandTotal - order.Discount
	order.FinalPrice = order.GrandTotal
	order.OrderDate = time.Now()

	if order.PaymentMethod == "COD" {
		if order.FinalPrice > 1000 {
			return models.Order{}, errors.New("cash on delivery is not allowed for orders above 1000")
		}
		order.PaymentMethodID = 1
		order.Method = "Cash"
		order.OrderStatus = "pending"
	}

	for _, item := range cartItems {
		availableStock, err := o.orderRepository.GetProductStock(item.ProductID)
		if err != nil || item.Quantity > availableStock {
			return models.Order{}, fmt.Errorf("insufficient stock for product ID %d", item.ProductID)
		}
		newStock := availableStock - item.Quantity
		err = o.orderRepository.UpdateProductStock(tx, item.ProductID, newStock)
		if err != nil {
			return models.Order{}, fmt.Errorf("failed to update stock for product %d", item.ProductID)
		}
	}

	orderID, err := o.orderRepository.CreateOrder(tx, order)
	if err != nil {
		return models.Order{}, err
	}

	var orderItems []domain.OrderItem
	for _, item := range cartItems {
		orderItems = append(orderItems, domain.OrderItem{
			OrderID:    orderID,
			ProductID:  item.ProductID,
			Quantity:   item.Quantity,
			TotalPrice: float64(item.Price) * float64(item.Quantity),
		})
	}

	err = o.orderRepository.CreateOrderItems(tx, orderItems)
	if err != nil {
		return models.Order{}, fmt.Errorf("failed to create order items: %w", err)
	}

	err = o.orderRepository.CommitTransaction(tx)
	if err != nil {
		return models.Order{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	orderSuccessResponse, err := o.orderRepository.GetBriefOrderDetails(orderID)
	if err != nil {
		return models.Order{}, fmt.Errorf("failed to fetch brief order details: %w", err)
	}

	return orderSuccessResponse, nil
}

func (o *OrderUseCase) GetOrderDetails(userID int) ([]models.FullOrderDetails, error) {

	fullOrderDetails, err := o.orderRepository.GetOrderDetails(userID)
	if err != nil {
		return []models.FullOrderDetails{}, err
	}
	return fullOrderDetails, nil
}

func (o *OrderUseCase) CancelOrders(orderID string, userID int) error {

	orderIDInt, err := strconv.Atoi(orderID)
	if err != nil {
		return fmt.Errorf("invalid order ID format: %w", err)
	}

	userTest, err := o.orderRepository.UserOrderRelationship(orderIDInt, userID)
	if err != nil {
		return err
	}
	if userTest != userID {
		log.Printf("Warning: User %d attempting to cancel order %d belonging to user %d", userID, orderIDInt, userTest)
	}

	tx, err := o.orderRepository.BeginTransaction()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer o.orderRepository.RollbackTransaction(tx)

	orderProductDetails, err := o.orderRepository.GetProductDetailsFromOrders(orderIDInt)
	if err != nil {
		return err
	}
	orderStatus, err := o.orderRepository.GetOrderStatus(orderIDInt)
	if err != nil {
		return err
	}
	if orderStatus == "delivered" {
		return errors.New("items already delivered, cannot cancel")
	}

	if orderStatus == "returned" || orderStatus == "Failed" {
		return fmt.Errorf("the order is in %s, so no point in cancelling", orderStatus)
	}
	if orderStatus == "cancelled" {
		return errors.New("the order is already cancelled, so no point in cancelling")
	}

	err = o.orderRepository.CancelOrders(orderIDInt)
	if err != nil {
		return err
	}

	for _, product := range orderProductDetails {
		availableStock, err := o.orderRepository.GetProductStock(product.ProductID)
		if err != nil {
			return err
		}

		newStock := availableStock + product.Quantity
		err = o.orderRepository.UpdateProductStock(tx, product.ProductID, newStock)
		if err != nil {
			return errors.New("failed to restore product stock")
		}
	}

	err = o.orderRepository.CommitTransaction(tx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	err = o.orderRepository.UpdateQuantityOfProduct(orderProductDetails)
	if err != nil {
		return err
	}

	return nil
}
