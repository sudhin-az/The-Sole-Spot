package usecase

import (
	"ecommerce_clean_architecture/pkg/domain"
	"ecommerce_clean_architecture/pkg/helper"
	"ecommerce_clean_architecture/pkg/repository/interfaces"
	"ecommerce_clean_architecture/pkg/utils"
	"ecommerce_clean_architecture/pkg/utils/models"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	userRepo interfaces.UserRepository
}

func NewUserUseCase(userRepo interfaces.UserRepository) *UserUseCase {
	return &UserUseCase{userRepo: userRepo}
}

func (uc *UserUseCase) IsEmailExists(email string) bool {
	return uc.userRepo.IsEmailExists(email)
}

func (uc *UserUseCase) IsPhoneExists(phone string) bool {
	return uc.userRepo.IsPhoneExists(phone)
}

func (uc *UserUseCase) UserSignup(user models.User) (models.TokenUsers, error) {
	if uc.userRepo.IsEmailExists(user.Email) || uc.userRepo.IsPhoneExists(user.Phone) {
		return models.TokenUsers{}, errors.New("user already exists")
	}
	fmt.Println("email", user.Email)
	hashedPassword, err := helper.HashPassword(user.Password)
	if err != nil {
		return models.TokenUsers{}, err
	}
	user.Password = hashedPassword

	otp := utils.GenerateOTP()
	otpExpiry := time.Now().Add(3 * time.Minute)

	err = uc.userRepo.SaveOrUpdateOTP(user.Email, otp, otpExpiry)
	fmt.Println("errrrrrrrrr", err)
	if err != nil {
		return models.TokenUsers{}, err
	}

	err = uc.userRepo.SaveTempUser(user)
	if err != nil {
		return models.TokenUsers{}, err
	}

	tokenusers, err := uc.userRepo.GetTempUserByEmail(user.Email)
	if err != nil {
		uc.userRepo.DeleteOTP(user.Email)
		return models.TokenUsers{}, err
	}

	err = utils.SendOTPEmail(user.Email, otp)
	if err != nil {
		uc.userRepo.DeleteOTP(user.Email)
		return models.TokenUsers{}, err
	}

	fmt.Println("Token users", tokenusers)
	return models.TokenUsers{}, nil
}
func ValidateUser(user models.User) error {
	var validationErrors []string

	if matched, _ := regexp.MatchString(`^[a-zA-Z\s]+$`, user.FirstName); !matched || len(user.FirstName) < 2 {
		validationErrors = append(validationErrors, "First name must contain only letters and spaces, and be at least 2 characters long")
	}

	if matched, _ := regexp.MatchString(`^[a-zA-Z\s]+$`, user.LastName); !matched || len(user.LastName) < 2 {
		validationErrors = append(validationErrors, "Last name must contain only letters and spaces, and be at least 2 characters long")
	}

	if matched, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, user.Email); !matched {
		validationErrors = append(validationErrors, "Invalid email format")
	}

	if matched, _ := regexp.MatchString(`^[0-9]{10}$`, user.Phone); !matched {
		validationErrors = append(validationErrors, "Invalid phone number format")
	}

	if len(user.Password) < 5 {
		validationErrors = append(validationErrors, "Password must be at least 5 characters long")
	}

	if len(validationErrors) > 0 {
		return errors.New(strings.Join(validationErrors, "; "))
	}

	return nil
}
func ValidatePassword(password models.NewPassword) error {
	var validationErrors []string
	if len(password.Password) < 5 {
		validationErrors = append(validationErrors, "Password must be at least 5 characters long")
	}
	if password.NewPassword != password.ReEnter {
		return errors.New("password do not match")
	}
	if password.NewPassword == "" || password.ReEnter == "" {
		return errors.New("password cannot be empty")
	}
	return nil
}

