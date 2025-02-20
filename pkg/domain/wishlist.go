package domain

type Wishlist struct {
	ID        uint     `gorm:"primary_key"`
	User      Users    `gorm:"foreignKey:UserID"`
	UserID    uint     `gorm:"not null"`
	Product   Products `gorm:"foreignKey:ProductID"`
	ProductID uint     `gorm:"not null"`
}
