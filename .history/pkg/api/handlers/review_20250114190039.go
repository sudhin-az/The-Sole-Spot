package handlers

import "ecommerce_clean_architecture/pkg/usecase"

type ReviewHandler struct {
	useCase usecase.ReviewUseCase
}

func NewReviewHandler(useCase usecase.ReviewUseCase) *ReviewHandler {
	return &ReviewHandler{
		useCase: useCase,
	}
}

func (r *ReviewHandler) addre {
	
}
