package usecase

import (
	"math"
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/internal/repository"
	apiErrors "github.com/aswinsreeraj/evntx/pkg/errors"
	"github.com/aswinsreeraj/evntx/pkg/logger"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WalletUsecase struct {
	repo               repository.WalletRepository
	roleRepo           repository.UserRoleRepository
	platformWalletRepo repository.PlatformWalletRepository
	razorpayService    repository.RazorpayService
}

func NewWalletUsecase(
	repo repository.WalletRepository,
	roleRepo repository.UserRoleRepository,
	platformWalletRepo repository.PlatformWalletRepository,
	razorpayService repository.RazorpayService,
) *WalletUsecase {
	return &WalletUsecase{repo: repo, roleRepo: roleRepo, platformWalletRepo: platformWalletRepo, razorpayService: razorpayService}
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
				wallet.TotalCredited = normalizeWalletAmount(wallet.TotalCredited + normalizedAmount)
			} else if referenceType == domain.WalletReferenceTypeSettlement {
				if wallet.PendingBalance < normalizedAmount {
					return apiErrors.ErrInsufficientBalance
				}
				wallet.PendingBalance = normalizeWalletAmount(wallet.PendingBalance - normalizedAmount)
				wallet.AvailableBalance = normalizeWalletAmount(wallet.AvailableBalance + normalizedAmount)
			} else {
				wallet.AvailableBalance = normalizeWalletAmount(wallet.AvailableBalance + normalizedAmount)
				wallet.TotalCredited = normalizeWalletAmount(wallet.TotalCredited + normalizedAmount)
			}
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

func (u *WalletUsecase) RequestPayout(userID string, amount float64, accountName, accountNumber, ifscCode string) error {
	if amount <= 0 {
		return apiErrors.New(400, apiErrors.InvalidRequestBody, "Amount must be greater than zero")
	}

	if accountName == "" || accountNumber == "" || ifscCode == "" {
		return apiErrors.New(400, apiErrors.InvalidRequestBody, "Bank details are required")
	}

	wallet, err := u.repo.GetWalletByUserID(userID)
	if err != nil {
		return err
	}

	normalizedAmount := normalizeWalletAmount(amount)

	lockedAmount := 0.0
	if wallet.ReserveBalance < 0 {
		lockedAmount = math.Abs(wallet.ReserveBalance)
	}

	withdrawableBalance := wallet.AvailableBalance - lockedAmount
	if withdrawableBalance < normalizedAmount {
		return apiErrors.New(400, apiErrors.InsufficientBalance, "Insufficient withdrawable balance. Please note any reserve deficits.")
	}

	if err := u.ApplyTransaction(
		wallet.ID,
		domain.WalletTransactionTypeDebit,
		normalizedAmount,
		"payout",
		userID,
	); err != nil {
		return err
	}

	if u.platformWalletRepo != nil {
		if notifyErr := u.platformWalletRepo.ApplyPlatformTransaction(
			domain.WalletTransactionTypeDebit,
			normalizedAmount,
			domain.PlatformRefTypePayout,
			userID,
		); notifyErr != nil {
			return notifyErr
		}
	}

	return nil
}

func normalizeWalletAmount(amount float64) float64 {
	return math.Round(amount*100) / 100
}

func (u *WalletUsecase) CreateAddFundOrder(userID string, amount float64) (*domain.PaymentOrderResponse, error) {
	if amount <= 0 {
		return nil, apiErrors.New(400, apiErrors.InvalidRequestBody, "Amount must be greater than zero")
	}

	wallet, err := u.repo.GetWalletByUserID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {

			wallet = &domain.Wallet{
				ID:               uuid.NewString(),
				UserID:           userID,
				AvailableBalance: 0,
				PendingBalance:   0,
				UpdatedAt:        time.Now(),
			}
			if createErr := u.repo.CreateWallet(wallet); createErr != nil {
				logger.Log.Error().Err(createErr).Str("user_id", userID).Msg("failed to create missing wallet during add-fund")
				return nil, apiErrors.Wrap(createErr, 500, apiErrors.InternalServerError, "Failed to initialize wallet")
			}
		} else {
			return nil, err
		}
	}

	amountPaise := int64(amount * 100)
	receipt := "fund_" + wallet.ID[:8]

	order, err := u.razorpayService.CreateOrder(amountPaise, receipt)
	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("user_id", userID).
			Float64("amount", amount).
			Msg("failed to create razorpay order in usecase")
		return nil, apiErrors.Wrap(err, 500, apiErrors.PaymentFailed, "Failed to create payment order")
	}

	return &domain.PaymentOrderResponse{
		OrderID:     order.ID,
		Amount:      order.Amount,
		Currency:    order.Currency,
		RazorpayKey: u.razorpayService.GetKeyID(),
	}, nil
}

func (u *WalletUsecase) VerifyAddFundPayment(userID string, razorpayOrderID, razorpayPaymentID, razorpaySignature string) error {
	wallet, err := u.repo.GetWalletByUserID(userID)
	if err != nil {
		return err
	}

	isValid, err := u.razorpayService.VerifySignature(razorpayOrderID, razorpayPaymentID, razorpaySignature)
	if err != nil {
		return apiErrors.Wrap(err, 500, apiErrors.PaymentFailed, "Failed to verify payment signature")
	}

	if !isValid {
		return apiErrors.New(400, apiErrors.PaymentFailed, "Invalid payment signature")
	}

	orderData, err := u.razorpayService.FetchOrder(razorpayOrderID)
	if err != nil {
		return apiErrors.Wrap(err, 500, apiErrors.PaymentFailed, "Failed to fetch order details")
	}

	amount := float64(orderData.Amount) / 100

	err = u.ApplyTransaction(
		wallet.ID,
		domain.WalletTransactionTypeCredit,
		amount,
		domain.WalletReferenceTypeFundAddition,
		razorpayOrderID,
	)
	if err != nil {
		return err
	}

	if u.platformWalletRepo != nil {
		if notifyErr := u.platformWalletRepo.ApplyPlatformTransaction(
			domain.WalletTransactionTypeCredit,
			amount,
			domain.PlatformRefTypeFundAddition,
			userID,
		); notifyErr != nil {
			logger.Log.Error().Err(notifyErr).Msg("Failed to apply platform transaction for fund addition")
		}
	}

	return nil
}
