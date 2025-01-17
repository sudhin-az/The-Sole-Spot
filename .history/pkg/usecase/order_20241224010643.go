package usecase

import (
	"ecommerce_clean_architecture/pkg/repository"
	"ecommerce_clean_architecture/pkg/usecase"
)

type OrderUseCase struct {
	orderRepository repository.OrderRepository
	userRepository  repository.UserRepository
	cartRepository  repository.CartRepository
}

func NewOrderUseCase(orderRepository repository.OrderRepository, userRepository repository.UserRepository, cartRepository repository.CartRepository) usecase.OrderUseCase {
	return &OrderUseCase{
		orderRepository: orderRepository,
		userRepository:  userRepository,
		cartRepository:  cartRepository,
	}
}
