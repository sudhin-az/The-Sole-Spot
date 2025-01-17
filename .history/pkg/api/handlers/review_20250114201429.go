package handlers

import (
	"ecommerce_clean_architecture/pkg/usecase"
	"ecommerce_clean_architecture/pkg/utils/models"
	"ecommerce_clean_architecture/pkg/utils/response"
	"net/http"

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
	ProductID := c.Query("product_id")
	if ProductID == "" {
		errRes := response.ClientResponse(http.StatusBadRequest, "Product ID is required", nil, nil)
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	userID, ok := c.Get("id")
	if !ok {
		errRes := response.ClientResponse(http.StatusUnauthorized, "User ID not found in context", nil, nil)
		c.JSON(http.StatusUnauthorized, errRes)
		return
	}
	userid := userID.(int)

	var review models.ReviewRequest
	if err := c.ShouldBindJSON(&review); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Invalid input data", nil, err)
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	err := r.useCase.AddReview(userid, ProductID, review.Rating, review.Comment)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Invalid input data", nil, err)
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "OTP sent successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)
}
