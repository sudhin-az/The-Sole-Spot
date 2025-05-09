package usecase

import (
	"bytes"
	"ecommerce_clean_arch/pkg/domain"
	"ecommerce_clean_arch/pkg/helper"
	"ecommerce_clean_arch/pkg/repository"
	"ecommerce_clean_arch/pkg/utils/models"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/wcharczuk/go-chart"
	"golang.org/x/crypto/bcrypt"
)

type AdminUseCase struct {
	adminrepository repository.AdminRepository
}

func NewAdminUseCase(adminrepository repository.AdminRepository) *AdminUseCase {
	return &AdminUseCase{
		adminrepository: adminrepository,
	}
}

func (ad *AdminUseCase) SignUpHandler(admin models.AdminSignUp) (domain.TokenAdmin, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), 10)
	if err != nil {
		return domain.TokenAdmin{}, errors.New("internal server error")
	}
	admin.Password = string(hashedPassword)

	adminDetails, err := ad.adminrepository.SignUpHandler(admin)
	if err != nil {
		return domain.TokenAdmin{}, err
	}

	return domain.TokenAdmin{
		Admin: adminDetails,
	}, nil

}

func (ad *AdminUseCase) LoginHandler(admin models.AdminLogin) (domain.TokenAdmin, error) {
	adminCompareDetails, err := ad.adminrepository.LoginHandler(admin)
	if err != nil {
		return domain.TokenAdmin{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(adminCompareDetails.Password), []byte(admin.Password))
	if err != nil {
		return domain.TokenAdmin{}, err
	}

	adminDetailsResponse := models.AdminDetailsResponse{
		ID:    adminCompareDetails.ID,
		Name:  adminCompareDetails.Name,
		Email: adminCompareDetails.Email,
	}

	tokenString, err := helper.GenerateTokenAdmin(adminDetailsResponse)
	if err != nil {
		return domain.TokenAdmin{}, err
	}

	return domain.TokenAdmin{
		Admin: adminDetailsResponse,
		Token: tokenString,
	}, nil
}

func (ad *AdminUseCase) GetUsers() ([]models.User, error) {
	userDetails, err := ad.adminrepository.GetUsers()
	if err != nil {
		return nil, err
	}
	return userDetails, nil
}

func (ad *AdminUseCase) BlockUser(userID int) error {
	user, err := ad.adminrepository.GetUserByID(userID)
	if err != nil {
		return errors.New("already blocked")
	} else {
		user.Blocked = true
	}
	err = ad.adminrepository.UpdateBlockUserByID(user)
	if err != nil {
		return err
	}
	return nil
}

func (ad *AdminUseCase) UnBlockUsers(userID int) error {
	user, err := ad.adminrepository.GetUserByID(userID)
	if err != nil {
		return errors.New("already unblocked")
	} else {
		user.Blocked = false
	}
	err = ad.adminrepository.UpdateBlockUserByID(user)
	if err != nil {
		return err
	}
	return nil
}
func (ad *AdminUseCase) GetAllOrderDetails() ([]models.FullOrderDetails, error) {
	fullOrderDetails, err := ad.adminrepository.GetAllOrderDetails()
	if err != nil {
		log.Println("Error in repository:", err)
		return nil, err
	}
	return fullOrderDetails, nil
}

func (ad *AdminUseCase) CancelOrders(orderID string, userID int) error {

	orderIDInt, err := strconv.Atoi(orderID)
	if err != nil {
		return fmt.Errorf("invalid order ID format: %w", err)
	}

	userTest, err := ad.adminrepository.AdminOrderRelationship(orderIDInt, userID)
	if err != nil {
		return err
	}
	if userTest != userID {
		log.Printf("Warning: User %d attempting to cancel order %d belonging to user %d", userID, orderIDInt, userTest)
	}

	tx, err := ad.adminrepository.BeginnTransaction()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer ad.adminrepository.RollbackkTransaction(tx)

	orderProductDetails, err := ad.adminrepository.GetProductDetailFromOrders(orderIDInt)
	if err != nil {
		return err
	}
	orderStatus, err := ad.adminrepository.Getorderstatus(orderIDInt)
	if err != nil {
		return err
	}
	if orderStatus == "delivered" {
		return errors.New("items already delivered, cannot cancel")
	}

	if orderStatus == "returned" || orderStatus == "Failed" {
		return fmt.Errorf("the order is in %s, so no point in cancelling", orderStatus)
	}
	if orderStatus == "cancelled" {
		return errors.New("the order is already cancelled, so no point in cancelling")
	}

	err = ad.adminrepository.Cancelorders(orderIDInt)
	if err != nil {
		return err
	}

	for _, product := range orderProductDetails {
		availableStock, err := ad.adminrepository.GetProductStockk(product.ProductID)
		if err != nil {
			return err
		}

		newStock := availableStock + product.Quantity
		err = ad.adminrepository.UpdateProductStockk(tx, product.ProductID, newStock)
		if err != nil {
			return errors.New("failed to restore product stock")
		}
	}

	err = ad.adminrepository.CommittTransaction(tx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	err = ad.adminrepository.UpdatequantityOfproduct(orderProductDetails)
	if err != nil {
		return err
	}

	return nil
}

func (ad *AdminUseCase) ChangeOrderStatus(orderID string, status string) (models.Order, error) {

	validStatuses := map[string]bool{
		"pending":   true,
		"shipped":   true,
		"cancelled": true,
		"return":    true,
		"delivered": true,
		"failed":    true,
	}
	if !validStatuses[status] {
		return models.Order{}, fmt.Errorf("invalid order status")
	}
	order, err := ad.adminrepository.GetOrderDetails(orderID)
	if err != nil {
		return models.Order{}, fmt.Errorf("order ID does not exist")
	}

	switch status {
	case "shipped":
		if order.OrderStatus == "cancelled" {
			return models.Order{}, fmt.Errorf("cannot ship an order that has already been cancelled")
		}

	case "cancelled":
		if order.OrderStatus == "cancelled" {
			return models.Order{}, fmt.Errorf("order has already been cancelled by user or Admin")
		}
		if order.OrderStatus == "delivered" {
			return models.Order{}, fmt.Errorf("cannot cancel the order as it is already delivered")
		}

	case "delivered":
		if order.OrderStatus == "delivered" {
			return models.Order{}, fmt.Errorf("order has already been delivered")
		}
		if order.OrderStatus != "shipped" {
			return models.Order{}, fmt.Errorf("cannot mark the order as delivered unless it is shipped")
		}

	case "failed":
		if order.OrderStatus == "failed" {
			return models.Order{}, fmt.Errorf("order has already failed")
		}
		if order.OrderStatus == "delivered" {
			return models.Order{}, fmt.Errorf("cannot mark an already delivered order as failed")
		}

	case "return":
		if order.OrderStatus == "return" {
			return models.Order{}, fmt.Errorf("order has already been marked as return")
		}
		if order.OrderStatus != "delivered" {
			return models.Order{}, fmt.Errorf("only delivered orders can be returned")
		}
	}

	updateOrder, err := ad.adminrepository.ChangeOrderStatus(orderID, status)
	if err != nil {
		return models.Order{}, fmt.Errorf("could not get updated order status: %w", err)
	}
	return updateOrder, nil
}

func (ad *AdminUseCase) GetDateRange(startDate, endDate, limit string) (string, string) {
	today := time.Now()
	switch limit {
	case "day":
		startDate = today.AddDate(0, 0, -1).Format("2006-01-02")
		endDate = today.Format("2006-01-02")
	case "week":
		startOfWeek := today.AddDate(0, 0, -int(today.Weekday()))
		startDate = startOfWeek.Format("2006-01-02")
		endDate = today.Format("2006-01-02")
	case "month":
		startDate = time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location()).Format("2006-01-02")
		endDate = time.Date(today.Year(), today.Month()+1, 0, 0, 0, 0, 0, today.Location()).Format("2006-01-02")
	case "year":
		startDate = time.Date(today.Year(), 1, 1, 0, 0, 0, 0, today.Location()).Format("2006-01-02")
		endDate = time.Date(today.Year(), 12, 31, 0, 0, 0, 0, today.Location()).Format("2006-01-02")
	}
	return startDate, endDate
}

func (ad *AdminUseCase) TotalOrders(fromDate, toDate, orderStatus string) (models.OrderCount, models.AmountInformation, error) {
	orders, amount, err := ad.adminrepository.GetTotalOrders(fromDate, toDate, orderStatus)
	if err != nil {
		return models.OrderCount{}, models.AmountInformation{}, fmt.Errorf("failed to get total orders: %w", err)
	}
	return orders, amount, nil
}

func (ad *AdminUseCase) GenerateSalesReportPDF(orderCount models.OrderCount, amountInfo models.AmountInformation, startDate, endDate, orderStatus string) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Report Header
	pdf.SetFont("Arial", "B", 20)
	pdf.CellFormat(0, 10, "Sales Report", "", 1, "C", false, 0, "")
	pdf.Ln(10)

	// Report Duration
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(40, 10, "Report Duration:", "", 1, "L", false, 0, "")
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(40, 10, "Start Date: "+startDate, "", 1, "L", false, 0, "")
	pdf.CellFormat(40, 10, "End Date: "+endDate, "", 1, "L", false, 0, "")
	pdf.CellFormat(40, 10, "Payment Status: "+orderStatus, "", 1, "L", false, 0, "")
	pdf.Ln(12)

	// Summary Information
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Summary Information")
	pdf.Ln(8)

	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(90, 10, "Description", "1", 0, "C", false, 0, "")
	pdf.CellFormat(60, 10, "Amount", "1", 0, "C", false, 0, "")
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 12)
	summaryData := map[string]string{
		"Total Orders":                  strconv.Itoa(int(orderCount.TotalOrder)),
		"Total Amount Before Deduction": fmt.Sprintf("%.2f", amountInfo.TotalAmountBeforeDeduction),
		"Total Coupon Deduction":        fmt.Sprintf("%.2f", amountInfo.TotalCouponDeduction),
		"Total Product Offer Deduction": fmt.Sprintf("%.2f", amountInfo.TotalProuctOfferDeduction),
		"Total Amount After Deduction":  fmt.Sprintf("%.2f", amountInfo.TotalAmountAfterDeduction),
	}

	for desc, amount := range summaryData {
		pdf.CellFormat(90, 10, desc, "1", 0, "L", false, 0, "")
		pdf.CellFormat(60, 10, amount, "1", 0, "R", false, 0, "")
		pdf.Ln(-1)
	}
	pdf.Ln(10)

	// Order History Table
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Order History Details")
	pdf.Ln(8)

	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(60, 10, "Status", "1", 0, "C", false, 0, "")
	pdf.CellFormat(60, 10, "Count", "1", 0, "C", false, 0, "")
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 12)
	orderHistory := map[string]int{
		"Pending":   int(orderCount.TotalPending),
		"Confirmed": int(orderCount.TotalConfirmed),
		"Shipped":   int(orderCount.TotalShipped),
		"Delivered": int(orderCount.TotalDelivered),
		"Cancelled": int(orderCount.TotalCancelled),
		"Returned":  int(orderCount.TotalReturned),
	}

	for status, count := range orderHistory {
		pdf.CellFormat(60, 10, status, "1", 0, "L", false, 0, "")
		pdf.CellFormat(60, 10, strconv.Itoa(count), "1", 0, "R", false, 0, "")
		pdf.Ln(-1)
	}
	pdf.Ln(10)

	// Bar Chart
	chartData := []chart.Value{
		{Value: float64(orderCount.TotalPending), Label: "Pending"},
		{Value: float64(orderCount.TotalConfirmed), Label: "Confirmed"},
		{Value: float64(orderCount.TotalShipped), Label: "Shipped"},
		{Value: float64(orderCount.TotalDelivered), Label: "Delivered"},
		{Value: float64(orderCount.TotalCancelled), Label: "Cancelled"},
		{Value: float64(orderCount.TotalReturned), Label: "Returned"},
	}

	if hasValidData(chartData) {
		barChart := chart.BarChart{
			Width:  500,
			Height: 300,
			Bars:   chartData,
			XAxis: chart.Style{
				Show: true,
			},
			YAxis: chart.YAxis{
				Style: chart.Style{
					Show: true,
				},
				Range: &chart.ContinuousRange{
					Min: 0,
					Max: float64(orderCount.TotalOrder),
				},
			},
		}

		var chartBuffer bytes.Buffer
		err := barChart.Render(chart.PNG, &chartBuffer)
		if err != nil {
			return nil, fmt.Errorf("failed to generate bar chart: %v", err)
		}

		chartFileName := "temp_chart.png"
		err = os.WriteFile(chartFileName, chartBuffer.Bytes(), 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to save chart image: %v", err)
		}
		defer os.Remove(chartFileName)

		pageWidth, pageHeight := pdf.GetPageSize()
		chartWidth := float64(150)
		chartHeight := chartWidth / 1.5
		remainingHeight := pageHeight - pdf.GetY() - 10
		if chartHeight > remainingHeight {
			pdf.AddPage()
			pdf.Ln(5)
		}

		pdf.SetFont("Arial", "B", 14)
		pdf.CellFormat(0, 10, "Order Status Distribution", "", 1, "C", false, 0, "")
		pdf.Ln(5)

		chartX := (pageWidth - chartWidth) / 2
		chartY := pdf.GetY() + 2

		pdf.ImageOptions(
			chartFileName,
			chartX,
			chartY,
			chartWidth,
			chartHeight,
			false,
			gofpdf.ImageOptions{ImageType: "PNG"},
			0,
			"",
		)
		pdf.SetY(chartY + chartHeight + 2)

	} else {
		pdf.SetFont("Arial", "I", 12)
		pdf.CellFormat(0, 10, "No data available for chart representation.", "", 1, "C", false, 0, "")
		pdf.Ln(10)
	}

	// Pie Chart
	pieChartData := []chart.Value{
		{Value: amountInfo.TotalAmountBeforeDeduction, Label: "Total Amount Before Deduction"},
		{Value: amountInfo.TotalCouponDeduction, Label: "Total Coupon Deduction"},
		{Value: amountInfo.TotalProuctOfferDeduction, Label: "Total Product Offer Deduction"},
		{Value: amountInfo.TotalAmountAfterDeduction, Label: "Total Amount After Deduction"},
	}

	if hasValidData(pieChartData) {
		pieChart := chart.PieChart{
			Width:  400,
			Height: 400,
			Values: pieChartData,
		}

		var pieChartBuffer bytes.Buffer
		err := pieChart.Render(chart.PNG, &pieChartBuffer)
		if err != nil {
			return nil, fmt.Errorf("failed to generate pie chart: %v", err)
		}

		pieChartFileName := "temp_pie_chart.png"
		err = os.WriteFile(pieChartFileName, pieChartBuffer.Bytes(), 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to save pie chart image: %v", err)
		}
		defer os.Remove(pieChartFileName)

		pageWidth, pageHeight := pdf.GetPageSize()
		pieChartWidth := float64(150)
		pieChartHeight := float64(150)
		remainingHeight := pageHeight - pdf.GetY() - 10
		if pieChartHeight > remainingHeight {
			pdf.AddPage()
			pdf.Ln(5)
		}

		pdf.SetFont("Arial", "B", 14)
		pdf.CellFormat(0, 10, "Summary Information Distribution", "", 1, "C", false, 0, "")
		pdf.Ln(5)

		pieChartX := (pageWidth - pieChartWidth) / 2
		pieChartY := pdf.GetY() + 2

		pdf.ImageOptions(
			pieChartFileName,
			pieChartX,
			pieChartY,
			pieChartWidth,
			pieChartHeight,
			false,
			gofpdf.ImageOptions{ImageType: "PNG"},
			0,
			"",
		)
		pdf.SetY(pieChartY + pieChartHeight + 2)

	} else {
		pdf.SetFont("Arial", "I", 12)
		pdf.CellFormat(0, 10, "No data available for pie chart representation.", "", 1, "C", false, 0, "")
		pdf.Ln(10)
	}

	// Generate PDF output
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("error generating PDF: %v", err)
	}

	return buf.Bytes(), nil
}

func hasValidData(data []chart.Value) bool {
	for _, d := range data {
		if d.Value > 0 {
			return true
		}
	}
	return false
}

func (ad *AdminUseCase) BestSellingProduct() ([]models.BestSellingProduct, error) {
	bestSellingProduct, err := ad.adminrepository.BestSellingProduct()
	if err != nil {
		return []models.BestSellingProduct{}, err
	}
	return bestSellingProduct, nil
}

func (ad *AdminUseCase) BestSellingCategory() ([]models.BestSellingCategory, error) {
	bestSellingCategory, err := ad.adminrepository.BestSellingCategory()
	if err != nil {
		return []models.BestSellingCategory{}, err
	}
	return bestSellingCategory, nil
}
