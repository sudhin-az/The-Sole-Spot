package usecase

import (
	"ecommerce_clean_architecture/pkg/helper"
	"ecommerce_clean_architecture/pkg/repository/interfaces"
	"ecommerce_clean_architecture/pkg/utils"
	"ecommerce_clean_architecture/pkg/utils/models"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	userRepo interfaces.UserRepository
}

func NewUserUseCase(userRepo interfaces.UserRepository) *UserUseCase {
	return &UserUseCase{userRepo: userRepo}
}

// IsEmailExists checks if the email already exists in the database
func (uc *UserUseCase) IsEmailExists(email string) bool {
	return uc.userRepo.IsEmailExists(email)
}

// IsPhoneExists checks if the phone number already exists in the database
func (uc *UserUseCase) IsPhoneExists(phone string) bool {
	return uc.userRepo.IsPhoneExists(phone)
}

func (uc *UserUseCase) UserSignUp(user models.UserSignUp) (models.TokenUsers, error) {
	// Check if the user already exists in the temp table
	if uc.userRepo.IsEmailExists(user.Email) || uc.userRepo.IsPhoneExists(user.Phone) {
		return models.TokenUsers{}, errors.New("user already exists")
	}
	// Hash the user's password
	fmt.Println("email", user.Email)
	hashedPassword, err := helper.HashPassword(user.Password)
	if err != nil {
		return models.TokenUsers{}, err
	}
	user.Password = hashedPassword

	// Generate OTP and save to database with expiry time
	otp := utils.GenerateOTP()
	otpExpiry := time.Now().Add(3 * time.Minute)

	// Save OTP
	fmt.Println("email", user.Email)
	err = uc.userRepo.SaveOrUpdateOTP(user.Email, otp, otpExpiry)
	fmt.Println("errrrrrrrrr", err)
	if err != nil {
		return models.TokenUsers{}, err
	}

	// Save user data temporarily
	fmt.Println("email", user.Email)
	err = uc.userRepo.SaveTempUser(user)
	if err != nil {
		// Cleanup OTP if saving temp user fails
		uc.userRepo.DeleteOTP(user.Email)
		return models.TokenUsers{}, err
	}

	tokenusers, err := uc.userRepo.GetTempUserByEmail(user.Email)
	if err != nil {
		// Cleanup OTP if saving temp user fails
		uc.userRepo.DeleteOTP(user.Email)
		return models.TokenUsers{}, err
	}

	// Send OTP to the user's email
	err = utils.SendOTPEmail(user.Email, otp)
	if err != nil {
		// Cleanup temporary user and OTP if sending email fails
		uc.userRepo.DeleteTempUser(user.Email)
		uc.userRepo.DeleteOTP(user.Email)
		return models.TokenUsers{}, err
	}

	userDetailsResponse := models.UserDetailsResponse{
		Email:     tokenusers.Email,
		FirstName: tokenusers.FirstName,
		LastName:  tokenusers.LastName,
		Phone:     tokenusers.Phone,
		Password:  tokenusers.Password,
	}

	fmt.Println("Token users", tokenusers)
	return models.TokenUsers{
		Users: userDetailsResponse,
	}, nil
}

func (uc *UserUseCase) VerifyOTP(email string, verify models.VerifyOTP) error {
	// Fetch OTP and expiry time from the database using the email
	otp, otpExpiry, err := uc.userRepo.GetOTP(email)
	if err != nil {
		return err
	}

	if otp != verify.OTP {
		fmt.Println(otp)
		return errors.New("invalid OTP")
	}

	if time.Now().After(otpExpiry) {
		fmt.Println(otp)
		return errors.New("expired OTP")
	}

	// Fetch user data from temporary storage
	tempUser, err := uc.userRepo.GetTempUserByEmail(email)
	if err != nil {
		return err
	}

	// Ensure user data does not already exist in the main table
	if uc.userRepo.IsEmailExists(tempUser.Email) || uc.userRepo.IsPhoneExists(tempUser.Phone) {
		// Clean up temporary user and OTP records
		uc.userRepo.DeleteTempUser(email)
		uc.userRepo.DeleteOTP(email)
		return errors.New("user already exists")
	}

	// Move user data from temporary table to main table

	err = uc.userRepo.CreateUser(tempUser)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	// Clean up OTP and temporary user record
	err = uc.userRepo.DeleteOTP(email)
	if err != nil {
		return err
	}
	err = uc.userRepo.DeleteTempUser(email)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UserUseCase) SaveTempUserAndGenerateOTP(user models.UserSignUp) (models.TokenUsers, error) {
	// Check if the user already exists in the main table
	if uc.userRepo.IsEmailExists(user.Email) || uc.userRepo.IsPhoneExists(user.Phone) {
		return models.TokenUsers{}, errors.New("user already exists")
	}

	// Hash the user's password
	hashedPassword, err := helper.HashPassword(user.Password)
	if err != nil {
		return models.TokenUsers{}, err
	}
	user.Password = hashedPassword

	// Save the user data temporarily
	err = uc.userRepo.SaveTempUser(user)
	if err != nil {
		return models.TokenUsers{}, err
	}

	// Generate and save OTP
	otp, otpExpiry, err := uc.generateAndSaveOTP(user.Email)
	if err != nil {
		return models.TokenUsers{}, err
	}
	fmt.Println("OTP Expiry time:", otpExpiry)

	// Send OTP to the user's email
	err = utils.SendOTPEmail(user.Email, otp)
	if err != nil {
		uc.userRepo.DeleteTempUser(user.Email)
		uc.userRepo.DeleteOTP(user.Email)
		return models.TokenUsers{}, err
	}

	tokenusers, err := uc.userRepo.GetTempUserByEmail(user.Email)
	if err != nil {
		// Cleanup OTP if saving temp user fails
		uc.userRepo.DeleteOTP(user.Email)
		return models.TokenUsers{}, err
	}

	userDetailsResponse := models.UserDetailsResponse{
		Email:     tokenusers.Email,
		FirstName: tokenusers.FirstName,
		LastName:  tokenusers.LastName,
		Phone:     tokenusers.Phone,
		Password:  tokenusers.Password,
	}

	return models.TokenUsers{
		Users: userDetailsResponse,
	}, nil
}

