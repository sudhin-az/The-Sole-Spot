package interfaces

import "ecommerce_clean_architecture/pkg/utils/models"

type WalletUseCase interface {
	GetUserWallet(userID string) (*models.UserWallet, error)
	GetWalletTransaction(userID string) (*[]models.WalletTransaction, error)
}
