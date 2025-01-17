package usecase

import (
	"ecommerce_clean_architecture/pkg/domain"
	"ecommerce_clean_architecture/pkg/repository"
	"ecommerce_clean_architecture/pkg/utils/models"
	"errors"
	"fmt"
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

func (o *OrderUseCase) FetchCartItemsForUser(userID int) ([]models.OrderFromCart, error) {
	return o.orderRepository.GetCartItems(userID)
}

func (o *OrderUseCase) OrderItemsFromCart(order models.Order, cartItems []models.OrderFromCart) (domain.OrderSuccessResponse, error) {
	// Check Address Validity
	addressExist, err := o.orderRepository.AddressExist(order.AddressID)
	if err != nil || !addressExist {
		return domain.OrderSuccessResponse{}, fmt.Errorf("invalid address ID")
	}

	// Begin Transaction
	tx, err := o.orderRepository.BeginTransaction()
	if err != nil {
		return domain.OrderSuccessResponse{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer o.orderRepository.RollbackTransaction(tx)

	// Create Order
	orderID, err := o.orderRepository.CreateOrder(tx, order)
	if err != nil {
		return domain.OrderSuccessResponse{}, fmt.Errorf("failed to create order: %w", err)
	}

	// Process Order Items
	var orderItems []domain.OrderItem
	for _, item := range cartItems {
		availableStock, err := o.orderRepository.GetProductStock(item.ProductID)
		if err != nil || item.Quantity > availableStock {
			return domain.OrderSuccessResponse{}, fmt.Errorf("insufficient stock for product ID %d", item.ProductID)
		}

		newStock := availableStock - item.Quantity
		if err = o.orderRepository.UpdateProductStock(tx, item.ProductID, newStock); err != nil {
			return domain.OrderSuccessResponse{}, fmt.Errorf("failed to update stock for product ID %d", item.ProductID)
		}

		orderItems = append(orderItems, domain.OrderItem{
			OrderID:    orderID,
			ProductID:  item.ProductID,
			Quantity:   item.Quantity,
			TotalPrice: float64(item.Price) * float64(item.Quantity),
		})
	}

	if err := o.orderRepository.CreateOrderItems(tx, orderItems); err != nil {
		return domain.OrderSuccessResponse{}, fmt.Errorf("failed to create order items: %w", err)
	}

	// Commit Transaction
	if err := o.orderRepository.CommitTransaction(tx); err != nil {
		return domain.OrderSuccessResponse{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Prepare Response
	return o.orderRepository.GetBriefOrderDetails(orderID)
}

func (o *OrderUseCase) GetOrderDetails(userID int) ([]models.FullOrderDetails, error) {

	fullOrderDetails, err := o.orderRepository.GetOrderDetails(userID)
	if err != nil {
		return []models.FullOrderDetails{}, err
	}
	return fullOrderDetails, nil
}

func (o *OrderUseCase) CancelOrders(orderID string, userID int) error {
	userTest, err := o.orderRepository.UserOrderRelationship(orderID, userID)
	if err != nil {
		return err
	}
	if userTest != userID {
		return errors.New("the order is not done by this user")
	}
	// Begin a transaction
	tx, err := o.orderRepository.BeginTransaction()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer o.orderRepository.RollbackTransaction(tx)
	orderProductDetails, err := o.orderRepository.GetProductDetailsFromOrders(orderID)
	if err != nil {
		return err
	}
	shipmentStatus, err := o.orderRepository.GetShipmentStatus(orderID)
	if err != nil {
		return err
	}
	if shipmentStatus == "delivered" {
		return errors.New("items already delivered, cannot cancel")
	}

	if shipmentStatus == "pending" || shipmentStatus == "returned" || shipmentStatus == "Failed" {
		message := fmt.Sprintf(shipmentStatus)

		return errors.New("the order is in" + message + ", so no point in cancelling")
	}
	if shipmentStatus == "cancelled" {
		return errors.New("the order is already cancelled, so no point in cancelling")
	}
	err = o.orderRepository.CancelOrders(orderID)
	if err != nil {
		return err
	}
	for _, product := range orderProductDetails {
		availableStock, err := o.orderRepository.GetProductStock(product.ProductID)
		if err != nil {
			return err
		}

		// Restore stock
		newStock := availableStock + product.Quantity
		err = o.orderRepository.UpdateProductStock(tx, product.ProductID, newStock)
		if err != nil {
			return errors.New("failed to restore product stock")
		}
	}
	// Commit the transaction
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
