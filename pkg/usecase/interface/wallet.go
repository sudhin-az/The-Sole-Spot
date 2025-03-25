package interfaces

import "ecommerce_clean_arch/pkg/utils/models"

type WalletUseCase interface {
	GetUserWallet(userID string) (*models.UserWallet, error)
	GetWalletTransaction(userID string) (*[]models.WalletTransaction, error)
}
