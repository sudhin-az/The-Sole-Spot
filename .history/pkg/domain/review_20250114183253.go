package domain

type Review struct {
	ID        uint     `gorm:"primary_key"`
	User      Users    `gorm:"foreignKey:UserID"`
	UserID    uint     `gorm:"not null"`
	Product   Products `gorm:"foreignKey:ProductID"`
	ProductID uint     `gorm:"not null"`
	Rating    int      `gorm:"not null"`
	Comment   string   `gorm:"type:text"`
}
