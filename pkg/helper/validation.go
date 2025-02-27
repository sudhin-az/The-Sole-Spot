package helper

import (
	"ecommerce_clean_architecture/pkg/utils/models"
	"fmt"

	"github.com/go-playground/validator/v10"
)

func ValidateAddress(input any) (string, error) {
	validate := validator.New()
	err := validate.Struct(input)
	if err != nil {
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
	// if password.Password != password.ConfirmPassword {
	// 	return errors.New("password do not match")
	// }
	// if password.Password == "" || password.ConfirmPassword == "" {
	// 	return errors.New("password cannot be empty")
	// }

	return nil
}

func ValidationErrorToText(err error) []string {
	var errors []string

	// Map field names to user-friendly labels
	fieldNames := map[string]string{
		"CouponCode":      "Coupon Code",
		"Discount":        "Discount",
		"MinimumRequired": "Minimum Required",
		"MaximumAllowed":  "Maximum Allowed",
		"MaximumUsage":    "Maximum Usage",
		"ExpireDate":      "Expire Date",
	}

	for _, err := range err.(validator.ValidationErrors) {
		var message string
		field := err.Field()

		// Get the user-friendly field name if it exists
		if name, exists := fieldNames[field]; exists {
			field = name
		}

		// Customize error messages for each validation tag
		switch err.Tag() {
		case "required":
			message = fmt.Sprintf("%s is required.", field)
		case "min":
			message = fmt.Sprintf("%s cannot be less than %s.", field, err.Param())
		case "max":
			message = fmt.Sprintf("%s cannot be greater than %s.", field, err.Param())
		case "oneof":
			message = fmt.Sprintf("%s must be one of %s.", field, err.Param())
		case "datetime":
			message = fmt.Sprintf("%s must be in the format YYYY-MM-DD.", field)
		case "gtcsfield":
			message = fmt.Sprintf("%s must be greater than %s.", field, err.Param())
		case "alphanum":
			message = fmt.Sprintf("%s must be alphanumeric.", field)
		case "len":
			message = fmt.Sprintf("%s must be exactly %s characters long.", field, err.Param())
		default:
			message = fmt.Sprintf("%s is invalid.", field)
		}

		errors = append(errors, message)
	}
	return errors
}
