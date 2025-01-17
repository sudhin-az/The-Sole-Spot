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
func (o *OrderUseCase) OrderItemsFromCart(order models.Order) (domain.OrderSuccessResponse, error) {
	cartExist, err := o.orderRepository.DoesCartExist(order.UserID)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}
	if !cartExist {
		return domain.OrderSuccessResponse{}, errors.New("cart is empty; cannot place order")
	}

	addressExist, err := o.orderRepository.AddressExist(int(order.AddressID))
	fmt.Println("AddressID:", order.AddressID)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}
	if !addressExist {
		return domain.OrderSuccessResponse{}, errors.New("address does not exist")
	}

	tx, err := o.orderRepository.BeginTransaction()
	fmt.Println("Transaction started:")
	if err != nil {
		return domain.OrderSuccessResponse{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer o.orderRepository.RollbackTransaction(tx)

	cartItems, err := o.orderRepository.FetchCartItem(order.UserID)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}
	for _, item := range cartItems {
		order.GrandTotal += float64(item.Price)
	}
	order.FinalPrice = order.GrandTotal

	if order.PaymentMethod == "COD" {
		if order.FinalPrice > 1000 {
			return domain.OrderSuccessResponse{}, errors.New("cash on delivery is not allowed for orders above 1000")
		}
		order.PaymentMethodID = 1
		order.PaymentStatus = "not paid"
		order.ShipmentStatus = "pending"
	}
	orderID, err := o.orderRepository.CreateOrder(tx, order)
	fmt.Println("Order Created")
	if err != nil {
		return domain.OrderSuccessResponse{}, fmt.Errorf("failed to create order: %w", err)
	}
	var orderItems []domain.OrderItem
	for _, item := range cartItems {
		// Check stock availability
		availableStock, err := o.orderRepository.GetProductStock(item.ProductID)
		fmt.Println("product stock")
		if err != nil {
			return domain.OrderSuccessResponse{}, fmt.Errorf("failed to fetch stock for product %d: %w", item.ProductID, err)
		}
		if item.Quantity > availableStock {
			return domain.OrderSuccessResponse{}, fmt.Errorf("insufficient stock for product ID %d; please remove it from the cart", item.ProductID)
		}

		newStock := availableStock - item.Quantity

		err = o.orderRepository.UpdateProductStock(tx, item.ProductID, newStock)
		fmt.Println("Stock updated")
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
	fmt.Println("Create order items")
	if err != nil {
		return domain.OrderSuccessResponse{}, fmt.Errorf("failed to create order items: %w", err)
	}

	err = o.orderRepository.CommitTransaction(tx)
	fmt.Println("commit transaction")
	if err != nil {
		return domain.OrderSuccessResponse{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	orderSuccessResponse, err := o.orderRepository.GetBriefOrderDetails(orderID)
	fmt.Println("Get brief order details")
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
