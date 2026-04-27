package repository

import "github.com/aswinsreeraj/evntx/internal/domain"

type PlatformWalletRepository interface {
	EnsureExists() error
	GetPlatformWallet() (*domain.PlatformWallet, error)
	GetPlatformWalletStats() (*domain.PlatformWalletStats, error)
	GetPlatformTransactions(page, limit int) ([]domain.PlatformWalletTransaction, int64, error)
	ApplyPlatformTransaction(
		txnType string,
		amount float64,
		referenceType string,
		referenceID string,
	) error
}
