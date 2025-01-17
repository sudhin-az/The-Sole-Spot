package handlers

import "ecommerce_clean_architecture/pkg/usecase"

type ReviewHandler struct {
	useCase usecase.ReviewUseCase
}

func NewReviewHandler(useCase us) {

}
