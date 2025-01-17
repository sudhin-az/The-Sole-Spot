package db

import (
	"ecommerce_clean_architecture/pkg/config"
	"ecommerce_clean_architecture/pkg/domain"
	"ecommerce_clean_architecture/pkg/utils/models"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(cfg config.Config) (*gorm.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s", cfg.DBHost, cfg.DBUser, cfg.DBName, cfg.DBPort, cfg.DBPassword)
	db, dbErr := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	db.AutoMigrate(&domain.AdminDetails{})
	db.AutoMigrate(&domain.Address{})
	db.AutoMigrate(&domain.Category{})
	db.AutoMigrate(&domain.Products{})
	db.AutoMigrate(&domain.Cart{})
	db.AutoMigrate(&domain.Order{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.OTP{})
	db.AutoMigrate(&models.TempUser{})
	db.AutoMigrate(&domain.Coupon{})
	db.AutoMigrate(&domain.UsedCoupon{})

	return db, dbErr

}
