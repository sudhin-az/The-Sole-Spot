package handlers

import (
	"ecommerce_clean_architecture/pkg/usecase"
	"ecommerce_clean_architecture/pkg/utils/models"
	"ecommerce_clean_architecture/pkg/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderUseCase usecase.OrderUseCase
}

func NewOrderHandler(useCase usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{
		orderUseCase: useCase,
	}
}

func (o *OrderHandler) OrderItemsFromCart(c *gin.Context) {
	// Extract user ID from context
	userID, ok := c.Get("id")
	if !ok {
		errRes := response.ClientResponse(http.StatusUnauthorized, "User ID not found in context", nil, nil)
		c.JSON(http.StatusUnauthorized, errRes)
		return
	}
	userid := userID.(int)

	// Prepare a slice to store cart items
	var cartItems []models.OrderFromCart

	// Bind JSON input (supports both single and multiple items)
	if err := c.ShouldBindJSON(&cartItems); err != nil {
		// Try binding a single object if the initial binding fails
		var singleCartItem models.OrderFromCart
		if singleErr := c.ShouldBindJSON(&singleCartItem); singleErr != nil {
			errorRes := response.ClientResponse(http.StatusBadRequest, "Invalid input data", nil, singleErr.Error())
			c.JSON(http.StatusBadRequest, errorRes)
			return
		}
		cartItems = append(cartItems, singleCartItem)
	}

	// Create an order object to pass
	order := models.Order{
		UserID: userid,
		// Set other fields if required (e.g., address ID, payment method, etc.)
	}

	// Call the use case
	orderSuccessResponse, err := o.orderUseCase.OrderItemsFromCart(order, cartItems)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "Could not process the order", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}

	// Return success response
	successRes := response.ClientResponse(http.StatusOK, "Successfully created the order", orderSuccessResponse, nil)
	c.JSON(http.StatusOK, successRes)
}

func (o *OrderHandler) ViewOrders(c *gin.Context) {
	userID, ok := c.Get("id")
	if !ok {
		errRes := response.ClientResponse(http.StatusUnauthorized, "User ID not found in context", nil, nil)
		c.JSON(http.StatusUnauthorized, errRes)
		return
	}
	userid := userID.(int)
	var orders models.Order
	if err := c.ShouldBindJSON(&orders); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are wrong", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	fullOrderDetails, err := o.orderUseCase.GetOrderDetails(userid)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "Could not do the order", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Full Order Details", fullOrderDetails, nil)
	c.JSON(http.StatusOK, successRes)
}

func (o *OrderHandler) CancelOrders(c *gin.Context) {
	orderID := c.Query("order_id")
	userID, ok := c.Get("id")
	if !ok {
		errRes := response.ClientResponse(http.StatusUnauthorized, "User ID not found in context", nil, nil)
		c.JSON(http.StatusUnauthorized, errRes)
		return
	}
	userid := userID.(int)

	err := o.orderUseCase.CancelOrders(orderID, userid)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "request not correct ", nil, err.Error())
		c.JSON(http.StatusBadGateway, errorRes)
		return
	}
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "failed to cannel the order ", nil, err.Error())
		c.JSON(http.StatusBadGateway, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Cancel Successfull", nil, nil)
	c.JSON(http.StatusOK, successRes)
}
