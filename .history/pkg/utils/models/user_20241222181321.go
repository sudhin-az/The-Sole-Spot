package models

import "time"

type UserSignUp struct {
	ID              uint   `gorm:"primarykey"`
	FirstName       string `json:"first_name" validate:"required"`
	LastName        string `json:"last_name" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Phone           string `json:"phone" validate:"required,e164"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
	Blocked         bool
}

type UserLogin struct {
	Email    string `json:"email" binding:"required" validate:"email"`
	Password string `json:"password" binding:"required"`
}

type TempUser struct {
	ID        uint   `json:"id" gorm:"primary key,not null"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" validate:"email"`
	Password  string `json:"password" validate:"min=8,max=20"`
	Phone     string `json:"phone"`
	Blocked   bool
}

type OTP struct {
	ID        int    `json:"id" gorm:"primary key,not null"`
	Email     string `json:"email" validate:"email"`
	OTP       string `json:"otp" validate:"otp"`
	OtpExpiry time.Time
}

type VerifyOTP struct {
	OTP string `json:"otp"`
}

type TokenUsers struct {
	Users        UserDetailsResponse
	AccessToken  string
	RefreshToken string
}
type UserDetailsResponse struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
}
type UserSignInResponse struct {
	Id       int    `json:"id"`
	UserID   int    `json:"user_id"`
	Name     string `json:"name"`
	Email    string `json:"email" validate:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}
type UserDetailsAtAdmin struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	BlockStatus bool   `json:"block_status"`
}
