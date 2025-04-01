package usecase

import (
	"ecommerce_clean_arch/pkg/domain"
	"ecommerce_clean_arch/pkg/repository"
	"ecommerce_clean_arch/pkg/utils"
	"ecommerce_clean_arch/pkg/utils/models"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jung-kurt/gofpdf"
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
	var rawTotal float64
	var categoryDiscount float64

	for _, item := range cartItems {
		grandTotal += item.TotalPrice
		rawTotal += item.Price * float64(item.Quantity)
		categoryDiscount += float64(item.CategoryDiscount)
	}

	order.GrandTotal = grandTotal
	order.RawTotal = rawTotal

	order.FinalPrice = utils.RoundToTwoDecimalPlaces(grandTotal - categoryDiscount)
	if order.FinalPrice < 0 {
		order.FinalPrice = 0
	}
	order.CategoryDiscount = categoryDiscount
	order.OrderDate = time.Now()

	deliveryCharge := 0.0
	if order.FinalPrice < 1000 {
		deliveryCharge = 50.0
	}
	order.DeliveryCharge = deliveryCharge
	order.FinalPrice += deliveryCharge

	if order.CouponCode != "" {
		couponData, err := o.CouponRepo.CheckCouponExpired(tx, order.CouponCode)
		if err != nil {
			return models.Order{}, fmt.Errorf("failed to fetch coupon details: %w", err)
		}

		if order.FinalPrice < float64(couponData.MinimumRequired) {
			return models.Order{}, fmt.Errorf("order price does not meet coupon requirements (Total Price: %f, Coupon: %s, Minimum Required: %d)", order.FinalPrice, order.CouponCode, couponData.MinimumRequired)
		}

		if couponData.EndDate.Before(time.Now()) {
			return models.Order{}, errors.New("Coupon has expired")
		}

		if exist := o.orderRepository.CheckCouponAppliedOrNot(tx, order.UserID, order.CouponCode); exist >= couponData.MaximumUsage {
			return models.Order{}, fmt.Errorf("coupon %s already applied %d times", order.CouponCode, exist)
		}

		order.CouponID = &couponData.ID
		order.Discount = float64(couponData.Discount)

		discount := (order.GrandTotal * float64(couponData.Discount)) / 100
		if discount > float64(couponData.MaximumAllowed) {
			discount = float64(couponData.MaximumAllowed)
		}

		order.FinalPrice = utils.RoundToTwoDecimalPlaces(order.FinalPrice - discount)
		if order.FinalPrice < 0 {
			order.FinalPrice = 0
		}

		order.DiscountAmount = utils.RoundToTwoDecimalPlaces(discount)
	}

	switch order.PaymentMethod {
	case "COD":
		if order.FinalPrice > 1000 {
			return models.Order{}, errors.New("cash on delivery is not allowed for orders above 1000")
		}
		order.PaymentMethodID = 1
		order.OrderStatus = "pending"
		order.PaymentStatus = "not paid"

	case "ONLINE":
		order.PaymentMethodID = 2
		order.OrderStatus = "pending"
		order.PaymentStatus = "not paid"

	case "WALLET":
		userWallet, err := o.orderRepository.GetWalletAmount(tx, order.UserID)
		if err != nil {
			return models.Order{}, err
		}

		if userWallet < order.FinalPrice {
			return models.Order{}, errors.New("wallet amount is less than total amount")
		}

		if userWallet < 0 {
			return models.Order{}, errors.New("wallet amount is invalid")
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
		order.PaymentStatus = "paid"
		order.OrderStatus = "success"

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
			TotalPrice: item.TotalPrice,
		})
	}

	err = o.orderRepository.CreateOrderItems(tx, orderItems)
	if err != nil {
		return models.Order{}, fmt.Errorf("failed to create order items: %w", err)
	}

	for _, item := range cartItems {
		err := o.cartRepository.RemoveFromCart(item.UserID, item.ProductID)
		if err != nil {
			return models.Order{}, err
		}
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
	log.Println("paymentStatus: ", paymentStatus)

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

	log.Println("Proceeding with cancellation...")

	err = o.orderRepository.CancelOrders(tx, orderIDInt)
	if err != nil {
		log.Println("7------------", err)
		return err
	}
	log.Println("Order status updated to cancelled.")

	err = o.orderRepository.UpdatePaymentStatus(tx, orderIDInt, "refunded")
	if err != nil {
		log.Println("8------------", err)
		return err
	}
	var totalRefundAmount float64
	log.Println("OrderProductDetails:", orderProductDetails)
	for _, product := range orderProductDetails {
		totalRefundAmount += product.FinalPrice
	}
	log.Println("totalrefundamount", totalRefundAmount)

	// if totalRefundAmount <= 0 {
	// 	return errors.New("refund amount is zero; cannot update wallet")
	// }

	newBalance, err := o.walletRepository.CreateOrUpdateWallet(tx, userID, uint(totalRefundAmount))
	if err != nil {
		log.Println("9------------", err)
		return err
	}
	log.Println("New wallet balance after refund:", newBalance)

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

func (o *OrderUseCase) GenerateInvoice(orderID string, userID int) (*gofpdf.Fpdf, error) {
	orderDetails, err := o.orderRepository.FetchOrderDetailsFromDB(orderID)
	if err != nil {
		return nil, errors.New("Unable to fetch order Details")
	}

	log.Println("Order ID:", orderID)
	log.Println("Fetched Order Status:", orderDetails.OrderStatus)

	if orderDetails.OrderStatus != models.Confirm {
		return nil, errors.New("order status is not success")
	}

	invoiceNumber := fmt.Sprintf("INV-%s", orderID)
	date := time.Now().Format("02 Jan 2006")
	dueDate := time.Now().AddDate(0, 0, 15).Format("02 Jan 2006")

	// Seller Info
	sellerInfo := fmt.Sprintf("Axis Bank\nAccount Name: Sole-Spot\nAccount No.: 123-456-7890\nPay by: %v", orderDetails.OrderDate)

	// Buyer Info
	buyerInfo := fmt.Sprintf("%s\nphone: %v\n%s %s\n%s, %s %s",
		orderDetails.CustomerName,
		orderDetails.CustomerPhoneNumber,
		orderDetails.CustomerAddress.HouseName,
		orderDetails.CustomerAddress.State,
		orderDetails.CustomerAddress.Street,
		orderDetails.CustomerAddress.City,
		orderDetails.CustomerAddress.Pin,
	)
	log.Println("CustomerName:", buyerInfo)
	// Initialize PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 20)
	pdf.CellFormat(0, 15, "Tax Invoice", "", 1, "C", false, 0, "")
	pdf.Ln(10)

	// Invoice and Date
	pdf.SetFont("Arial", "", 12)
	pdf.SetXY(130, 20)
	pdf.CellFormat(0, 6, fmt.Sprintf("Invoice No: %s", invoiceNumber), "", 1, "R", false, 0, "")
	pdf.SetXY(130, 26)
	pdf.CellFormat(0, 6, fmt.Sprintf("Date: %s", date), "", 1, "R", false, 0, "")
	pdf.SetXY(130, 32)
	pdf.CellFormat(0, 6, fmt.Sprintf("Due Date: %s", dueDate), "", 1, "R", false, 0, "")

	// Billing and Shipping
	pdf.Ln(15)
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 10, "Billed To:", "", 1, "", false, 0, "")
	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(0, 10, sellerInfo, "", "L", false)
	pdf.Ln(5)
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 10, "Shipped To:", "", 1, "", false, 0, "")
	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(0, 10, buyerInfo, "", "L", false)

	// Items Table
	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 12)
	pdf.SetFillColor(230, 230, 230)
	pdf.CellFormat(80, 10, "Item", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 10, "Quantity", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 10, "Price", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 10, "Total", "1", 1, "C", true, 0, "")

	pdf.SetFont("Arial", "", 12)
	var subtotal float64

	for _, item := range orderDetails.Items {
		total := float64(item.Quantity) * item.Price
		subtotal += total
		pdf.CellFormat(80, 10, item.Name, "1", 0, "", false, 0, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("%d", item.Quantity), "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 10, fmt.Sprintf("%.2f", item.Price), "1", 0, "R", false, 0, "")
		pdf.CellFormat(40, 10, fmt.Sprintf("%.2f", total), "1", 1, "R", false, 0, "")
	}

	// Additional Charges and Final Total
	totalOfferAmount := orderDetails.RawAmount - orderDetails.GrandTotal
	totalDiscountAmount := orderDetails.Discount
	categoryDiscount := orderDetails.CategoryDiscount
	totalDeliveryCharge := orderDetails.DeliveryCharge
	finalGrandTotal := orderDetails.FinalPrice

	// Totals Section
	pdf.Ln(5)
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(150, 10, "Subtotal", "1", 0, "R", false, 0, "")
	pdf.CellFormat(40, 10, fmt.Sprintf("%.2f", subtotal), "1", 1, "R", false, 0, "")
	pdf.CellFormat(150, 10, "Total Offer Amount", "1", 0, "R", false, 0, "")
	pdf.CellFormat(40, 10, fmt.Sprintf("%.2f", totalOfferAmount), "1", 1, "R", false, 0, "")
	pdf.CellFormat(150, 10, "Total Discount Amount", "1", 0, "R", false, 0, "")
	pdf.CellFormat(40, 10, fmt.Sprintf("%.2f", totalDiscountAmount), "1", 1, "R", false, 0, "")
	pdf.CellFormat(150, 10, "Total Category Discount Amount", "1", 0, "R", false, 0, "")
	pdf.CellFormat(40, 10, fmt.Sprintf("%.2f", categoryDiscount), "1", 1, "R", false, 0, "")
	pdf.CellFormat(150, 10, "Total Delivery Charge", "1", 0, "R", false, 0, "")
	pdf.CellFormat(40, 10, fmt.Sprintf("%.2f", totalDeliveryCharge), "1", 1, "R", false, 0, "")
	pdf.CellFormat(150, 10, "Grand Total", "1", 0, "R", false, 0, "")
	pdf.CellFormat(40, 10, fmt.Sprintf("%.2f", finalGrandTotal), "1", 1, "R", false, 0, "")

	return pdf, nil
}
