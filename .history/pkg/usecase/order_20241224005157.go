package usecase

import (
	"ecommerce_clean_architecture/pkg/repository"
)

type OrderUseCase struct {
	orderRepository repository.OrderRepository
	userRepository  repository.UserRepository
	cartRepository  repository.CartRepository
}
