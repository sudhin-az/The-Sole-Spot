package domain

type Category struct {
	ID               int    `json:"id" gorm:"primarykey;not null"`
	Category         string `json:"category"`
	Description      string `json:"description"`
	CategoryDiscount int    `json:"category_discount"`
}
