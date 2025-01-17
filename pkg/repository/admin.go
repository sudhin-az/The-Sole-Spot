package repository

import (
	"ecommerce_clean_architecture/pkg/domain"
	"ecommerce_clean_architecture/pkg/utils/models"

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

func (ad *AdminRepository) GetUsers(listusers models.UserSignUp) (models.UserSignUp, error) {
	var users domain.Users
	err := ad.DB.Raw("SELECT * FROM users", users).Error
	if err != nil {
		return models.UserSignUp{}, err
	}
	return models.UserSignUp{}, nil
}
