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
	orderRepository  repository.OrderRepository
	userRepository   repository.UserRepository
	cartRepository   repository.CartRepository
	walletRepository repository.WalletRepository
	WalletUseCase    WalletUseCase
	CouponRepo       repository.CouponRepository
}

func NewOrderUseCase(orderRepository repository.OrderRepository, userRepository repository.UserRepository, cartRepository repository.CartRepository, walletRepository repository.WalletRepository, walletUseCase WalletUseCase, couponRepository repository.CouponRepository) *OrderUseCase {
	return &OrderUseCase{
		orderRepository:  orderRepository,
		userRepository:   userRepository,
		cartRepository:   cartRepository,
		walletRepository: walletRepository,
		WalletUseCase:    walletUseCase,
		CouponRepo:       couponRepository,
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
		if err != nil {
			_ = o.orderRepository.RollbackTransaction(tx)
		}
	}()

	cartItems, err := o.orderRepository.FetchCartItem(tx, order.UserID)
	if err != nil {
		return models.Order{}, err
	}

	var grandTotal float64
	for _, item := range cartItems {
		grandTotal += item.TotalPrice
	}
	order.GrandTotal = grandTotal
	order.FinalPrice = order.GrandTotal
	order.OrderDate = time.Now()

	if order.CouponCode != "" {
		couponData, err := o.CouponRepo.CheckCouponExpired(tx, order.CouponCode)
		if err != nil {
			return models.Order{}, fmt.Errorf("failed to fetch coupon details: %w", err)
		}
		if order.FinalPrice < float64(couponData.MinimumRequired) || order.FinalPrice >= float64(couponData.MaximumAllowed) {
			return models.Order{}, fmt.Errorf("order price does not meet coupon requirements (Total Price: %f, Coupon: %s, Maximum Allowed: %d)", order.FinalPrice, order.CouponCode, couponData.MaximumAllowed)
		}
		if couponData.EndDate.Before(time.Now()) {
			return models.Order{}, errors.New("Coupon has expired")
		}
		if exist := o.orderRepository.CheckCouponAppliedOrNot(tx, order.UserID, order.CouponCode); exist >= couponData.MaximumUsage {
			return models.Order{}, fmt.Errorf("coupon %s already applied %d times", order.CouponCode, exist)
		}
		order.CouponID = &couponData.ID
		order.Discount = float64(couponData.Discount)

		discount := (order.FinalPrice * float64(couponData.Discount)) / 100
		if discount > float64(couponData.MaximumAllowed) {
			discount = float64(couponData.MaximumAllowed)
		}
		order.FinalPrice = order.FinalPrice - discount
		if order.FinalPrice < 0 {
			order.FinalPrice = 0
		}
	}

	//COD
	switch order.PaymentMethod {
	case "COD":
		if order.FinalPrice > 1000 {
			return models.Order{}, errors.New("cash on delivery is not allowed for orders above 1000")
		}
		order.PaymentMethodID = 1
		order.Method = "Cash"
		order.OrderStatus = "pending"
		order.PaymentStatus = "not paid"
		//Online
	case "ONLINE":
		order.PaymentMethodID = 2
		order.Method = "Razorpay"
		order.OrderStatus = "pending"
		order.PaymentStatus = "not paid"

		//Wallet
	case "WALLET":
		userWallet, err := o.orderRepository.GetWalletAmount(tx, order.UserID)
		if err != nil {
			return models.Order{}, err
		}
		if userWallet < order.FinalPrice {
			return models.Order{}, errors.New("wallet amount is less than total amount")
		}
		neweBalance := userWallet - order.FinalPrice
		err = o.orderRepository.UpdateWalletAmount(tx, neweBalance, order.UserID)
		if err != nil {
			return models.Order{}, fmt.Errorf("failed to update wallet: %w", err)
		}
		walletTxn := models.WalletTransaction{
			UserID:      order.UserID,
			Debit:       uint(order.FinalPrice),
			EventDate:   time.Now(),
			TotalAmount: uint(neweBalance),
		}
		if err := o.walletRepository.WalletTransaction(tx, walletTxn); err != nil {
			return models.Order{}, fmt.Errorf("failed to record wallet transaction: %w", err)
		}
		order.PaymentMethodID = 3
		order.Method = "Wallet"
		order.PaymentStatus = "success"
	default:
		return models.Order{}, errors.New("unsupported payment method")
	}

	for _, item := range cartItems {
		availableStock, err := o.orderRepository.GetProductStock(tx, item.ProductID)
		if err != nil {
			return models.Order{}, fmt.Errorf("failed to fetch stock for product ID %d: %w", item.ProductID, err)
		}
		if item.Quantity > availableStock {
			return models.Order{}, fmt.Errorf("insufficient stock for product ID %d", item.ProductID)
		}

		newStock := availableStock - item.Quantity
		err = o.orderRepository.UpdateProductStock(tx, item.ProductID, newStock)
		if err != nil {
			return models.Order{}, fmt.Errorf("failed to update stock for product ID %d", item.ProductID)
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
			TotalPrice: (float64(item.Price) * float64(item.Quantity)) - (float64(item.Price)*float64(item.Quantity))*(order.Discount/100),
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
		log.Println("1------------", err)
		return fmt.Errorf("invalid order ID format: %w", err)
	}

	userTest, err := o.orderRepository.UserOrderRelationship(orderIDInt, userID)
	if err != nil {
		log.Println("2------------", err)
		return err
	}
	if userTest != userID {
		log.Printf("Warning: User %d attempting to cancel order %d belonging to user %d", userID, orderIDInt, userTest)
	}

	tx, err := o.orderRepository.BeginTransaction()
	if err != nil {
		log.Println("3------------", err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = o.orderRepository.RollbackTransaction(tx)
	}()

	orderProductDetails, err := o.orderRepository.GetProductDetailsFromOrders(tx, orderIDInt)
	if err != nil {
		log.Println("4------------", err)
		return err
	}

	paymentStatus, err := o.orderRepository.GetPaymentStatus(tx, orderID)
	if err != nil {
		log.Println("5------------", err)
		return errors.New("cannot show the payment status")
	}
	fmt.Println("paymentStatus: ", paymentStatus)

	orderStatus, err := o.orderRepository.GetOrderStatus(tx, orderIDInt)
	if err != nil {
		log.Println("6------------", err)
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

	fmt.Println("Proceeding with cancellation...")

	err = o.orderRepository.CancelOrders(tx, orderIDInt)
	if err != nil {
		log.Println("7------------", err)
		return err
	}
	fmt.Println("Order status updated to cancelled.")

	err = o.orderRepository.UpdatePaymentStatus(tx, orderIDInt, "refunded")
	if err != nil {
		log.Println("8------------", err)
		return err
	}
	var totalRefundAmount float64
	fmt.Println("OrderProductDetails:", orderProductDetails)
	for _, product := range orderProductDetails {
		totalRefundAmount += product.FinalPrice
	}
	fmt.Println("totalrefundamount", totalRefundAmount)

	// if totalRefundAmount <= 0 {
	// 	return errors.New("refund amount is zero; cannot update wallet")
	// }

	newBalance, err := o.walletRepository.CreateOrUpdateWallet(tx, userID, uint(totalRefundAmount))
	if err != nil {
		log.Println("9------------", err)
		return err
	}
	fmt.Println("New wallet balance after refund:", newBalance)

	walletTxn := models.WalletTransaction{
		UserID:      userID,
		Credit:      uint(totalRefundAmount),
		Debit:       0,
		EventDate:   time.Now(),
		TotalAmount: newBalance,
	}

	err = o.walletRepository.WalletTransaction(tx, walletTxn)
	if err != nil {
		log.Println("wallet transaction error:", err)
		return err
	}
	for _, product := range orderProductDetails {
		availableStock, err := o.orderRepository.GetProductStock(tx, product.ProductID)
		if err != nil {
			log.Println("11------------", err)
			return err
		}
		newStock := availableStock + product.Quantity
		err = o.orderRepository.UpdateProductStock(tx, product.ProductID, newStock)
		if err != nil {
			log.Println("12------------", err)
			return errors.New("failed to restore product stock")
		}
	}

	err = o.orderRepository.UpdateQuantityOfProduct(tx, orderProductDetails)
	if err != nil {
		log.Println("13------------", err)
		return err
	}

	err = o.orderRepository.CommitTransaction(tx)
	if err != nil {
		log.Println("14------------", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (o *OrderUseCase) CancelOrderItem(orderItemID string, userID int) (domain.OrderItem, error) {
	orderItemIDInt, err := strconv.Atoi(orderItemID)
	if err != nil {
		return domain.OrderItem{}, fmt.Errorf("invalid order item ID format: %w", err)
	}

	orderUserID, err := o.orderRepository.UserOrderRelationship(orderItemIDInt, userID)
	if err != nil {
		return domain.OrderItem{}, errors.New("order item does not exist")
	}
	if orderUserID != userID {
		return domain.OrderItem{}, errors.New("you are not authorized to return this order!")
	}

	tx, err := o.orderRepository.BeginTransaction()
	if err != nil {
		return domain.OrderItem{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = o.orderRepository.RollbackTransaction(tx)
	}()

	orderItemStatus, err := o.orderRepository.GetOrderStatus(tx, orderItemIDInt)
	if err != nil {
		return domain.OrderItem{}, err
	}
	if orderItemStatus == "cancelled" {
		return domain.OrderItem{}, errors.New("order item was already cancelled!")
	}
	if orderItemStatus == "returned" {
		return domain.OrderItem{}, errors.New("order item already returned!")
	}
	if orderItemStatus != "delivered" {
		return domain.OrderItem{}, errors.New("order item not delivered, cannot cancel!")
	}

	err = o.orderRepository.CancelOrderItem(tx, orderItemIDInt)
	if err != nil {
		return domain.OrderItem{}, errors.New("failed to cancel order item!")
	}

	refundAmount, err := o.orderRepository.GetOrderItemPrice(tx, orderItemIDInt)
	if err != nil {
		return domain.OrderItem{}, err
	}

	newBalance, err := o.walletRepository.CreateOrUpdateWallet(tx, userID, uint(refundAmount))
	if err != nil {
		return domain.OrderItem{}, err
	}
	walletTxn := models.WalletTransaction{
		UserID:      userID,
		Credit:      uint(refundAmount),
		Debit:       0,
		EventDate:   time.Now(),
		TotalAmount: newBalance,
	}

	err = o.walletRepository.WalletTransaction(tx, walletTxn)
	if err != nil {
		return domain.OrderItem{}, err
	}

	prodctID, quantity, err := o.orderRepository.GetOrderItemDetails(tx, orderItemIDInt)
	if err != nil {
		return domain.OrderItem{}, err
	}

	availableStock, err := o.orderRepository.GetProductStock(tx, prodctID)
	if err != nil {
		return domain.OrderItem{}, err
	}
	newStock := availableStock + quantity
	err = o.orderRepository.UpdateProductStock(tx, prodctID, newStock)
	if err != nil {
		return domain.OrderItem{}, errors.New("failed to restore product stock")
	}

	err = o.orderRepository.CommitTransaction(tx)
	if err != nil {
		return domain.OrderItem{}, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return domain.OrderItem{}, nil
}

func (o *OrderUseCase) ReturnUserOrder(orderID string, userID int) error {
	orderIDInt, err := strconv.Atoi(orderID)
	if err != nil {
		return fmt.Errorf("invalid order ID format: %w", err)
	}

	orderUserID, err := o.orderRepository.UserOrderRelationship(orderIDInt, userID)
	if err != nil {
		return errors.New("order item does not exist")
	}
	if orderUserID != userID {
		return errors.New("you are not authorized to return this order!")
	}

	tx, err := o.orderRepository.BeginTransaction()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = o.orderRepository.RollbackTransaction(tx)
	}()

	orderStatus, err := o.orderRepository.GetOrderStatus(tx, orderIDInt)
	if err != nil {
		return err
	}
	if orderStatus == "cancelled" {
		return errors.New("order was cancelled, return not possible!")
	}
	if orderStatus == "returned" {
		return errors.New("order returned already!")
	}
	if orderStatus != "delivered" {
		return errors.New("order not delivered!")
	}

	err = o.orderRepository.UpdateUserOrderReturn(tx, orderIDInt, userID)
	if err != nil {
		return errors.New("return failed!")
	}

	orderProductDetails, err := o.orderRepository.GetProductDetailsFromOrders(tx, orderIDInt)
	if err != nil {
		return err
	}

	var totalRefundAmount float64

	for _, product := range orderProductDetails {
		totalRefundAmount += product.FinalPrice
	}

	// if totalRefundAmount <= 0 {
	// 	return errors.New("refund amount is zero; cannot update wallet")
	// }

	newBalance, err := o.walletRepository.CreateOrUpdateWallet(tx, userID, uint(totalRefundAmount))
	if err != nil {
		return err
	}
	walletTxn := models.WalletTransaction{
		UserID:      userID,
		Credit:      uint(totalRefundAmount),
		Debit:       0,
		EventDate:   time.Now(),
		TotalAmount: newBalance,
	}

	err = o.walletRepository.WalletTransaction(tx, walletTxn)
	if err != nil {
		return err
	}
	for _, product := range orderProductDetails {
		availableStock, err := o.orderRepository.GetProductStock(tx, product.ProductID)
		if err != nil {
			return err
		}
		newStock := availableStock + product.Quantity
		err = o.orderRepository.UpdateProductStock(tx, product.ProductID, newStock)
		if err != nil {
			return errors.New("failed to restore product stock")
		}
	}

	err = o.orderRepository.UpdateQuantityOfProduct(tx, orderProductDetails)
	if err != nil {
		return err
	}

	err = o.orderRepository.CommitTransaction(tx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil

}
