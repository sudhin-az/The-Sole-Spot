package repository

import "gorm.io/gorm"

type ReviewRepository struct {
	DB *gorm.DB
}
