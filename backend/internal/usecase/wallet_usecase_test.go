package usecase_test

import (
	"context"
	"testing"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/internal/repository"
	"github.com/aswinsreeraj/evntx/internal/repository/mocks"
	"github.com/aswinsreeraj/evntx/internal/usecase"
	apiErrors "github.com/aswinsreeraj/evntx/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWalletUsecase_ApplyTransaction(t *testing.T) {
	mockWalletRepo := mocks.NewWalletRepository(t)

	walletUsecase := usecase.NewWalletUsecase(
		mockWalletRepo,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	t.Run("Success_Credit", func(t *testing.T) {
		walletID := "wallet-123"
		amount := 100.0

		wallet := &domain.Wallet{
			ID:               walletID,
			AvailableBalance: 50.0,
			TotalCredited:    50.0,
		}

		mockWalletRepo.On("WithTransaction", mock.AnythingOfType("func(repository.WalletRepository) error")).
			Run(func(args mock.Arguments) {
				fn := args.Get(0).(func(repository.WalletRepository) error)
				fn(mockWalletRepo)
			}).
			Return(nil).Once()

		mockWalletRepo.On("GetWalletByID", walletID).Return(wallet, nil).Once()
		mockWalletRepo.On("CreateTransaction", mock.AnythingOfType("*domain.WalletTransaction")).Return(nil).Once()
		mockWalletRepo.On("UpdateWallet", mock.AnythingOfType("*domain.Wallet")).Return(nil).Once()

		err := walletUsecase.ApplyTransaction(walletID, domain.WalletTransactionTypeCredit, amount, domain.WalletReferenceTypeFundAddition, "ref-123")

		assert.NoError(t, err)
		assert.Equal(t, float64(150.0), wallet.AvailableBalance)
		assert.Equal(t, float64(150.0), wallet.TotalCredited)
	})

	t.Run("Failure_InsufficientBalance_Debit", func(t *testing.T) {
		walletID := "wallet-123"
		amount := 100.0

		wallet := &domain.Wallet{
			ID:               walletID,
			AvailableBalance: 50.0, 
		}

		mockWalletRepo.On("WithTransaction", mock.AnythingOfType("func(repository.WalletRepository) error")).
			Run(func(args mock.Arguments) {
				fn := args.Get(0).(func(repository.WalletRepository) error)
				fn(mockWalletRepo)
			}).
			Return(apiErrors.ErrInsufficientBalance).Once()

		mockWalletRepo.On("GetWalletByID", walletID).Return(wallet, nil).Once()

		err := walletUsecase.ApplyTransaction(walletID, domain.WalletTransactionTypeDebit, amount, domain.WalletReferenceTypePurchase, "ref-123")

		assert.ErrorIs(t, err, apiErrors.ErrInsufficientBalance)
	})
}

func TestWalletUsecase_RequestPayout(t *testing.T) {
	mockWalletRepo := mocks.NewWalletRepository(t)
	mockPayoutRepo := mocks.NewPayoutRepository(t)

	walletUsecase := usecase.NewWalletUsecase(
		mockWalletRepo,
		nil,
		nil,
		nil,
		nil,
		mockPayoutRepo,
		nil,
	)

	t.Run("Success_RequestPayout", func(t *testing.T) {
		ctx := context.Background()
		userID := "user-123"
		amount := 50.0

		wallet := &domain.Wallet{
			ID:               "wallet-123",
			UserID:           userID,
			AvailableBalance: 100.0,
			PendingBalance:   0.0,
			ReserveBalance:   0.0, 
		}

		mockPayoutRepo.On("GetCredentialByUserID", ctx, userID).Return(&domain.PayoutCredential{}, nil).Once()
		mockWalletRepo.On("GetWalletByUserID", userID).Return(wallet, nil).Once()

		mockWalletRepo.On("WithTransaction", mock.AnythingOfType("func(repository.WalletRepository) error")).
			Run(func(args mock.Arguments) {
				fn := args.Get(0).(func(repository.WalletRepository) error)
				fn(mockWalletRepo)
			}).
			Return(nil).Once()

		mockWalletRepo.On("UpdateWallet", mock.AnythingOfType("*domain.Wallet")).Return(nil).Once()
		mockPayoutRepo.On("CreatePayoutRequest", ctx, mock.AnythingOfType("*domain.PayoutRequest")).Return(nil).Once()
		mockWalletRepo.On("CreateTransaction", mock.AnythingOfType("*domain.WalletTransaction")).Return(nil).Once()

		err := walletUsecase.RequestPayout(ctx, userID, amount)

		assert.NoError(t, err)
		assert.Equal(t, float64(50.0), wallet.AvailableBalance)
		assert.Equal(t, float64(50.0), wallet.PendingBalance)
	})

	t.Run("Failure_Deficit", func(t *testing.T) {
		ctx := context.Background()
		userID := "user-123"
		amount := 50.0

		wallet := &domain.Wallet{
			ID:               "wallet-123",
			UserID:           userID,
			AvailableBalance: 100.0,
			ReserveBalance:   -60.0, 
		}

		mockPayoutRepo.On("GetCredentialByUserID", ctx, userID).Return(&domain.PayoutCredential{}, nil).Once()
		mockWalletRepo.On("GetWalletByUserID", userID).Return(wallet, nil).Once()

		err := walletUsecase.RequestPayout(ctx, userID, amount)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Insufficient withdrawable balance")
	})
}
