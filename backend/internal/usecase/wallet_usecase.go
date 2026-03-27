package usecase

import (
	"math"
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/internal/repository"
	apiErrors "github.com/aswinsreeraj/evntx/pkg/errors"
	"github.com/google/uuid"
)

type WalletUsecase struct {
	repo repository.WalletRepository
}

func NewWalletUsecase(repo repository.WalletRepository) *WalletUsecase {
	return &WalletUsecase{repo: repo}
}

func (u *WalletUsecase) GetWalletByUserID(userID string) (*domain.Wallet, error) {
	return u.repo.GetWalletByUserID(userID)
}

func (u *WalletUsecase) ApplyTransaction(
	walletID string,
	txnType string,
	amount float64,
	referenceType string,
	referenceID string,
) error {
	if txnType != domain.WalletTransactionTypeCredit && txnType != domain.WalletTransactionTypeDebit {
		return apiErrors.New(400, apiErrors.InvalidRequestBody, "Invalid wallet transaction type")
	}

	if amount <= 0 {
		return apiErrors.New(400, apiErrors.InvalidRequestBody, "Amount must be greater than zero")
	}

	if referenceType == "" {
		return apiErrors.New(400, apiErrors.InvalidRequestBody, "Reference type is required")
	}

	if referenceID == "" {
		return apiErrors.New(400, apiErrors.InvalidRequestBody, "Reference ID is required")
	}

	normalizedAmount := normalizeWalletAmount(amount)
	now := time.Now()

	return u.repo.WithTransaction(func(repo repository.WalletRepository) error {
		wallet, err := repo.GetWalletByID(walletID)
		if err != nil {
			return err
		}

		if txnType == domain.WalletTransactionTypeDebit && wallet.AvailableBalance < normalizedAmount {
			return apiErrors.ErrInsufficientBalance
		}

		txn := &domain.WalletTransaction{
			ID:            uuid.NewString(),
			WalletID:      wallet.ID,
			Type:          txnType,
			Amount:        normalizedAmount,
			ReferenceType: referenceType,
			ReferenceID:   referenceID,
			Status:        domain.WalletTransactionStatusCompleted,
			CreatedAt:     now,
		}

		if err := repo.CreateTransaction(txn); err != nil {
			return err
		}

		switch txnType {
		case domain.WalletTransactionTypeCredit:
			if referenceType == domain.WalletReferenceTypeEarning {
				wallet.PendingBalance = normalizeWalletAmount(wallet.PendingBalance + normalizedAmount)
			} else {
				wallet.AvailableBalance = normalizeWalletAmount(wallet.AvailableBalance + normalizedAmount)
			}
			wallet.TotalCredited = normalizeWalletAmount(wallet.TotalCredited + normalizedAmount)
		case domain.WalletTransactionTypeDebit:
			wallet.AvailableBalance = normalizeWalletAmount(wallet.AvailableBalance - normalizedAmount)
			wallet.TotalDebited = normalizeWalletAmount(wallet.TotalDebited + normalizedAmount)
		}

		if wallet.AvailableBalance < 0 {
			return apiErrors.ErrInsufficientBalance
		}

		wallet.UpdatedAt = now

		return repo.UpdateWallet(wallet)
	})
}

func (u *WalletUsecase) GetTransactionsByUserID(
	userID string,
	filters domain.WalletTransactionFilter,
	page int,
	limit int,
) ([]domain.WalletTransaction, int64, error) {
	wallet, err := u.repo.GetWalletByUserID(userID)
	if err != nil {
		return nil, 0, err
	}

	return u.repo.GetTransactionsByWalletID(wallet.ID, filters, page, limit)
}

func normalizeWalletAmount(amount float64) float64 {
	return math.Round(amount*100) / 100
}
