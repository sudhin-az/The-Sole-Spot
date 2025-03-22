package handlers

import (
	"ecommerce_clean_architecture/pkg/usecase"
	"ecommerce_clean_architecture/pkg/utils/models"
	"ecommerce_clean_architecture/pkg/utils/response"
	"log"
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

// @Summary Create a new admin
// @Description Creates a new admin account
// @Tags Admin
// @Accept json
// @Produce json
// @Param admin body models.AdminSignUp true "Admin details"
// @Success 201 {object} response.Response{}
// @Failure 400 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /admin/signup [post]

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

//AdminLogin godoc
// @Summary Login an admin
// @Description Logs in an admin account
// @Tags Admin
// @Accept json
// @Produce json
// @Param admin body models.AdminLogin true "Admin login details"
// @Success 201 {object} response.Response{}
// @Failure 400 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /admin/login [post]

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

// GetUsers godoc
// @Summary Get all users
// @Description Retrieves a list of all users
// @Tags Admin
// @Produce json
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /admin/users [get]

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

// BlockUser godoc
// @Summary Block a user
// @Description Blocks a user account
// @Tags Admin
// @Produce json
// @Param id query int true "User ID to block"
// @Success 200 {object} response.Response{}
// @Failure 400 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /admin/block-user [get]

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

// UnBlockUser godoc
// @Summary Unblock a user
// @Description Unblocks a user account
// @Tags Admin
// @Produce json
// @Param id query int true "User ID to unblock"
// @Success 200 {object} response.Response{}
// @Failure 400 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /admin/unblock-user [get]

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

// ListOrders godoc
// @Summary List all orders
// @Description Retrieves a list of all orders
// @Tags Admin
// @Produce json
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /admin/orders [get]

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

// AdminCancelOrders godoc
// @Summary Cancel an order
// @Description Cancels an order
// @Tags Admin
// @Produce json
// @Param order_id query string true "Order ID to cancel"
// @Success 200 {object} response.Response{}
// @Failure 400 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /admin/cancel-order [patch]

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

// ChangeOrderStatus godoc
// @Summary Change order status
// @Description Changes the status of an order
// @Tags Admin
// @Accept json
// @Produce json
// @Param order_id query string true "Order ID to update"
// @Param order_status body string true "New order status"
// @Success 200 {object} response.Response{}
// @Failure 400 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /admin/update-order-status [put]

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

// SalesReport godoc
// @Summary Generate sales report
// @Description Generates a sales report for a given date range
// @Tags Admin
// @Produce json
// @Param start_date query string false "Start date of the report period"
// @Param end_date query string false "End date of the report period"
// @Param limit query string false "Limit the report to a specific time period (day, week, month, year)"
// @Param order_status query string false "Filter the report by order status"
// @Success 200 {object} response.Response{}
// @Failure 400 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /admin/sales-report [get]

func (ad *AdminHandler) SalesReport(c *gin.Context) {
	adminID, ok := c.Get("id")
	if !ok {
		errRes := response.ClientResponse(http.StatusUnauthorized, "Admin ID not found in context", nil, nil)
		c.JSON(http.StatusUnauthorized, errRes)
		return
	}
	log.Println("AdminId: ", adminID)

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	limit := c.Query("limit")
	orderStatus := c.Query("order_status")

	if startDate == "" && endDate == "" && limit == "" {
		errRes := response.ClientResponse(http.StatusBadRequest, "Please provide start and end date or a valid limit (day, week, month, year).", nil, nil)
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	startDate, endDate = ad.adminUseCase.GetDateRange(startDate, endDate, limit)

	result, amount, err := ad.adminUseCase.TotalOrders(startDate, endDate, orderStatus)
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

// GenerateSalesReport godoc
// @Summary Generate sales report PDF
// @Description Generates a sales report PDF for a given date range
// @Tags Admin
// @Produce application/pdf
// @Param start_date query string false "Start date of the report period"
// @Param end_date query string false "End date of the report period"
// @Param limit query string false "Limit the report to a specific time period (day, week, month, year)"
// @Param order_status query string false "Filter the report by order status"
// @Success 200 {file} response.Response{}
// @Failure 400 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /admin/sales-report-pdf [get]

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
	orderStatus := c.Query("order_status")

	if startDate == "" && endDate == "" && limit == "" {
		errRes := response.ClientResponse(http.StatusBadRequest, "Please provide start and end date or a valid limit (day, week, month, year).", nil, nil)
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	startDate, endDate = ad.adminUseCase.GetDateRange(startDate, endDate, limit)
	result, amount, err := ad.adminUseCase.TotalOrders(startDate, endDate, orderStatus)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "Error processing orders", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}

	// Call PDF Generation Function
	pdfData, err := ad.adminUseCase.GenerateSalesReportPDF(result, amount, startDate, endDate, orderStatus)
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

// BestSellingProduct godoc
// @Summary Get best selling products
// @Description Retrieves a list of best selling products
// @Tags Admin
// @Produce json
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /admin/best-selling-products [get]

func (ad *AdminHandler) BestSellingProduct(c *gin.Context) {
	_, ok := c.Get("id")
	if !ok {
		errRes := response.ClientResponse(http.StatusUnauthorized, "Admin ID not found in context", nil, nil)
		c.JSON(http.StatusUnauthorized, errRes)
		return
	}
	bestSellingProduct, err := ad.adminUseCase.BestSellingProduct()
	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "could not getting best selling products", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Best Selling Products Retrieved Successfully", bestSellingProduct, nil)
	c.JSON(http.StatusOK, successRes)
}

// BestSellingCategory godoc
// @Summary Get best selling categories
// @Description Retrieves a list of best selling categories
// @Tags Admin
// @Produce json
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /admin/best-selling-categories [get]

func (ad *AdminHandler) BestSellingCategory(c *gin.Context) {
	_, ok := c.Get("id")
	if !ok {
		errRes := response.ClientResponse(http.StatusUnauthorized, "Admin ID not found in context", nil, nil)
		c.JSON(http.StatusUnauthorized, errRes)
		return
	}
	bestSellingCategory, err := ad.adminUseCase.BestSellingCategory()
	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "could not getting best selling categories", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Best Selling categories Retrieved Successfully", bestSellingCategory, nil)
	c.JSON(http.StatusOK, successRes)
}
