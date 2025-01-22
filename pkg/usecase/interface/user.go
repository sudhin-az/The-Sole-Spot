package interfaces

import (
	"ecommerce_clean_architecture/pkg/domain"
	"ecommerce_clean_architecture/pkg/utils/models"
	"time"
)

type UserUseCaseInterface interface {
	IsEmailExists(email string) bool
	IsPhoneExists(phone string) bool
	UserSignup(user models.User) (models.TokenUsers, error)
	VerifyOTP(email string, verify models.VerifyOTP) error
	SaveTempUserAndGenerateOTP(user models.User) (models.TokenUsers, error)
	generateAndSaveOTP(email string) (string, time.Time, error)
	VerifyOTPAndRegisterUser(email string, otp string) (models.TokenUsers, error)
	ResendOTP(string) error
	UserLogin(user models.User, input models.User) (models.TokenUsers, models.User, error)
	GetProducts() ([]models.ProductResponse, error)
	ListCategory() ([]domain.Category, error)
	UserProfile(userID string) (*models.User, error)
	UpdateProfile(editProfile models.User) (*models.User, error)
	ForgotPassword(userID int, input models.NewPassword) (models.User, error)
	AddAddress(userID int, address models.AddAddress) (models.AddAddress, error)
	UpdateAddress(userID int, address domain.Address) (domain.Address, error)
	DeleteAddress(userID int) error
	GetAllAddresses(id int) ([]domain.Address, error)
}
