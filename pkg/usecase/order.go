package usecase

import "ecommerce_clean_architecture/pkg/repository"

type OrderUseCase struct {
	orderRepository repository.OrderRepository
}

func NewOrderUseCase(orderRepository repository.OrderRepository) *OrderUseCase {
	return &OrderUseCase{
		orderRepository: orderRepository,
	}
}
