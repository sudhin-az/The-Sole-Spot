package usecase

import (
	"ecommerce_clean_architecture/pkg/domain"
	"ecommerce_clean_architecture/pkg/helper"
	"ecommerce_clean_architecture/pkg/repository"
	"ecommerce_clean_architecture/pkg/utils/models"
	"errors"
	"fmt"
	"log"
	"strconv"

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
		fmt.Println("Error in repository:", err)
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
	}
	updateOrder, err := ad.adminrepository.ChangeOrderStatus(orderID, status)
	if err != nil {
		return models.Order{}, fmt.Errorf("could not getting updated order status: %w", err)
	}
	return updateOrder, nil
}
