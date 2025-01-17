package repository

import (
	"ecommerce_clean_architecture/pkg/utils"
	"ecommerce_clean_architecture/pkg/utils/models"
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserByEmail(email string) (models.UserSignUp, error) {
	var user models.UserSignUp
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return models.UserSignUp{}, err
	}
	return user, nil
}
func (r *UserRepository) SaveOTP(email, otp string, expiry time.Time) error {
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
	// Update the existing OTP record
	otpRecord.Email = email
	otpRecord.OTP = otp
	otpRecord.OtpExpiry = otpExpiry
	log.Println("OTP saved successfully")
	return r.db.Save(&otpRecord).Error
}

func (r *UserRepository) CreateUser(user models.TempUser) error {
	return r.db.Save(&user).Error
}

func (r *UserRepository) SaveTempUserAndGenerateOTP(user models.UserSignUp) error {

	if err := r.db.Create(&user).Error; err != nil {
		return fmt.Errorf("failed to save temporary user: %w", err)
	}
	otp := utils.GenerateOTP()
	otpExpiry := time.Now().Add(3 * time.Minute)
	fmt.Println("hiiiiiiii", user)
	if err := r.SaveOrUpdateOTP(user.Email, otp, otpExpiry); err != nil {
		return fmt.Errorf("failed to save OTP: %w", err)
	}
	return nil
}

func (r *UserRepository) VerifyAndMoveUser(tempUser models.UserSignUp) error {
	// Create a permanent user using the data from tempUser
	permanentUser := models.UserSignUp{
		Email:    tempUser.Email,
		Password: tempUser.Password,
	}
	// Save the permanent user to the database
	if err := r.db.Create(&permanentUser).Error; err != nil {
		return fmt.Errorf("failed to create permanent user: %w", err)
	}
	// Delete the temporary user after successfully creating the permanent user
	if err := r.db.Delete(&tempUser).Error; err != nil {
		return fmt.Errorf("failed to delete temporary user: %w", err)
	}
	return nil
}

// SaveTempUser saves user data to a temporary table
func (r *UserRepository) SaveTempUser(user models.UserSignUp) error {
	tempUser := &models.TempUser{FirstName: user.FirstName, LastName: user.LastName, Email: user.Email,
		Password: user.Password, Phone: user.Phone}
	return r.db.Create(&tempUser).Error
}

func (r *UserRepository) UpdateOTP(otp models.OTP) error {
	// Ensure the WHERE clause is correctly formed to target the correct record
	result := r.db.Model(&models.OTP{}).
		Where("email = ?", otp.Email). // This is the WHERE clause
		Updates(map[string]interface{}{
			"otp":        otp.OTP,
			"otp_expiry": otp.OtpExpiry,
		})

	// Check if the update affected any rows
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
	result := r.db.Where("LOWER(email) = LOWER(?)", email).Order("created_at desc").First(&otp)
	if result.Error != nil {
		return models.OTP{}, result.Error
	}
	return otp, nil
}

func (r *UserRepository) GetTempUserByEmail(email string) (models.TempUser, error) {
	fmt.Println("Email being queried:", email)
	var user models.TempUser
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return models.TempUser{}, err
	}
	return user, nil
}

// DeleteTempUser removes a user from the temporary table
func (r *UserRepository) DeleteTempUser(email string) error {
	return r.db.Table("temp_users").Where("email = ?", email).Delete(&models.UserSignUp{}).Error
}

func (r *UserRepository) GetOTP(email string) (string, time.Time, error) {
	// Log before querying
	log.Printf("Fetching OTP for email: %s", email)

	var otpRecord models.OTP
	err := r.db.Where("email = ?", email).First(&otpRecord).Error
	fmt.Println("email", email)
	if err != nil {
		log.Printf("Error fetching OTP for email %s: %s", email, err.Error())
		return "", time.Time{}, err
	}

	log.Printf("Fetched OTP: %s, Expiry: %s for email: %s", otpRecord.OTP, otpRecord.OtpExpiry.String(), email)
	return otpRecord.OTP, otpRecord.OtpExpiry, nil
}
func (r *UserRepository) DeleteOTP(email string) error {
	return r.db.Where("email = ?", email).Delete(&models.OTP{}).Error
}

func (r *UserRepository) IsEmailExists(email string) bool {
	var count int64
	r.db.Model(&models.UserSignUp{}).Where("email = ?", email).Count(&count)
	return count > 0
}

func (r *UserRepository) IsPhoneExists(phone string) bool {
	var count int64
	r.db.Model(&models.UserSignUp{}).Where("phone = ?", phone).Count(&count)
	return count > 0
}

// ResendOTP generates a new OTP, updates it in the database, and sends the OTP via email.
func (r *UserRepository) ResendOTP(email string) error {
	var otpRecord models.OTP

	// Check if the OTP record exists for the given email
	err := r.db.Where("email = ?", email).First(&otpRecord).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("no OTP record found for the provided email")
		}
		return err
	}

	// Generate a new OTP and update the expiry time
	newOTP := utils.GenerateOTP()
	otpRecord.OTP = newOTP
	otpRecord.OtpExpiry = time.Now().Add(2 * time.Minute)

	// Update the OTP record in the database
	err = r.db.Save(&otpRecord).Error
	if err != nil {
		return err
	}

	// Send the new OTP via email
	err = utils.SendOTPEmail(email, newOTP)
	if err != nil {
		return err
	}
	return nil
}
