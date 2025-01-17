package usecase

import (
	"ecommerce_clean_architecture/pkg/repository"
	"ecommerce_clean_architecture/pkg/repository/interfaces"
)

type OrderUseCase struct {
	orderRepository repository.OrderRepository
	userRepository  repo.UserRepository
	cartRepository  interfaces.CartRepository
}
