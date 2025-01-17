package models

type AdminSignUp struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}
type AdminDetailsResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name" `
	Email string `json:"email" `
}
type AdminLogin struct {
	Email    string `json:"email"  binding:"required" validate:"required"`
	Password string `json:"password"  binding:"required" validate:"min=8,max=20"`
}
