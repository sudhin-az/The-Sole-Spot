package repository

import (
	"ecommerce_clean_architecture/pkg/domain"
	"ecommerce_clean_architecture/pkg/utils"
	"ecommerce_clean_architecture/pkg/utils/models"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
func (r *UserRepository) SaveOTP(email, otp string, expiry time.Time) error {
	fmt.Println("email", email)
	newOTP := models.OTP{
		Email:     email,
		OTP:       otp,
		OtpExpiry: expiry,
	}

	result := r.db.Create(&newOTP)
	return result.Error
}

func (r *UserRepository) SaveOrUpdateOTP(email string, otp string, otpExpiry time.Time) error {
	var otpRecord models.OTP
	log.Printf("Saving OTP for email: %s, OTP: %s, Expiry: %s", email, otp, otpExpiry.String())
	otpRecord.Email = email
	otpRecord.OTP = otp
	otpRecord.OtpExpiry = otpExpiry
	log.Println("OTP saved successfully")
	return r.db.Save(&otpRecord).Error
}

func (r *UserRepository) CreateUser(user models.User) error {
	return r.db.Save(&user).Error
}

func (r *UserRepository) SaveTempUserAndGenerateOTP(user models.User) error {

	if err := r.db.Create(&user).Error; err != nil {
		return fmt.Errorf("failed to save temporary user: %w", err)
	}
	otp := utils.GenerateOTP()
	otpExpiry := time.Now().Add(3 * time.Minute)
	if err := r.SaveOrUpdateOTP(user.Email, otp, otpExpiry); err != nil {
		return fmt.Errorf("failed to save OTP: %w", err)
	}
	return nil
}

func (r *UserRepository) VerifyOTPAndMoveUser(email string, otp string) error {
	var otpRecord models.OTP
	err := r.db.Where("email = ? AND otp = ?", email, otp).First(&otpRecord).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("invalid or expired OTP")
		}
		return err
	}
	if time.Now().After(otpRecord.OtpExpiry) {
		return fmt.Errorf("OTP has expired")
	}
	var tempUser models.TempUser
	err = r.db.Where("email = ?", email).First(&tempUser).Error
	if err != nil {
		return err
	}
	mainUser := models.TempUser{
		FirstName: tempUser.FirstName,
		LastName:  tempUser.LastName,
		Email:     tempUser.Email,
		Password:  tempUser.Password,
		Phone:     tempUser.Phone,
	}
	err = r.db.Create(&mainUser).Error
	if err != nil {
		return err
	}
	r.db.Delete(&otpRecord)

	return nil
}

func (r *UserRepository) SaveTempUser(user models.User) error {
	tempUser := &models.TempUser{FirstName: user.FirstName, LastName: user.LastName, Email: user.Email,
		Password: user.Password, Phone: user.Phone}
	return r.db.Create(&tempUser).Error
}

func (r *UserRepository) UpdateOTP(otp models.OTP) error {
	result := r.db.Model(&models.OTP{}).
		Where("email = ?", otp.Email).
		Updates(map[string]interface{}{
			"otp":        otp.OTP,
			"otp_expiry": otp.OtpExpiry,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows were updated, check your WHERE conditions")
	}
	return result.Error
}

func (r *UserRepository) GetOTPByEmail(email string) (models.OTP, error) {
	var otp models.OTP
	result := r.db.Where("email = ?", email).Order("created_at desc").First(&otp)
	if result.Error != nil {
		return models.OTP{}, result.Error
	}
	return otp, nil
}

func (r *UserRepository) GetTempUserByEmail(email string) (models.TempUser, error) {
	fmt.Println("Email being queried:", email)
	var user models.TempUser
	email = strings.ToLower(strings.TrimSpace(email))
	err := r.db.Where("LOWER(email) = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("Record not found for email", email)
			return models.TempUser{}, fmt.Errorf("temporary user not found for email %s", email)
		}
		fmt.Println("Error querrying temp_users table:", err)
		return models.TempUser{}, err
	}
	return user, nil
}
func (r *UserRepository) DeleteTempUser(email string) error {
	return r.db.Where("email = ?", email).Delete(&models.TempUser{}).Error
}

func (r *UserRepository) GetOTP(email string) (string, time.Time, error) {

	email = strings.TrimSpace(email)

	log.Printf("Fetching OTP for email: %s", email)

	var otpRecord models.OTP
	err := r.db.Where("email = ?", email).First(&otpRecord).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("NO OTP found for email %s", email)
			return "", time.Time{}, fmt.Errorf("no OTP found for email: %s", email)
		}
		log.Printf("Error fetching OTP for email %s: %s", email, err.Error())
		return "", time.Time{}, err
	}
	if time.Now().After(otpRecord.OtpExpiry) {
		log.Printf("OTP for email %s has expired", email)
		return "", time.Time{}, fmt.Errorf("OTP has expired")
	}

	log.Printf("Fetched OTP: %s, Expiry: %s for email: %s", otpRecord.OTP, otpRecord.OtpExpiry.String(), email)
	return otpRecord.OTP, otpRecord.OtpExpiry, nil
}
func (r *UserRepository) DeleteOTP(email string) error {
	return r.db.Where("email = ?", email).Delete(&models.OTP{}).Error
}

