package repository

import "github.com/aswinsreeraj/evntx/internal/domain"

type WalletRepository interface {
	CreateWallet(wallet *domain.Wallet) error
	GetWalletByUserID(userID string) (*domain.Wallet, error)
	GetWalletByID(walletID string) (*domain.Wallet, error)
	UpdateWallet(wallet *domain.Wallet) error
	CreateTransaction(txn *domain.WalletTransaction) error
	GetTransactionsByWalletID(
		walletID string,
		filters domain.WalletTransactionFilter,
		page int,
		limit int,
	) ([]domain.WalletTransaction, int64, error)
	WithTransaction(fn func(repo WalletRepository) error) error
}
