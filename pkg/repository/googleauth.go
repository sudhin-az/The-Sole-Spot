package repository

import (
	"ecommerce_clean_architecture/pkg/domain"
	"ecommerce_clean_architecture/pkg/repository/interfaces"

	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) interfaces.AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) GetUserByEmail(email string) (domain.Users, error) {
	var user domain.Users
	err := r.db.Where("email = ? AND deleted_at IS NULL", email).First(&user).Error
	return user, err
}

func (r *AuthRepository) CreateUser(user domain.Users) error {
	return r.db.Create(&user).Error
}
