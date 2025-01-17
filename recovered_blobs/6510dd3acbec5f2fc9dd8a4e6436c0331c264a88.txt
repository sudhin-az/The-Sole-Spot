package repository

import (
	"ecommerce_clean_architecture/pkg/domain"
	"ecommerce_clean_architecture/pkg/utils/models"
	"fmt"

	"gorm.io/gorm"
)

type AdminRepository struct {
	DB *gorm.DB
}

func NewAdminRepository(DB *gorm.DB) *AdminRepository {
	return &AdminRepository{
		DB: DB,
	}
}

func (ad *AdminRepository) CheckAdminAvailability(admin models.AdminSignUp) bool {
	var count int

	if err := ad.DB.Raw("select count(*) from admin_details where email = ?", admin.Email).Scan(&count).Error; err != nil {
		return false
	}

	return count > 0
}

func (ad *AdminRepository) SignUpHandler(admin models.AdminSignUp) (models.AdminDetailsResponse, error) {
	var adminDetails models.AdminDetailsResponse

	if err := ad.DB.Raw("insert into admin_details(name, email, password) values(?, ?, ?) returning id, name, email", admin.Name, admin.Email, admin.Password).Scan(&adminDetails).Error; err != nil {
		return models.AdminDetailsResponse{}, err
	}
	return adminDetails, nil
}

func (ad *AdminRepository) LoginHandler(admin models.AdminLogin) (domain.AdminDetails, error) {

	var adminCompareDetails domain.AdminDetails
	if err := ad.DB.Raw("select * from admin_details where email = ?", admin.Email).Scan(&adminCompareDetails).Error; err != nil {
		return domain.AdminDetails{}, err
	}
	return adminCompareDetails, nil
}

func (ad *AdminRepository) GetUsers() ([]models.User, error) {
	var listofusers []models.User
	err := ad.DB.Raw("SELECT * FROM user_sign_ups").Scan(&listofusers).Error
	if err != nil {
		return nil, err
	}
	return listofusers, nil
}

func (ad *AdminRepository) GetUserByID(userID int) (models.User, error) {
	querry := fmt.Sprintf("SELECT * FROM user_sign_ups WHERE id = '%d'", userID)
	var userDetails models.User
	if err := ad.DB.Raw(querry).Scan(&userDetails).Error; err != nil {
		return models.User{}, err
	}
	return userDetails, nil
}

func (ad *AdminRepository) UpdateBlockUserByID(user models.User) error {
	err := ad.DB.Exec("UPDATE user_sign_ups SET blocked = ? WHERE id = ?", user.Blocked, user.ID).Error
	if err != nil {
		return err
	}
	return nil
}
