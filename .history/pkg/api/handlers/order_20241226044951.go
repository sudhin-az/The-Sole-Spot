package handlers

import (
	"ecommerce_clean_architecture/pkg/domain"
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
	userID, ok := c.Get("id")
	if !ok {
		errRes := response.ClientResponse(http.StatusUnauthorized, "User ID not found in context", nil, nil)
		c.JSON(http.StatusUnauthorized, errRes)
		return
	}
	userid := userID.(int)
	var orderFromCart models.OrderFromCart
	if err := c.ShouldBindJSON(&orderFromCart); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are wrong", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	orderSuccessResponse, err := o.orderUseCase.OrderItemsFromCart(orderFromCart, userid)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "Could not do the order", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}

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
	var orders domain.Order
	if err := c.ShouldBindJSON(&orders); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are wrong", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	if len(orders) == 0 {
		errRes := response.ClientResponse(http.StatusBadRequest, "Cart is empty", nil, nil)
		c.JSON(http.StatusBadRequest, errRes)
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
