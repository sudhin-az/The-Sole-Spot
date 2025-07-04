package db

import (
	"ecommerce_clean_arch/pkg/config"
	"ecommerce_clean_arch/pkg/domain"
	"ecommerce_clean_arch/pkg/helper"
	"ecommerce_clean_arch/pkg/utils/models"
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

	log.Println("✅ Database migrated successfully!")

	// ✅ Insert default admin if not exists
	var count int64
	db.Model(&domain.AdminDetails{}).Where("email = ?", "sudhin@gmail.com").Count(&count)

	if count == 0 {
		// IMPORTANT: use a password with at least 8 chars if your validation needs it
		password := "sudhin123"
		hashedPassword, hashErr := helper.HashPassword(password)
		if hashErr != nil {
			log.Fatalf("failed to hash password: %v", hashErr)
		}

		admin := domain.AdminDetails{
			Name:     "Sudhin",
			Phone:    "1234567890",
			Email:    "sudhin@gmail.com",
			Password: hashedPassword,
		}

		if err := db.Create(&admin).Error; err != nil {
			log.Fatalf("failed to insert admin: %v", err)
		}

		log.Println("✅ Default admin inserted.")
	} else {
		log.Println("✅ Admin already exists, skipping insert.")
	}

	return db, nil
}

