package handlers

import (
	"ecommerce_clean_architecture/pkg/usecase"
	"ecommerce_clean_architecture/pkg/utils/models"
	"ecommerce_clean_architecture/pkg/utils/response"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminUseCase usecase.AdminUseCase
}

func NewAdminHandler(usecase usecase.AdminUseCase) *AdminHandler {
	return &AdminHandler{
		adminUseCase: usecase,
	}
}

func (ad *AdminHandler) SignUpHandler(c *gin.Context) {
	var admin models.AdminSignUp

	if err := c.ShouldBindJSON(&admin); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "fields provided are wrong", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	adminDetails, err := ad.adminUseCase.SignUpHandler(admin)
	if err != nil {
		errRes := response.ClientResponse(http.StatusInternalServerError, "cannot authenticate Admin", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusCreated, "Successfully signed up the user", adminDetails, nil)
	c.JSON(http.StatusCreated, successRes)
}

func (ad *AdminHandler) LoginHandler(c *gin.Context) {
	var admin models.AdminLogin

	if err := c.ShouldBindJSON(&admin); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	adminDetails, err := ad.adminUseCase.LoginHandler(admin)
	if err != nil {
		errRes := response.ClientResponse(http.StatusInternalServerError, "cannot authenticate Admin", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusCreated, "Successfully login the user", adminDetails, nil)
	c.JSON(http.StatusCreated, successRes)
}

func (ad *AdminHandler) GetUsers(c *gin.Context) {

	users, err := ad.adminUseCase.GetUsers()
	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "could not retrieve records of users", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Successfully retrieved the users", users, nil)
	c.JSON(http.StatusOK, successRes)
}

func (ad *AdminHandler) BlockUser(c *gin.Context) {
	UserID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "check the parameters", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}

	err = ad.adminUseCase.BlockUser(UserID)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not block the user", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "the user is blockes", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

func (ad *AdminHandler) UnBlockUsers(c *gin.Context) {
	UserID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "check the parameters", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}

	err = ad.adminUseCase.UnBlockUsers(UserID)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not unblock the user", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "the user is unblocked", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

func (ad *AdminHandler) ListOrders(c *gin.Context) {
	fullOrderDetails, err := ad.adminUseCase.GetAllOrderDetails()
	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "Could not fetch the order details", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}

	if len(fullOrderDetails) == 0 {
		successRes := response.ClientResponse(http.StatusOK, "No orders found", nil, nil)
		c.JSON(http.StatusOK, successRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Full Order Details", fullOrderDetails, nil)
	c.JSON(http.StatusOK, successRes)
}

func (ad *AdminHandler) AdminCancelOrders(c *gin.Context) {
	orderID := c.Query("order_id")
	if orderID == "" {
		errRes := response.ClientResponse(http.StatusBadRequest, "Order ID is required", nil, nil)
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	ID, ok := c.Get("id")
	if !ok {
		errRes := response.ClientResponse(http.StatusUnauthorized, "ID not found in context", nil, nil)
		c.JSON(http.StatusUnauthorized, errRes)
		return
	}
	id := ID.(int)

	err := ad.adminUseCase.CancelOrders(orderID, id)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Failed to cancel order", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Order cancelled successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

func (ad *AdminHandler) ChangeOrderStatus(c *gin.Context) {
	orderID := c.Query("order_id")
	if orderID == "" {
		errRes := response.ClientResponse(http.StatusBadRequest, "Order ID is required", nil, nil)
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	var input struct {
		OrderStatus string `json:"order_status"binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Invalid input format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	order, err := ad.adminUseCase.ChangeOrderStatus(orderID, input.OrderStatus)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "Failed to update order status", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}

	response := gin.H{
		"status_code": 200,
		"message":     "Order status updated successfully",
		"error":       nil,
		"data": gin.H{
			"order_id":     order.OrderId,
			"order_status": order.OrderStatus,
		},
	}

	c.JSON(http.StatusOK, response)
}

func (ad *AdminHandler) SalesReport(c *gin.Context) {
	adminID, ok := c.Get("id")
	if !ok {
		errRes := response.ClientResponse(http.StatusUnauthorized, "Admin ID not found in context", nil, nil)
		c.JSON(http.StatusUnauthorized, errRes)
		return
	}
	fmt.Println("AdminId: ", adminID)

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	limit := c.Query("limit")
	paymentStatus := c.Query("payment_status")

	if startDate == "" && endDate == "" && limit == "" {
		errRes := response.ClientResponse(http.StatusBadRequest, "Please provide start and end date or a valid limit (day, week, month, year).", nil, nil)
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	startDate, endDate = ad.adminUseCase.GetDateRange(startDate, endDate, limit)

	result, amount, err := ad.adminUseCase.TotalOrders(startDate, endDate, paymentStatus)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "Error processing orders", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully created sales report", gin.H{
		"result": result,
		"amount": amount,
	}, nil)

	c.JSON(http.StatusOK, successRes)
}

func (ad *AdminHandler) GenerateSalesReport(c *gin.Context) {
	_, ok := c.Get("id")
	if !ok {
		errRes := response.ClientResponse(http.StatusUnauthorized, "Admin ID not found in context", nil, nil)
		c.JSON(http.StatusUnauthorized, errRes)
		return
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	limit := c.Query("limit")
	paymentStatus := c.Query("payment_status")

	if startDate == "" && endDate == "" && limit == "" {
		errRes := response.ClientResponse(http.StatusBadRequest, "Please provide start and end date or a valid limit (day, week, month, year).", nil, nil)
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	startDate, endDate = ad.adminUseCase.GetDateRange(startDate, endDate, limit)
	result, amount, err := ad.adminUseCase.TotalOrders(startDate, endDate, paymentStatus)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "Error processing orders", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}

	// Call PDF Generation Function
	pdfData, err := ad.adminUseCase.GenerateSalesReportPDF(result, amount, startDate, endDate, paymentStatus)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "Error generating PDF", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}

	// Send PDF Response
	c.Header("Content-Disposition", "attachment; filename=sales_report.pdf")
	c.Header("Content-Type", "application/pdf")
	c.Data(http.StatusOK, "application/pdf", pdfData)
}
