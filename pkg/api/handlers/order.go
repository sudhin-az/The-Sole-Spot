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
	couponCode := c.Query("coupon_code")
	if couponCode == "" {
		errRes := response.ClientResponse(http.StatusBadRequest, "Coupon Code is required", nil, nil)
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

	var orderRequest struct {
		UserID        int
		AddressID     int    `json:"address_id" binding:"required"`
		PaymentMethod string `json:"payment_method" binding:"required"`
	}

	if err := c.ShouldBindJSON(&orderRequest); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Invalid input data", nil, err)
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	addressId := orderRequest.AddressID
	paymentMethod := orderRequest.PaymentMethod
	order := models.Order{
		UserID:        userid,
		AddressID:     uint(addressId),
		PaymentMethod: paymentMethod,
		CouponCode:    couponCode,
	}
	orderSuccessResponse, err := o.orderUseCase.OrderItemsFromCart(order)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "Could not process the order", nil, err.Error())
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
	if !ok {
		errRes := response.ClientResponse(http.StatusUnauthorized, "invalid User ID Format", nil, nil)
		c.JSON(http.StatusUnauthorized, errRes)
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
	if orderID == "" {
		errRes := response.ClientResponse(http.StatusBadRequest, "Order ID is required", nil, nil)
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

	err := o.orderUseCase.CancelOrders(orderID, userid)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Request not correct", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Cancel Successful", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

func (o *OrderHandler) CancelOrderItem(c *gin.Context) {
	orderItemID := c.Query("order_item_id")
	if orderItemID == "" {
		errRes := response.ClientResponse(http.StatusBadRequest, "Order Item ID is required", nil, nil)
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

	orderDetails, err := o.orderUseCase.CancelOrderItem(orderItemID, userid)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Failed to cancel order item", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Order item cancelled successfully", orderDetails, nil)
	c.JSON(http.StatusOK, successRes)
}

func (o *OrderHandler) ReturnUserOrder(c *gin.Context) {
	orderID := c.Query("order_id")
	if orderID == "" {
		errRes := response.ClientResponse(http.StatusBadRequest, "Order ID is required", nil, nil)
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

	err := o.orderUseCase.ReturnUserOrder(orderID, userid)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "failed to Return the order", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "order returned successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

func (o *OrderHandler) GenerateInvoice(c *gin.Context) {
	userID, ok := c.Get("id")
	if !ok {
		errRes := response.ClientResponse(http.StatusUnauthorized, "User ID not found in context", nil, nil)
		c.JSON(http.StatusUnauthorized, errRes)
		return
	}
	userid := userID.(int)

	orderID := c.Query("order_id")

	if orderID == "" {
		errRes := response.ClientResponse(http.StatusBadRequest, "order ID is required", nil, nil)
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	pdf, err := o.orderUseCase.GenerateInvoice(orderID, userid)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not generate invoice", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}

	c.Header("Content-Disposition", "attachment; filename=invoice.pdf")
	c.Header("Content-Type", "application/pdf")

	err = pdf.Output(c.Writer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "the invoice is generated", pdf, nil)
	c.JSON(http.StatusOK, successRes)
}
