package domain

import "time"

type Wallet struct {
	WalletID uint  `gorm:"primaryKey;autoIncrement"`
	UserID   uint  `gorm:"uniqueIndex"`
	User     Users `gorm:"foreignkey:UserID;association_foreignkey:ID"`
	Balance  uint
}

type WalletTransaction struct {
	TransactionID uint      `gorm:"primaryKey;autoIncrement" json:"transactionID"`
	UserID        uint      `json:"userID"`
	Credit        uint      `json:"credit,omitempty"`
	Debit         uint      `json:"debit,omitempty"`
	EventDate     time.Time `json:"eventDate"`
	TotalAmount   uint      `json:"totalAmount"`
}
