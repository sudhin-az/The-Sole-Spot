package interfaces

import "ecommerce_clean_architecture/pkg/utils/models"

type WalletRepository interface {
	CreateOrUpdateWallet(userID string, creditAmount uint) (uint, error)
	GetWalletbalance(userID string) (uint, error)
	WalletTransaction(transaction models.WalletTransaction) error
	GetWalletTransaction(userID string) (*[]models.WalletTransaction, error)
	GetFinalPriceByOrderID(orderID string) (uint, error)
	GetWallet(userID string) (*models.UserWallet, error)
	UpdateWalletReduceBalance(userID string, amount uint) error
}
