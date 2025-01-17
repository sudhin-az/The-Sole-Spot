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
	err := ad.DB.Raw("SELECT * FROM users").Scan(&listofusers).Error
	if err != nil {
		return nil, err
	}
	return listofusers, nil
}

func (ad *AdminRepository) GetUserByID(userID int) (models.User, error) {
	querry := fmt.Sprintf("SELECT * FROM users WHERE id = '%d'", userID)
	var userDetails models.User
	if err := ad.DB.Raw(querry).Scan(&userDetails).Error; err != nil {
		return models.User{}, err
	}
	return userDetails, nil
}

func (ad *AdminRepository) UpdateBlockUserByID(user models.User) error {
	err := ad.DB.Exec("UPDATE users SET blocked = ? WHERE id = ?", user.Blocked, user.ID).Error
	if err != nil {
		return err
	}
	return nil
}

func (o *OrderRepository) BeginTransaction() (*gorm.DB, error) {
	tx := o.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (o *OrderRepository) CommitTransaction(tx *gorm.DB) error {
	return tx.Commit().Error
}

func (o *OrderRepository) RollbackTransaction(tx *gorm.DB) error {
	return tx.Rollback().Error
}
func (o *OrderRepository) GetProductStock(productID int) (int, error) {
	var stock int
	err := o.DB.Raw("select stock from products where id = ?", productID).Scan(&stock).Error
	if err != nil {
		return 0, err
	}
	return stock, nil
}
func (o *OrderRepository) UpdateProductStock(tx *gorm.DB, productID int, newStock int) error {
	return tx.Model(&domain.Products{}).Where("id = ?", productID).Update("stock", newStock).Error
}
