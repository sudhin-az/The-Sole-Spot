package utils

import (
	"errors"
	"regexp"
)

// ValidateFirstName validates that the first name contains only letters and spaces
func ValidateFirstName(firstName string) error {
	var validNamePattern = `^[a-zA-Z\s]+$`
	matched, _ := regexp.MatchString(validNamePattern, firstName)
	if !matched || len(firstName) < 2 {
		return errors.New("first name must contain only letters and spaces, and be at least 2 characters long")
	}
	return nil
}

// ValidateLastName validates that the last name contains only letters and spaces
func ValidateLastName(lastName string) error {
	var validNamePattern = `^[a-zA-Z\s]+$`
	matched, _ := regexp.MatchString(validNamePattern, lastName)
	if !matched || len(lastName) < 2 {
		return errors.New("last name must contain only letters and spaces, and be at least 2 characters long")
	}
	return nil
}

// ValidatePhoneNumber ensures the phone number is numeric and exactly 10 digits
func ValidatePhoneNumber(phone string) error {
	var validPhonePattern = `^[0-9]{10}$`
	matched, _ := regexp.MatchString(validPhonePattern, phone)
	if !matched {
		return errors.New("phone number must be exactly 10 numeric digits")
	}
	return nil
}

// ValidateEmail checks if the email is in the correct format
func ValidateEmail(email string) error {
	var validEmailPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(validEmailPattern, email)
	if !matched {
		return errors.New("invalid email format")
	}
	return nil
}