func (uc *UserUseCase) generateAndSaveOTP(email string) (string, time.Time, error) {
	otp := utils.GenerateOTP()
	otpExpiry := time.Now().Add(3 * time.Minute)

	err := uc.userRepo.SaveOrUpdateOTP(email, otp, otpExpiry)
	if err != nil {
		return "", time.Time{}, err
	}

	return otp, otpExpiry, nil
}

func (uc *UserUseCase) VerifyOTPAndRegisterUser(email string, otp string) (models.TokenUsers, error) {
	// Verify the OTP
	err := uc.VerifyOTP(email, models.VerifyOTP{OTP: otp})
	fmt.Println(err)
	if err != nil {
		return models.TokenUsers{}, errors.New("OTP verification failed")
	}

	// Retrieve temporary user data
	tempUser, err := uc.userRepo.GetTempUserByEmail(email)
	if err != nil {
		return models.TokenUsers{}, errors.New("temporary user not found")
	}

	// Move data to main user table
	err = uc.userRepo.CreateUser(tempUser)
	if err != nil {
		return models.TokenUsers{}, err
	}

	// Delete temporary user data and OTP after successful registration
	err = uc.userRepo.DeleteTempUser(email)
	if err != nil {
		return models.TokenUsers{}, err
	}
	err = uc.userRepo.DeleteOTP(email)
	if err != nil {
		return models.TokenUsers{}, err
	}

	// Generate and return JWT token or any response needed
	token, err := helper.GenerateTokenUsers(tempUser.ID, tempUser.Email, time.Now())
	if err != nil {
		return models.TokenUsers{}, err
	}

	return models.TokenUsers{AccessToken: token}, nil
}

func (uc *UserUseCase) ResendOTP(email string) error {
	otp := utils.GenerateOTP()
	otpExpiry := time.Now().Add(5 * time.Minute)

	err := uc.userRepo.UpdateOTP(models.OTP{
		Email:     email,
		OTP:       otp,
		OtpExpiry: otpExpiry,
	})
	if err != nil {
		return err
	}

	return utils.SendOTPEmail(email, otp)
}

func (uc *UserUseCase) UserLogin(user models.UserLogin) (models.TokenUsers, error) {
	userDetails, err := uc.userRepo.GetUserByEmail(user.Email)
	if err != nil {
		return models.TokenUsers{}, errors.New("user does not exist")
	}

	err = bcrypt.CompareHashAndPassword([]byte(userDetails.Password), []byte(user.Password))
	if err != nil {
		return models.TokenUsers{}, errors.New("incorrect password")
	}

	// Convert UserSignUp to UserDetailsResponse
	userDetailsResponse := models.UserDetailsResponse{
		Email:     userDetails.Email,
		FirstName: userDetails.FirstName,
		LastName:  userDetails.LastName,
		Phone:     userDetails.Phone,
		Password:  userDetails.Password,
	}

	// Generate access and refresh tokens
	accessToken, err := helper.GenerateAccessToken(userDetailsResponse)
	if err != nil {
		return models.TokenUsers{}, errors.New("failed to generate access token")
	}

	refreshToken, err := helper.GenerateRefreshToken(userDetailsResponse)
	if err != nil {
		return models.TokenUsers{}, errors.New("failed to generate refresh token")
	}

	return models.TokenUsers{
		Users:        userDetailsResponse,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
