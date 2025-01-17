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
	reviews, err := r.useCase.AddReview(userid, ProductID, review.Rating, review.Comment)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, err.Error(), nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Review Added Successfully", reviews, nil)
	c.JSON(http.StatusOK, successRes)
}

func (r *ReviewHandler) GetReviewsByProductID(c *gin.Context) {
	ProductID := c.Query("product_id")
	if ProductID == "" {
		errRes := response.ClientResponse(http.StatusBadRequest, "Product ID is required", nil, nil)
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	review, err := r.useCase.GetReviewsByProductID(ProductID)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "cannot get reviews", nil, err)
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Reviews Retrieved Successfully", review, nil)
	c.JSON(http.StatusOK, successRes)
}
func (r *ReviewHandler) DeleteReview(c *gin.Context) {
	ReviewID := c.Query("id")
	if ReviewID == "" {
		errRes := response.ClientResponse(http.StatusBadRequest, "Review ID is required", nil, nil)
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	err := r.useCase.DeleteReview(ReviewID)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "cannot delete reviews", nil, err)
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Review Deleted Successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)
}
func (r *ReviewHandler) GetAverageRating(c *gin.Context) {
	ProductID := c.Query("product_id")
	if ProductID == "" {
		errRes := response.ClientResponse(http.StatusBadRequest, "Product ID is required", nil, nil)
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	avgRating, err := r.useCase.GetAverageRating(ProductID)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "cannot getting avg rating", nil, err)
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	responseData := map[string]float64{
		"average_rating": avgRating,
	}

	successRes := response.ClientResponse(http.StatusOK, "avg rating retrieved Successfully", responseData, nil)
	c.JSON(http.StatusOK, successRes)
}
