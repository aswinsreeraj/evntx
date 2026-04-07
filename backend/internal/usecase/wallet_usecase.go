package usecase

import (
	"context"
	"math"
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/internal/repository"
	"github.com/aswinsreeraj/evntx/pkg/encryption"
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
	bookingRepo        repository.BookingRepository
	payoutRepo         repository.PayoutRepository
	refundRepo         repository.RefundRepository
}

func NewWalletUsecase(
	repo repository.WalletRepository,
	roleRepo repository.UserRoleRepository,
	platformWalletRepo repository.PlatformWalletRepository,
	razorpayService repository.RazorpayService,
	bookingRepo repository.BookingRepository,
	payoutRepo repository.PayoutRepository,
	refundRepo repository.RefundRepository,
) *WalletUsecase {
	return &WalletUsecase{repo: repo, roleRepo: roleRepo, platformWalletRepo: platformWalletRepo, razorpayService: razorpayService, bookingRepo: bookingRepo, payoutRepo: payoutRepo, refundRepo: refundRepo}
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

	txns, total, err := u.repo.GetTransactionsByWalletID(wallet.ID, filters, page, limit)
	if err != nil {
		return nil, 0, err
	}

	var bookingIDs []string
	for _, txn := range txns {
		if txn.ReferenceType == domain.WalletReferenceTypePurchase ||
			txn.ReferenceType == domain.WalletReferenceTypeUserCancellation ||
			txn.ReferenceType == domain.WalletReferenceTypeOrganizerCancellation ||
			txn.ReferenceType == domain.WalletReferenceTypeEarning ||
			txn.ReferenceType == domain.WalletReferenceTypeRefund {
			bookingIDs = append(bookingIDs, txn.ReferenceID)
		}
	}

	var contexts map[string]domain.BookingContextDetails
	if len(bookingIDs) > 0 && u.bookingRepo != nil {
		contexts, _ = u.bookingRepo.GetBookingContextsByIDs(context.Background(), bookingIDs)
	}

	for i, txn := range txns {
		txnCtx := &domain.WalletTransactionContext{
			Type: txn.ReferenceType,
		}

		if (txn.ReferenceType == domain.WalletReferenceTypePurchase ||
			txn.ReferenceType == domain.WalletReferenceTypeUserCancellation ||
			txn.ReferenceType == domain.WalletReferenceTypeOrganizerCancellation ||
			txn.ReferenceType == domain.WalletReferenceTypeEarning ||
			txn.ReferenceType == domain.WalletReferenceTypeRefund) &&
			contexts != nil {

			if bDetail, ok := contexts[txn.ReferenceID]; ok {
				txnCtx.Details = bDetail
			}
		} else if txn.ReferenceType == domain.WalletReferenceTypePayout {
			txnCtx.Details = domain.PayoutContextDetails{
				Amount:      txn.Amount,
				Status:      txn.Status,
				ProcessedAt: txn.CreatedAt,
			}
		} else if txn.ReferenceType == domain.WalletReferenceTypeFundAddition {
			txnCtx.Details = map[string]interface{}{
				"method": "razorpay",
				"id":     txn.ReferenceID,
			}
		}

		txns[i].Context = txnCtx
	}

	return txns, total, nil
}

func (u *WalletUsecase) AddPayoutCredentials(ctx context.Context, userID, holderName, accNumber, ifsc, upi string) error {
	if holderName == "" || accNumber == "" || ifsc == "" {
		return apiErrors.New(400, apiErrors.InvalidRequestBody, "Account Name, Number and IFSC are required")
	}

	encAcc, err := encryption.EncryptAES(accNumber)
	if err != nil {
		return apiErrors.ErrInternalServerError
	}
	encIFSC, err := encryption.EncryptAES(ifsc)
	if err != nil {
		return apiErrors.ErrInternalServerError
	}
	var encUPI *string
	if upi != "" {
		e, err := encryption.EncryptAES(upi)
		if err != nil {
			return apiErrors.ErrInternalServerError
		}
		encUPI = &e
	}

	cred := &domain.PayoutCredential{
		ID:                     uuid.NewString(),
		UserID:                 userID,
		AccountHolderName:      holderName,
		AccountNumberEncrypted: encAcc,
		IFSCCodeEncrypted:      encIFSC,
		UPIIDEncrypted:         encUPI,
		CreatedAt:              time.Now(),
		UpdatedAt:              time.Now(),
	}

	return u.payoutRepo.UpsertCredential(ctx, cred)
}

