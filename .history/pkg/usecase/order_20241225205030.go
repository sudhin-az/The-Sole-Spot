package usecase

import (
	"ecommerce_clean_architecture/pkg/domain"
	"ecommerce_clean_architecture/pkg/repository"
	"ecommerce_clean_architecture/pkg/utils/models"
	"errors"

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
	addressExist, err := 
}