func (r *UserRepository) IsEmailExists(email string) bool {
	var count int64
	r.db.Model(&models.User{}).Where("email = ?", email).Count(&count)
	return count > 0
}

func (r *UserRepository) IsPhoneExists(phone string) bool {
	var count int64
	r.db.Model(&models.User{}).Where("phone = ?", phone).Count(&count)
	return count > 0
}

func (r *UserRepository) ResendOTP(email string) error {
	var otpRecord models.OTP
	err := r.db.Where("email = ?", email).First(&otpRecord).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("no OTP record found for the provided email")
		}
		return err
	}
	newOTP := utils.GenerateOTP()
	otpRecord.OTP = newOTP
	otpRecord.OtpExpiry = time.Now().Add(3 * time.Minute)

	err = r.db.Save(&otpRecord).Error
	if err != nil {
		return err
	}
	err = utils.SendOTPEmail(email, newOTP)
	if err != nil {
		return err
	}
	return nil
}
func (r *UserRepository) UnblockUser(email string) error {
	result := r.db.Model(&models.UserLogin{}).Where("email = ?", email).Update("blocked", false)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no user found with the given email")
	}
	return nil
}

func (r *UserRepository) GetProducts() ([]models.ProductResponse, error) {
	var listproducts []models.ProductResponse
	err := r.db.Raw("SELECT * FROM products WHERE deleted_at IS NULL").Scan(&listproducts).Error
	if err != nil {
		return []models.ProductResponse{}, err
	}
	return listproducts, nil
}
func (r *UserRepository) ListCategory() ([]domain.Category, error) {
	var listcategory []domain.Category
	err := r.db.Raw("SELECT * FROM categories WHERE deleted_at IS NULL").Scan(&listcategory).Error
	if err != nil {
		return []domain.Category{}, err
	}
	return listcategory, nil
}

func (r *UserRepository) UserProfile(userID int) (*models.User, error) {
	var userProfile models.User
	result := r.db.Where("id = ?", userID).First(&userProfile)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with ID %d not found", userID)
		}
		return nil, result.Error
	}

	return &userProfile, nil
}

func (r *UserRepository) UpdateProfile(editProfile models.User) (*models.User, error) {
	var profile models.User
	query := "UPDATE users set first_name = ?, last_name = ?, email = ?, phone = ? WHERE id = ? RETURNING *;"
	result := r.db.Raw(query, editProfile.FirstName, editProfile.LastName, editProfile.Email, editProfile.Phone, editProfile.ID).Scan(&profile)
	if result.Error != nil {
		return &models.User{}, errors.New("face some issue while update profile")
	}
	if result.RowsAffected == 0 {
		return &models.User{}, errors.New("No rows are affected")
	}
	return &profile, nil
}
func (r *UserRepository) GetUserByID(userID int) (models.User, error) {
	var user models.User

	err := r.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *UserRepository) UpdatePassword(userID int, newPassword string) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("password", newPassword).Error
}
func (r *UserRepository) GetPassword(userID int) (models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", userID).Select("password").First(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *UserRepository) AddAddress(userID int, address models.AddAddress) (models.AddAddress, error) {
	var userAddress models.AddAddress
	err := r.db.Raw(`
		INSERT INTO addresses (user_id, house_name, street, city, district, state, pin) 
		VALUES (?, ?, ?, ?, ?, ?, ?) 
		RETURNING user_id, house_name, street, city, district, state, pin`,
		userID, address.HouseName, address.Street, address.City, address.District, address.State, address.Pin).
		Scan(&userAddress).Error

	if err != nil {
		return models.AddAddress{}, err
	}
	return userAddress, nil
}
func (r *UserRepository) UpdateAddress(userID int, address domain.Address) (domain.Address, error) {
	var userAddress domain.Address
	query := `
		UPDATE addresses 
		SET house_name = ?, street = ?, city = ?, district = ?, state = ?, pin = ? 
		WHERE user_id = ? AND id = ? 
		RETURNING *;`
	result := r.db.Raw(query, address.HouseName, address.Street, address.City, address.District, address.State, address.Pin, userID, address.ID).Scan(&userAddress)

	if result.Error != nil {
		return domain.Address{}, errors.New("some issue occurred while updating the address")
	}
	if result.RowsAffected == 0 {
		return domain.Address{}, errors.New("no rows are affected")
	}
	return userAddress, nil
}

func (r *UserRepository) DeleteAddress(userID int) error {
	result := r.db.Where("id = ?", userID).Delete(&models.Address{})
	if result.RowsAffected < 1 {
		return errors.New("the id does not exist")
	}
	return result.Error
}
func (c *UserRepository) GetAllAddresses(id int) ([]domain.Address, error) {
	var addresses []domain.Address
	if err := c.db.Raw("SELECT * FROM addresses WHERE deleted_at IS NULL").Scan(&addresses).Error; err != nil {
		return []domain.Address{}, err
	}
	return addresses, nil
}
