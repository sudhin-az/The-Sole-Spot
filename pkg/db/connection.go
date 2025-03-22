package db

import (
	"ecommerce_clean_architecture/pkg/config"
	"ecommerce_clean_architecture/pkg/domain"
	"ecommerce_clean_architecture/pkg/utils/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(cfg config.Config) (*gorm.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s", cfg.DBHost, cfg.DBUser, cfg.DBName, cfg.DBPort, cfg.DBPassword)
	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&domain.AdminDetails{},
		&domain.Address{},
		&domain.Category{},
		&domain.Products{},
		&domain.Review{},
		&domain.Cart{},
		&domain.PaymentMethod{},
		&domain.RazorPay{},
		&models.Order{},
		&domain.OrderItem{},
		&domain.Wallet{},
		&domain.WalletTransaction{},
		&domain.Coupons{},
		&domain.Wishlist{},
		&models.User{},
		&models.OTP{},
		&models.TempUser{},
	)
	if err != nil {
		return nil, err
	}

	log.Println("Database migrated successfully!")
	return db, nil
}
