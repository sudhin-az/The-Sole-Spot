package interfaces

import "ecommerce_clean_architecture/pkg/utils/models"

type UserUseCaseInterface interface {
	IsEmailExists(email string) bool
	IsPhoneExists(phone string) bool
	UserSignUp(user models.UserSignUp) (models.TokenUsers, error)
	VerifyOTP(verify models.VerifyOTP) error
	ResendOTP(string) error
	UserLogin(user models.UserLogin) (models.TokenUsers, error)
}
