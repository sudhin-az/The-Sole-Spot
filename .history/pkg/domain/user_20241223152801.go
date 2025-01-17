package domain

import (
	"time"

	"gorm.io/gorm"
)

type Users struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Name     string `json:"name" gorm:"type:varchar(255)"`
	Email    string `json:"email" gorm:"type:varchar(255);unique"`
	Phone    string `json:"phone" gorm:"type:varchar(15);unique"`
	Password string `json:"password"`
}

type UserLoginMethod struct {
	ID                   uint   `gorm:"unique"`
	UserLoginMethodEmail string `json:"user_login_method_email" validate:"user_login_method_email"`
	LoginMethod          string
}

type GoogleResponse struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Picture string `json:"picture"`
}

type Address struct {
	ID        int    `json:"id" gorm:"primaryKey;not null"`
	UserID    int    `json:"user_id"`
	Users     Users  `json:"-" gorm:"foreignkey:UserID"`
	HouseName string `json:"house_name" validate:"required,min=3,max=100"`
	Street    string `json:"street" validate:"required,min=3,max=100"`
	City      string `json:"city" validate:"required,min=2,max=50"`
	District  string `json:"district" validate:"required,min=3,max=50"`
	State     string `json:"state" validate:"required,len=3"`
	Pin       string `json:"pin_code" validate:"required,len=6,numeric"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
}
