package usecase

import (
	"ecommerce_clean_arch/pkg/repository"
	"ecommerce_clean_arch/pkg/utils/models"
)

type WalletUseCase struct {
	repository repository.WalletRepository
}

func NewWalletUseCase(repo repository.WalletRepository) *WalletUseCase {
	return &WalletUseCase{repository: repo}
}

func (wal *WalletUseCase) GetUserWallet(userID int) (*models.UserWallet, error) {
	userWallet, err := wal.repository.GetWallet(userID)
	if err != nil {
		return nil, err
	}
	return userWallet, nil
}

func (wal *WalletUseCase) GetWalletTransaction(userID int) (*[]models.WalletTransaction, error) {
	transaction, err := wal.repository.GetWalletTransaction(userID)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}
