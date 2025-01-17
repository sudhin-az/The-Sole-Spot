package handlers

import (
	"ecommerce_clean_architecture/pkg/usecase"
	"ecommerce_clean_architecture/pkg/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	cartUseCase usecase.CartUseCase
}

func NewCartHandler(usecase usecase.CartUseCase) *CartHandler {
	return &CartHandler{
		cartUseCase: usecase,
	}
}

func (rt *CartHandler) AddToCart(c *gin.Context) {
	userID, ok := c.Get("id")
	if !ok {
		errRes := response.ClientResponse(http.StatusUnauthorized, "User ID not found in context", nil, nil)
		c.JSON(http.StatusUnauthorized, errRes)
		return
	}
	ID := userID.(int)

	type AddToCartRequest struct {
		ProductID int `json:"product_id" binding:"required"`
		Quantity  int `json:"quantity" binding:"required"`
	}

	var req AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Invalid input data", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	cart, err := rt.cartUseCase.AddToCart(ID, req.ProductID, req.Quantity)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "The product cannot be added to cart", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "The product is added to the cart successfully", cart, nil)
	c.JSON(http.StatusOK, successRes)
}
func (rt *CartHandler) RemoveFromCart(c *gin.Context) {
	userID, ok := c.Get("id")
	if !ok {
		errRes := response.ClientResponse(http.StatusUnauthorized, "User ID not found in context", nil, nil)
		c.JSON(http.StatusUnauthorized, errRes)
		return
	}
	ID := userID.(int)

	productID := c.Query("product_id")
	if productID == "" {
		errRes := response.ClientResponse(http.StatusBadRequest, "product ID is required", nil, nil)
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	product_id, err := strconv.Atoi(productID)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Invalid product ID", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	cart, err := rt.cartUseCase.RemoveProductFromCart(ID, product_id)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "could not remove the product from cart", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "product is removed successfully", cart, nil)
	c.JSON(http.StatusOK, successRes)
}
func (rt *CartHandler) DisplayCart(c *gin.Context) {
	userID, ok := c.Get("id")
	if !ok {
		errRes := response.ClientResponse(http.StatusUnauthorized, "User ID not found in context", nil, nil)
		c.JSON(http.StatusUnauthorized, errRes)
		return
	}
	ID := userID.(int)
	cart, err := rt.cartUseCase.DisplayCart(ID)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "could not displayed cart", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	if len(cart) == 0 {
		errRes := response.ClientResponse(http.StatusBadRequest, "Cart is empty", nil, nil)
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "cart items displayed successfully", cart, nil)
	c.JSON(http.StatusOK, successRes)
}
