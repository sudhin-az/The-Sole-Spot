package repository

import (
	"ecommerce_clean_architecture/pkg/domain"
	"ecommerce_clean_architecture/pkg/utils/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type WalletRepository struct {
	DB *gorm.DB
}

func NewWalletRepository(DB *gorm.DB) *WalletRepository {
	return &WalletRepository{DB: DB}
}

func (wal *WalletRepository) CreateOrUpdateWallet(tx *gorm.DB, userID int, creditAmount uint) (uint, error) {
	var wallet domain.Wallet
	err := tx.Where("user_id = ?", userID).First(&wallet).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			wallet = domain.Wallet{
				UserID:  uint(userID),
				Balance: creditAmount,
			}
			if err := tx.Create(&wallet).Error; err != nil {
				return 0, fmt.Errorf("failed to create wallet: %w", err)
			}
		} else {
			return 0, fmt.Errorf("failed to query wallet: %w", err)
		}
	} else {
		wallet.Balance += creditAmount
		if err := tx.Model(&wallet).Update("balance", wallet.Balance).Error; err != nil {
			return 0, fmt.Errorf("failed to update wallet balance: %w", err)
		}
	}
	return wallet.Balance, nil
}

func (wal *WalletRepository) GetWalletbalance(tx *gorm.DB, userID int) (uint, error) {
	var currentBalance uint
	err := tx.Raw("SELECT balance FROM wallets WHERE user_id = ?", userID).Scan(&currentBalance).Error
	if err != nil {
		return 0, fmt.Errorf("failed to fetch wallet balance: %w", err)
	}
	// if currentBalance == 0 {
	// 	return 0, errors.New("wallet balance is zero")
	// }
	return currentBalance, nil
}
func (wal *WalletRepository) WalletTransaction(tx *gorm.DB, transaction models.WalletTransaction) error {
	result := tx.Omit("TransactionID").Create(&transaction)
	return result.Error
}

func (wal *WalletRepository) GetWalletTransaction(userID int) (*[]models.WalletTransaction, error) {
	var transaction *[]models.WalletTransaction
	query := "SELECT * FROM wallet_transactions WHERE user_id = ?"
	result := wal.DB.Raw(query, userID).Scan(&transaction)
	if result.Error != nil {
		return nil, fmt.Errorf("face some issue while fetch user wallet transaction %w", result.Error)
	}
	// if result.RowsAffected == 0 {
	// 	return nil, errors.New("No rows affected while getting wallet transaction")
	// }
	return transaction, nil
}

func (wal *WalletRepository) GetFinalPriceByOrderID(orderID string) (uint, error) {
	var finalPrice uint
	query := "SELECT SUM(final_price) FROM orders WHERE order_id = ?"
	result := wal.DB.Raw(query, orderID).Scan(&finalPrice)
	if result.Error != nil {
		return 0, fmt.Errorf("face some issue while getting total amount of order by using order id %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return 0, errors.New("No rows affected while getting the price")
	}
	return finalPrice, nil
}

func (wal *WalletRepository) GetWallet(userID int) (*models.UserWallet, error) {
	var userWallet models.UserWallet
	query := "SELECT COALESCE(balance, 0),* FROM wallets WHERE user_id = ?"
	result := wal.DB.Raw(query, userID).Scan(&userWallet)
	if result.Error != nil {
		return nil, fmt.Errorf("face some issue while get user wallet %w", result.Error)
	}
	return &userWallet, nil
}

func (wal *WalletRepository) UpdateWalletReduceBalance(userID string, amount uint) error {
	query := "UPDATE wallets SET balance = ? WHERE user_id = ?"
	result := wal.DB.Raw(query, amount, userID)
	if result.Error != nil {
		return fmt.Errorf("face some issue while update wallet balance %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("No rows affected while updating the wallet balance")
	}
	return nil
}
