package helper

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func ValidateAddress(input any) (string, error) {
	validate := validator.New()
	err := validate.Struct(input)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "HouseName":
				return "House name must be between 3 and 100 characters", fmt.Errorf("invalid house name")
			case "Street":
				return "Street address must be between 3 and 100 characters", fmt.Errorf("invalid street address")
			case "City":
				return "City must be between 2 and 50 characters", fmt.Errorf("invalid city")
			case "District":
				return "District must be between 3 and 50 characters", fmt.Errorf("invalid district")
			case "State":
				return "State must be exactly 3 characters", fmt.Errorf("invalid state code")
			case "Email":
				return "Invalid email format", fmt.Errorf("invalid email")
			case "Password":
				return "Password must be at least 5 characters long", fmt.Errorf("invalid password")
			case "Pin":
				return "Pin code must be exactly 6 digits", fmt.Errorf("invalid pin code")
			default:
				return fmt.Sprintf("Invalid input: %s", err.Field()), fmt.Errorf("validation failed")
			}
		}
	}
	return "", nil
}