func (uc *UserUseCase) VerifyOTP(email string, verify models.VerifyOTP) error {
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

	tempUser, err := uc.userRepo.GetTempUserByEmail(email)
	if err != nil {
		return err
	}

	if uc.userRepo.IsEmailExists(tempUser.Email) || uc.userRepo.IsPhoneExists(tempUser.Phone) {
		uc.userRepo.DeleteOTP(email)
		return errors.New("user already exists")
	}

	tempUser = models.TempUser{}

	err = uc.userRepo.DeleteOTP(email)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UserUseCase) SaveTempUserAndGenerateOTP(user models.User) (models.TokenUsers, error) {
	if uc.userRepo.IsEmailExists(user.Email) || uc.userRepo.IsPhoneExists(user.Phone) {
		return models.TokenUsers{}, errors.New("user already exists")
	}

	hashedPassword, err := helper.HashPassword(user.Password)
	if err != nil {
		return models.TokenUsers{}, err
	}
	user.Password = hashedPassword

	err = uc.userRepo.SaveTempUser(user)
	if err != nil {
		return models.TokenUsers{}, err
	}

	otp, otpExpiry, err := uc.generateAndSaveOTP(user.Email)
	if err != nil {
		return models.TokenUsers{}, err
	}
	fmt.Println("OTP Expiry time:", otpExpiry)

	err = utils.SendOTPEmail(user.Email, otp)
	if err != nil {
		uc.userRepo.DeleteTempUser(user.Email)
		uc.userRepo.DeleteOTP(user.Email)
		return models.TokenUsers{}, err
	}

	tokenusers, err := uc.userRepo.GetTempUserByEmail(user.Email)
	if err != nil {
		uc.userRepo.DeleteOTP(user.Email)
		return models.TokenUsers{}, err
	}

	userDetailsResponse := models.UserDetailsResponse{
		Id:        int(tokenusers.ID),
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
	err := uc.VerifyOTP(email, models.VerifyOTP{OTP: otp})
	fmt.Println(err)
	if err != nil {
		return models.TokenUsers{}, errors.New("OTP verification failed")
	}

	tempUser, err := uc.userRepo.GetTempUserByEmail(email)
	if err != nil {
		return models.TokenUsers{}, errors.New("temporary user not found")
	}

	User := models.User{
		ID:        int(tempUser.ID),
		FirstName: tempUser.FirstName,
		LastName:  tempUser.LastName,
		Email:     tempUser.Email,
		Phone:     tempUser.Phone,
		Password:  tempUser.Password,
	}
	err = uc.userRepo.CreateUser(User)
	if err != nil {
		return models.TokenUsers{}, err
	}

	err = uc.userRepo.DeleteTempUser(email)
	if err != nil {
		return models.TokenUsers{}, err
	}
	err = uc.userRepo.DeleteOTP(email)
	if err != nil {
		return models.TokenUsers{}, err
	}

	token, err := helper.GenerateTokenUsers(tempUser.ID, tempUser.Email, time.Now())
	if err != nil {
		return models.TokenUsers{}, err
	}

	return models.TokenUsers{AccessToken: token, RefreshToken: token, Users: models.UserDetailsResponse{
		Id:        int(tempUser.ID),
		FirstName: tempUser.FirstName,
		LastName:  tempUser.LastName,
		Email:     tempUser.Email,
		Phone:     tempUser.Phone,
		Password:  tempUser.Password,
	}}, nil
}

func (uc *UserUseCase) ResendOTP(email string) error {
	fmt.Println("Email:", email)
	otp := utils.GenerateOTP()
	otpExpiry := time.Now().Add(3 * time.Minute)
	fmt.Println("Email, Otp, OtpExpiry", email, otp, otpExpiry)
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

	if userDetails.Blocked {
		return models.TokenUsers{}, errors.New("user is blocked, so couldn't be logged in")
	}

	err = bcrypt.CompareHashAndPassword([]byte(userDetails.Password), []byte(user.Password))
	if err != nil {
		return models.TokenUsers{}, errors.New("incorrect password")
	}
	userDetailsResponse := models.UserDetailsResponse{
		Id:        int(userDetails.ID),
		Email:     userDetails.Email,
		FirstName: userDetails.FirstName,
		LastName:  userDetails.LastName,
		Phone:     userDetails.Phone,
		Password:  userDetails.Password,
	}

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

func (uc *UserUseCase) GetProducts() ([]models.ProductResponse, error) {
	productDetails, err := uc.userRepo.GetProducts()
	if err != nil {
		return []models.ProductResponse{}, err
	}
	return productDetails, nil
}
func (cat *UserUseCase) ListCategory() ([]domain.Category, error) {
	categoryDetails, err := cat.userRepo.ListCategory()
	if err != nil {
		return []domain.Category{}, err
	}
	return categoryDetails, nil
}

func (uc *UserUseCase) UserProfile(userID int) (*models.User, error) {
	userprofile, err := uc.userRepo.UserProfile(userID)
	if err != nil {
		return nil, fmt.Errorf("user with ID %d not found", userID)
	}
	return userprofile, nil
}

func (uc *UserUseCase) UpdateProfile(editProfile models.User) (*models.User, error) {
	userProfile, err := uc.userRepo.UpdateProfile(editProfile)
	if err != nil {
		return &models.User{}, err
	}
	return userProfile, nil
}

func (uc *UserUseCase) ForgotPassword(userID int, input models.NewPassword) (models.User, error) {
	user, err := uc.userRepo.GetUserByID(userID)
	if err != nil {
		return models.User{}, errors.New("failed to retrieve user")
	}
	password, err := uc.userRepo.GetPassword(user.ID)
	if err != nil {
		return models.User{}, errors.New("failed to retrieve password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(password.Password), []byte(input.Password)); err != nil {
		return models.User{}, errors.New("invalid password")
	}
	if input.NewPassword != input.ReEnter {
		return models.User{}, errors.New("passwords do not match")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, errors.New("failed to hash password")
	}
	if err := uc.userRepo.UpdatePassword(userID, string(hashedPassword)); err != nil {
		return models.User{}, errors.New("failed to update password")
	}
	return models.User{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
		Password:  user.Password,
		Blocked:   user.Blocked,
	}, nil
}

func (u *UserUseCase) AddAddress(userID int, address models.AddAddress) (models.AddAddress, error) {
	return u.userRepo.AddAddress(userID, address)
}
func (u *UserUseCase) UpdateAddress(userID int, address domain.Address) (domain.Address, error) {
	return u.userRepo.UpdateAddress(userID, address)
}
func (u *UserUseCase) DeleteAddress(userID int) error {
	return u.userRepo.DeleteAddress(userID)
}
func (u *UserUseCase) GetAllAddresses(id int) ([]domain.Address, error) {
	address, err := u.userRepo.GetAllAddresses(id)
	if err != nil {
		return []domain.Address{}, err
	}
	return address, nil
}
