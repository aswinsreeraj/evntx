package repository

import "github.com/aswinsreeraj/evntx/internal/domain"

type PlatformWalletRepository interface {
	EnsureExists() error
	GetPlatformWallet() (*domain.PlatformWallet, error)
	ApplyPlatformTransaction(
		txnType string,
		amount float64,
		referenceType string,
		referenceID string,
	) error
}
