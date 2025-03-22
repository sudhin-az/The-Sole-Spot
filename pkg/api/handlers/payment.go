package handlers

import (
	"ecommerce_clean_architecture/pkg/usecase"
	"ecommerce_clean_architecture/pkg/utils/models"
	"ecommerce_clean_architecture/pkg/utils/response"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	PaymentUsecase usecase.PaymentUsecase
}

func NewPaymentHandler(paymentUsecase usecase.PaymentUsecase) *PaymentHandler {
	return &PaymentHandler{PaymentUsecase: paymentUsecase}
}

// @Summary Create a new payment
// @Description Creates a new payment order
// @Tags Payment
// @Accept json
// @Produce json
// @Param order_id query string true "Order ID"
// @Param user_id query string true "User ID"
// @Success 200 {object} response.Response{}
// @Failure 400 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /payment [get]

func (pay *PaymentHandler) CreatePayment(c *gin.Context) {
	orderID := c.Query("order_id")
	userID := c.Query("user_id")
	user_ID, err := strconv.Atoi(userID)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Invalid user_id", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	orderDetail, razorID, err := pay.PaymentUsecase.CreatePayment(orderID, user_ID)
	if err != nil {
		if strings.Contains(err.Error(), "Payment failed") {
			errorRes := response.ClientResponse(http.StatusInternalServerError, "Payment failed", nil, err.Error())
			c.JSON(http.StatusInternalServerError, errorRes)
			return
		} else {
			errorRes := response.ClientResponse(http.StatusInternalServerError, "could not generate order details", nil, err.Error())
			c.JSON(http.StatusInternalServerError, errorRes)
			return
		}
	}
	log.Println("OrderDetails: ", orderDetail)
	log.Println("OrderID is: ", orderID)
	log.Println("razorID: ", razorID)
	c.HTML(
		http.StatusOK, "index.html", gin.H{
			"final_price": orderDetail.FinalPrice * 100,
			"razor_id":    razorID,
			"user_id":     userID,
			"order_id":    orderDetail.OrderId,
			"user_name":   orderDetail.Name,
			"total":       orderDetail.FinalPrice,
		})
}

// @Summary Verify online payment
// @Description Verifies an online payment transaction
// @Tags Payment
// @Accept json
// @Produce json
// @Param body body models.OnlinePaymentVerification true "Payment Verification Details"
// @Success 200 {object} response.Response{}
// @Failure 400 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /payment/verify [post]

func (pay *PaymentHandler) OnlinePaymentVerification(c *gin.Context) {
	var onlinePaymentDetails models.OnlinePaymentVerification

	if err := c.ShouldBindJSON(&onlinePaymentDetails); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid request data", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	order, err := pay.PaymentUsecase.OnlinePaymentVerification(onlinePaymentDetails)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "could not update payment details", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Successfully updated payment details", order, nil)
	c.JSON(http.StatusOK, successRes)
}

// @Summary Payment success page
// @Description Renders the success page after payment
// @Tags Payment
// @Produce html
// @Success 200 {string} string "HTML Page"
// @Router /payment/success [get]

func (pay *PaymentHandler) PaymentSuccess(c *gin.Context) {
	c.HTML(http.StatusOK, "success.html", gin.H{})
}
