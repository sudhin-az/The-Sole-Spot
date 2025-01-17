package usecase

import (
	"ecommerce_clean_architecture/pkg/domain"
	"ecommerce_clean_architecture/pkg/helper"
	"ecommerce_clean_architecture/pkg/repository"
	"ecommerce_clean_architecture/pkg/utils/models"
	"errors"

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

	// Convert User to UserDetailsResponse
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
