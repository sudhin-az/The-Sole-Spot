package handlers

import "ecommerce_clean_architecture/pkg/usecase"

type OrderHandler struct {
	orderUseCase usecase.OrderUseCase
}

func NewOrderHandler(usecase usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{
		orderUseCase: usecase,
	}
}
