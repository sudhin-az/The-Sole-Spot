package domain

import "ecommerce_clean_architecture/pkg/utils/models"

type AdminDetails struct {
	ID       int    `json:"id" gorm:"primary key,not null"`
	Name     string `json:"name" gorm:"validate:required"`
	Phone    string `json:"phone" gorm:"validate:required"`
	Email    string `json:"email" gorm:"validate:required"`
	Password string `json:"password" gorm:"validate:required"`
}

type TokenAdmin struct {
	Admin models.AdminDetailsResponse
	Token string
}
