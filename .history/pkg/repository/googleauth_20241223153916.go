package repository

import (
	"ecommerce_clean_architecture/pkg/domain"

	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) GetUserByEmail(email string) (domain.Users, error) {
	var user domain.Users
	err := r.db.Where("email = ? AND deleted_at IS NULL", email).First(&user).Error
	return user, err
}

func (r *AuthRepository) CreateUser(newuser domain.Users) error {
	return r.db.Create(&newuser).Error
}
