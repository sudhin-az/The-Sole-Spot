package helper

import (
	"ecommerce_clean_architecture/pkg/utils/models"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

func ValidateAddress(input any) (string, error) {
	validate := validator.New()
	err := validate.Struct(input)
	if err != nil {
		// Return only the first validation error
		for _, fieldErr := range err.(validator.ValidationErrors) {
			switch fieldErr.Field() {
			case "HouseName":
				return "House name must be between 3 and 100 characters", fmt.Errorf("invalid house name")
			case "Street":
				return "Street address must be between 3 and 100 characters", fmt.Errorf("invalid street address")
			case "City":
				return "City must be between 2 and 50 characters", fmt.Errorf("invalid city")
			case "District":
				return "District must be between 3 and 50 characters", fmt.Errorf("invalid district")
			case "Pin":
				return "Pin code must be 5 or 6 characters", fmt.Errorf("invalid pin code")
			case "State":
				return "Invalid state code (use 2-letter ISO code)", fmt.Errorf("invalid state code")
			case "Email":
				return "Invalid email format", fmt.Errorf("invalid email")
			case "Password":
				return "Password must be at least 5 characters long", fmt.Errorf("invalid password")
			case "FirstName":
				return "First name must contain only letters and spaces, and be at least 2 characters long", fmt.Errorf("invalid first name")
			case "LastName":
				return "Last name must contain only letters and spaces, and be at least 2 characters long", fmt.Errorf("invalid last name")
			case "PhoneNumber":
				return "Invalid phone number format", fmt.Errorf("invalid phone number")
			}
			// Return a generic error if no specific field match
			return "Invalid input", fmt.Errorf("validation failed")
		}
	}
	return "", nil
}
func ValidatePassword(password models.ForgotPassword) error {
	var validationErrors []string
	if len(password.Password) < 5 {
		validationErrors = append(validationErrors, "Password must be at least 5 characters long")
	}
	if password.Password != password.ConfirmPassword {
		return errors.New("password do not match")
	}
	if password.Password == "" || password.ConfirmPassword == "" {
		return errors.New("password cannot be empty")
	}
	if password.Otp != repository.
	return nil
}
