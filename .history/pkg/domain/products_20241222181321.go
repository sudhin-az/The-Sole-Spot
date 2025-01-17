package domain

type Products struct {
	ID         int      `json:"id" gorm:"primarykey;not null"`
	CategoryID int      `json:"category_id"`
	Category   Category `json:"-" gorm:"foreignkey:CategoryID;constraint:OnDelete:CASCADE"`
	Name       string   `json:"name" validate:"required"`
	Quantity   int      `json:"quantity"`
	Stock      int      `json:"stock"`
	Price      float64  `json:"price"`
}
