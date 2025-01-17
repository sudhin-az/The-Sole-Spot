package domain

type Users struct {
	ID       int    `json:"id" gorm:"primary key,not null"`
	Name     string `json:"name"`
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"min=8,max=20"`
	Phone    string `json:"phone"`
	Blocked  bool   `json:"blocked" gorm:"default:false"`
}

type UserLoginMethod struct {
	ID                   uint   `gorm:"unique"`
	UserLoginMethodEmail string `json:"user_login_method_email" validate:"user_login_method_email"`
	LoginMethod          string
}

type GoogleResponse struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Picture string `json:"picture"`
}