func (u *WalletUsecase) RequestPayout(ctx context.Context, userID string, amount float64) error {
	if amount <= 0 {
		return apiErrors.New(400, apiErrors.InvalidRequestBody, "Amount must be greater than zero")
	}

	_, err := u.payoutRepo.GetCredentialByUserID(ctx, userID)
	if err != nil {
		return apiErrors.New(400, apiErrors.InvalidRequestBody, "Payout credentials not found. Please set them up first.")
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

	return u.repo.WithTransaction(func(txRepo repository.WalletRepository) error {
		wallet.AvailableBalance -= normalizedAmount
		wallet.PendingBalance += normalizedAmount
		wallet.UpdatedAt = time.Now()

		if err := txRepo.UpdateWallet(wallet); err != nil {
			return err
		}

		payoutReq := &domain.PayoutRequest{
			ID:          uuid.NewString(),
			UserID:      userID,
			Amount:      normalizedAmount,
			Status:      domain.PayoutStatusPending,
			RequestedAt: time.Now(),
			CreatedAt:   time.Now(),
		}

		if err := u.payoutRepo.CreatePayoutRequest(ctx, payoutReq); err != nil {
			return err
		}

		txn := &domain.WalletTransaction{
			ID:            uuid.NewString(),
			WalletID:      wallet.ID,
			Type:          domain.WalletTransactionTypeDebit,
			Amount:        normalizedAmount,
			ReferenceType: domain.WalletReferenceTypePayout,
			ReferenceID:   payoutReq.ID,
			Status:        "pending",
			CreatedAt:     time.Now(),
		}

		if err := txRepo.CreateTransaction(txn); err != nil {
			return err
		}

		return nil
	})
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

func (u *WalletUsecase) AdminApprovePayout(ctx context.Context, adminID, payoutID string) error {
	p, err := u.payoutRepo.GetPayoutRequestByID(ctx, payoutID)
	if err != nil {
		return err
	}
	if p.Status != domain.PayoutStatusPending {
		return apiErrors.New(400, apiErrors.InvalidStateTransition, "Only pending payouts can be approved")
	}

	return u.payoutRepo.UpdatePayoutRequestStatus(ctx, payoutID, domain.PayoutStatusApproved, &adminID, nil)
}

func (u *WalletUsecase) AdminRejectPayout(ctx context.Context, adminID, payoutID, reason string) error {
	p, err := u.payoutRepo.GetPayoutRequestByID(ctx, payoutID)
	if err != nil {
		return err
	}
	if p.Status != domain.PayoutStatusPending {
		return apiErrors.New(400, apiErrors.InvalidStateTransition, "Only pending payouts can be rejected")
	}

	wallet, err := u.repo.GetWalletByUserID(p.UserID)
	if err != nil {
		return err
	}

	return u.repo.WithTransaction(func(txRepo repository.WalletRepository) error {
		wallet.PendingBalance -= p.Amount
		wallet.AvailableBalance += p.Amount
		wallet.UpdatedAt = time.Now()

		if err := txRepo.UpdateWallet(wallet); err != nil {
			return err
		}

		txn := &domain.WalletTransaction{
			ID:            uuid.NewString(),
			WalletID:      wallet.ID,
			Type:          domain.WalletTransactionTypeCredit,
			Amount:        p.Amount,
			ReferenceType: "payout_refund",
			ReferenceID:   payoutID,
			Status:        "completed",
			CreatedAt:     time.Now(),
		}

		if err := txRepo.CreateTransaction(txn); err != nil {
			return err
		}

		return u.payoutRepo.UpdatePayoutRequestStatus(ctx, payoutID, domain.PayoutStatusRejected, &adminID, &reason)
	})
}

func (u *WalletUsecase) AdminBulkApprovePayouts(ctx context.Context, adminID string, payoutIDs []string) error {
	for _, id := range payoutIDs {
		_ = u.AdminApprovePayout(ctx, adminID, id)
	}
	return nil
}

func (u *WalletUsecase) GetPayoutRequestsByUser(ctx context.Context, userID string, page, limit int) ([]domain.PayoutRequest, int64, error) {
	return u.payoutRepo.GetPayoutRequestsByUserID(ctx, userID, page, limit)
}

func (u *WalletUsecase) AdminGetPayoutRequests(ctx context.Context, status string, page, limit int) ([]domain.AdminPayoutDetail, int64, error) {
	return u.payoutRepo.AdminGetPayoutRequests(ctx, status, page, limit)
}

func (u *WalletUsecase) AdminGetRefundRequests(ctx context.Context, status string, page, limit int) ([]domain.AdminRefundDetail, int64, error) {
	return u.refundRepo.AdminGetRefundRequests(ctx, status, page, limit)
}

func (u *WalletUsecase) AdminProcessRefundRequest(ctx context.Context, adminID, refundID string) error {
	req, err := u.refundRepo.GetRefundRequestByID(ctx, refundID)
	if err != nil {
		return err
	}
	if req.Status != domain.RefundStatusPending {
		return apiErrors.New(400, apiErrors.InvalidStateTransition, "Refund is already processed")
	}

	if err := u.platformWalletRepo.ApplyPlatformTransaction(
		domain.WalletTransactionTypeDebit,
		req.Amount,
		domain.PlatformRefTypeRefund,
		refundID,
	); err != nil {
		return err
	}

	return u.refundRepo.UpdateRefundRequestStatus(ctx, refundID, domain.RefundStatusProcessed, &adminID)
}

func (u *WalletUsecase) AutoProcessApprovedPayouts(ctx context.Context) error {
	payouts, _, err := u.AdminGetPayoutRequests(ctx, string(domain.PayoutStatusApproved), 1, 1000)
	if err != nil {
		return err
	}

	systemReason := "Automated processing success"
	for _, p := range payouts {
		if p.Status == domain.PayoutStatusApproved {
			logger.Log.Info().Str("payout_id", p.ID).Msg("Auto-processing approved payout via Cron")
			u.payoutRepo.UpdatePayoutRequestStatus(ctx, p.ID, domain.PayoutStatusCompleted, nil, &systemReason)
		}
	}
	return nil
}
