package repository

import "github.com/aswinsreeraj/evntx/internal/domain"

type WalletRepository interface {
	CreateWallet(wallet *domain.Wallet) error
	GetWalletByUserID(userID string) (*domain.Wallet, error)
}
