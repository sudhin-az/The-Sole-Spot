package handlers

import (
	"ecommerce_clean_architecture/pkg/usecase"

	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	useCase usecase.ReviewUseCase
}

func NewReviewHandler(useCase usecase.ReviewUseCase) *ReviewHandler {
	return &ReviewHandler{
		useCase: useCase,
	}
}

func (r *ReviewHandler) AddReview(c *gin.Context) {

}
