package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        int    `gorm:"primarykey"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Phone     string `json:"phone" validate:"required,e164"`
	Password  string `json:"password" validate:"required,min=8"`
	Blocked   bool
}
type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=5"`
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
type AddAddress struct {
	UserID    int    `json:"user_id"`
	HouseName string `json:"house_name" validate:"required,min=3,max=100"`
	Street    string `json:"street" validate:"required,min=3,max=100"`
	City      string `json:"city" validate:"required,min=2,max=50"`
	District  string `json:"district" validate:"required,min=3,max=50"`
	State     string `json:"state" validate:"required,len=3"`
	Pin       string `json:"pin" validate:"required,len=6"`
}

type Address struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"not null"`
	HouseName string `gorm:"not null"`
	Street    string `gorm:"not null"`
	City      string `gorm:"not null"`
	District  string `gorm:"not null"`
	State     string `gorm:"not null"`
	Pin       string `gorm:"not null"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
}

type NewPassword struct {
	Password    string `json:"password" validate:"required,min=8,max=32"`
	NewPassword string `json:"newpassword" validate:"required,min=8,max=32"`
	ReEnter     string `json:"reenter" validate:"required,min=8,max=32"`
}

type PaymentDetails struct {
	ID           uint   `json:"id"`
	Payment_Name string `json:"payment_name"`
}
type CheckoutDetails struct {
	Address        []Address
	Payment_Method []PaymentDetails
	Cart           []Cart
	Grand_Total    float64
	Total_Price    float64
}
