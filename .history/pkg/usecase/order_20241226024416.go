package usecase

import (
	"ecommerce_clean_architecture/pkg/domain"
	"ecommerce_clean_architecture/pkg/helper"
	"ecommerce_clean_architecture/pkg/repository"
	"ecommerce_clean_architecture/pkg/utils/models"
	"errors"
	"fmt"

	"github.com/jinzhu/copier"
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

func (o *OrderUseCase) OrderItemsFromCart(orderFromCart models.OrderFromCart, userID int) (domain.OrderSuccessResponse, error) {
	var orderBody models.OrderIncoming
	err := copier.Copy(&orderBody, &orderFromCart)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}
	orderBody.UserID = uint(userID)
	cartExist, err := o.orderRepository.DoesCartExist(userID)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}
	if !cartExist {
		return domain.OrderSuccessResponse{}, errors.New("cart empty can't order")
	}
	addressExist, err := o.orderRepository.AddressExist(orderBody)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}
	if !addressExist {
		return domain.OrderSuccessResponse{}, errors.New("address does not exist")
	}
	cartItems, err := o.cartRepository.GetAllItemsFromCart(userID)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}

	for _, c := range cartItems {
		// Fetch product stock from the repository
		availableStock, err := o.orderRepository.GetProductStock(c.ProductID)
		if err != nil {
			return domain.OrderSuccessResponse{}, err
		}
		if c.Quantity > availableStock {
			return domain.OrderSuccessResponse{}, errors.New("Insufficient Stock")
		}
	}

	var orderDetails domain.Order
	var orderItemDetails domain.OrderItem

	orderDetails = helper.CopyOrderDetails(orderDetails, orderBody)

	for _, c := range cartItems {
		orderDetails.GrandTotal += c.TotalPrice
	}
	orderDetails.FinalPrice = orderDetails.GrandTotal

	//for cash on delivery
	if orderBody.PaymentID == 1 {

		if orderDetails.FinalPrice > 1000 {
			return domain.OrderSuccessResponse{}, errors.New("cash on delivery is not possible")
		}
		orderDetails.PaymentStatus = "not paid"
		orderDetails.ShipmentStatus = "pending"
	}
	err = o.orderRepository.CreateOrder(orderDetails)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}
	for _, c := range cartItems {
		// for each order save details of products and associated details and use order_id as foreign key ( for each order multiple product will be there)
		orderItemDetails.OrderID = orderDetails.OrderId
		orderItemDetails.ProductID = uint(c.ProductID)
		orderItemDetails.Quantity = c.Quantity
		orderItemDetails.TotalPrice = c.TotalPrice

		err := o.orderRepository.AddOrderItems(orderItemDetails, userID, c.ProductID, float64(c.Quantity))
		if err != nil {
			return domain.OrderSuccessResponse{}, err
		}
	}
	orderSuccessResponse, err := o.orderRepository.GetBriefOrderDetails(orderDetails.OrderId)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}
	return orderSuccessResponse, nil
}

func (o *OrderUseCase) GetOrderDetails(userID int, page int, count int) ([]models.FullOrderDetails, error) {

	fullOrderDetails, err := o.orderRepository.GetOrderDetails(userID, page, count)
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

	orderProductDetails, err := o.orderRepository.GetProductDetailsFromOrders(orderID)
	if err != nil {
		return err
	}
	shipmentStatus, err := o.orderRepository.GetShipmentStatus(orderID)
	if err != nil {
		return err
	}
	if shipmentStatus == "delivered" {
		return errors.New("item already delivered, cannot cancel")
	}

	if shipmentStatus == "returned" || shipmentStatus == "return" {
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
	err = o.orderRepository.UpdateQuantityOfProduct(orderProductDetails)
	if err != nil {
		return err
	}
	return nil
}
