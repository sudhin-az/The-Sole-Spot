package usecase

import (
	"ecommerce_clean_architecture/pkg/domain"
	"ecommerce_clean_architecture/pkg/helper"
	"ecommerce_clean_architecture/pkg/repository"
	"ecommerce_clean_architecture/pkg/utils/models"
	"errors"
	"fmt"
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

	// Convert UserSignUp to UserDetailsResponse
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
	// Validate if orderID is an integer
	orderIDInt, err := strconv.Atoi(orderID)
	if err != nil {
		return fmt.Errorf("invalid order ID format: %w", err)
	}

	userTest, err := o.orderRepository.UserOrderRelationship(orderIDInt, userID)
	if err != nil {
		return err
	}
	if userTest != userID {
		return errors.New("the order is not done by this user")
	}

	// Begin a transaction
	tx, err := o.orderRepository.BeginTransaction()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer o.orderRepository.RollbackTransaction(tx)

	orderProductDetails, err := o.orderRepository.GetProductDetailsFromOrders(orderIDInt)
	if err != nil {
		return err
	}
	shipmentStatus, err := o.orderRepository.GetShipmentStatus(orderIDInt)
	if err != nil {
		return err
	}
	if shipmentStatus == "delivered" {
		return errors.New("items already delivered, cannot cancel")
	}

	if shipmentStatus == "returned" || shipmentStatus == "Failed" {
		return fmt.Errorf("the order is in %s, so no point in cancelling", shipmentStatus)
	}
	if shipmentStatus == "cancelled" {
		return errors.New("the order is already cancelled, so no point in cancelling")
	}

	err = o.orderRepository.CancelOrders(orderIDInt)
	if err != nil {
		return err
	}

	for _, product := range orderProductDetails {
		availableStock, err := o.orderRepository.GetProductStock(product.ProductID)
		if err != nil {
			return err
		}

		// Restore stock
		newStock := availableStock + product.Quantity
		err = o.orderRepository.UpdateProductStock(tx, product.ProductID, newStock)
		if err != nil {
			return errors.New("failed to restore product stock")
		}
	}

	// Commit the transaction
	err = o.orderRepository.CommitTransaction(tx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	err = o.orderRepository.UpdateQuantityOfProduct(orderProductDetails)
	if err != nil {
		return err
	}

	return nil
}
