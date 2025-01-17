package interfaces

import (
	"ecommerce_clean_architecture/pkg/utils/models"
	"time"
)

type UserRepository interface {
	IsEmailExists(email string) bool
	IsPhoneExists(phone string) bool
	SaveTempUser(user models.UserSignUp) error
	GetTempUserByEmail(email string) (models.TempUser, error)
	DeleteTempUser(email string) error
	SaveOrUpdateOTP(email string, otp string, otpExpiry time.Time) error
	GetOTP(email string) (string, time.Time, error)
	SaveOTP(email, otp string, expiry time.Time) error
	DeleteOTP(email string) error
	UpdateOTP(otp models.OTP) error
	CreateUser(user models.TempUser) error
	GetUserByEmail(email string) (models.UserSignUp, error)
}
