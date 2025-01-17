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

func (o *OrderUseCase) OrderItemsFromCart(order models.Order, cartItems []models.OrderFromCart) (domain.OrderSuccessResponse, error) {
	cartExist, err := o.orderRepository.DoesCartExist(order.UserID)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}
	if !cartExist {
		return domain.OrderSuccessResponse{}, errors.New("cart is empty; cannot place order")
	}

	addressExist, err := o.orderRepository.AddressExist(models.OrderFromCart{})
	fmt.Println("Address", cartItems.AddressID)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}
	if !addressExist {
		return domain.OrderSuccessResponse{}, errors.New("address does not exist")
	}

	tx, err := o.orderRepository.BeginTransaction()
	if err != nil {
		return domain.OrderSuccessResponse{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer o.orderRepository.RollbackTransaction(tx)

	orderID, err := o.orderRepository.CreateOrder(tx, order)
	if err != nil {
		return domain.OrderSuccessResponse{}, fmt.Errorf("failed to create order: %w", err)
	}

	var orderItems []domain.OrderItem
	for _, item := range cartItems {
		// Check stock availability
		availableStock, err := o.orderRepository.GetProductStock(item.ProductID)
		if err != nil {
			return domain.OrderSuccessResponse{}, fmt.Errorf("failed to fetch stock for product %d: %w", item.ProductID, err)
		}
		if item.Quantity > availableStock {
			return domain.OrderSuccessResponse{}, fmt.Errorf("insufficient stock for product ID %d; please remove it from the cart", item.ProductID)
		}

		newStock := availableStock - item.Quantity
		err = o.orderRepository.UpdateProductStock(tx, item.ProductID, newStock)
		if err != nil {
			return domain.OrderSuccessResponse{}, fmt.Errorf("failed to update stock for product %d: %w", item.ProductID, err)
		}

		orderItems = append(orderItems, domain.OrderItem{
			OrderID:    orderID,
			ProductID:  item.ProductID,
			Quantity:   item.Quantity,
			TotalPrice: float64(item.Price) * float64(item.Quantity),
		})
	}

	err = o.orderRepository.CreateOrderItems(tx, orderItems)
	if err != nil {
		return domain.OrderSuccessResponse{}, fmt.Errorf("failed to create order items: %w", err)
	}

	err = o.orderRepository.CommitTransaction(tx)
	if err != nil {
		return domain.OrderSuccessResponse{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	order.GrandTotal = 0
	for _, item := range orderItems {
		order.GrandTotal += item.TotalPrice
	}
	order.FinalPrice = order.GrandTotal

	if order.PaymentMethod == "COD" {
		if order.FinalPrice > 1000 {
			return domain.OrderSuccessResponse{}, errors.New("cash on delivery is not allowed for orders above 1000")
		}
		order.PaymentStatus = "not paid"
		order.ShipmentStatus = "pending"
	}

	orderSuccessResponse, err := o.orderRepository.GetBriefOrderDetails(orderID)
	if err != nil {
		return domain.OrderSuccessResponse{}, fmt.Errorf("failed to fetch brief order details: %w", err)
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
