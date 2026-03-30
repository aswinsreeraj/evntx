package repository

import (
	"encoding/json"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
	apiErrors "github.com/aswinsreeraj/evntx/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func getPlatformFeeRate() float64 {
	val := os.Getenv("PLATFORM_FEE_PERCENTAGE")
	if val == "" {
		return 0.05
	}
	f, err := strconv.ParseFloat(val, 64)
	if err != nil || f < 0 || f > 100 {
		return 0.05
	}
	return f / 100
}

type PaymentModel struct {
	ID                string          `gorm:"type:uuid;primaryKey" json:"id"`
	BookingID         string          `gorm:"type:uuid;index;not null" json:"booking_id"`
	Provider          string          `gorm:"not null" json:"provider"`
	ProviderReference string          `gorm:"uniqueIndex;not null" json:"provider_reference"`
	Amount            float64         `gorm:"type:numeric(12,2);not null" json:"amount"`
	Status            string          `gorm:"not null" json:"status"`
	RawResponse       json.RawMessage `gorm:"type:jsonb" json:"raw_response"`
	CreatedAt         time.Time       `json:"created_at"`
}

type paymentGormRepository struct {
	db *gorm.DB
}

func NewPaymentGormRepository(db *gorm.DB) *paymentGormRepository {
	return &paymentGormRepository{db: db}
}

func (r *paymentGormRepository) CreatePayment(payment *domain.Payment) error {
	model := PaymentModel{
		ID:                payment.ID,
		BookingID:         payment.BookingID,
		Provider:          payment.Provider,
		ProviderReference: payment.ProviderReference,
		Amount:            payment.Amount,
		Status:            payment.Status,
		RawResponse:       payment.RawResponse,
		CreatedAt:         payment.CreatedAt,
	}

	return r.db.Create(&model).Error
}

func (r *paymentGormRepository) FindByProviderReference(orderID string) (*domain.Payment, error) {
	var model PaymentModel

	if err := r.db.Where("provider_reference = ?", orderID).First(&model).Error; err != nil {
		return nil, err
	}

	return &domain.Payment{
		ID:                model.ID,
		BookingID:         model.BookingID,
		Provider:          model.Provider,
		ProviderReference: model.ProviderReference,
		Amount:            model.Amount,
		Status:            model.Status,
		RawResponse:       model.RawResponse,
		CreatedAt:         model.CreatedAt,
	}, nil
}

func (r *paymentGormRepository) FindByBookingID(bookingID string) (*domain.Payment, error) {
	var model PaymentModel

	if err := r.db.Where("booking_id = ?", bookingID).Order("created_at DESC").First(&model).Error; err != nil {
		return nil, err
	}

	return &domain.Payment{
		ID:                model.ID,
		BookingID:         model.BookingID,
		Provider:          model.Provider,
		ProviderReference: model.ProviderReference,
		Amount:            model.Amount,
		Status:            model.Status,
		RawResponse:       model.RawResponse,
		CreatedAt:         model.CreatedAt,
	}, nil
}

func (r *paymentGormRepository) UpdateStatus(paymentID string, status string) error {
	return r.db.Model(&PaymentModel{}).
		Where("id = ?", paymentID).
		Update("status", status).Error
}

func (r *paymentGormRepository) MarkPaymentSuccess(paymentID string, bookingID string, organizerID string, amount float64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		paymentResult := tx.Model(&PaymentModel{}).
			Where("id = ?", paymentID).
			Update("status", domain.PaymentStatusSuccess)
		if paymentResult.Error != nil {
			return paymentResult.Error
		}
		if paymentResult.RowsAffected == 0 {
			return apiErrors.ErrResourceNotFound
		}

		bookingResult := tx.Model(&BookingModel{}).
			Where("id = ? AND status IN ?", bookingID, []string{"reserved", "paid"}).
			Update("status", "paid")
		if bookingResult.Error != nil {
			return bookingResult.Error
		}
		if bookingResult.RowsAffected == 0 {
			return apiErrors.ErrInvalidStateTransition
		}

		normalizedAmount := math.Round(amount*100) / 100
		now := time.Now()

		var booking BookingModel
		if err := tx.Where("id = ?", bookingID).First(&booking).Error; err != nil {
			return err
		}

		var totalTickets int64
		if err := tx.Model(&BookingTicketModel{}).Where("booking_id = ?", bookingID).Select("COALESCE(SUM(quantity), 0)").Scan(&totalTickets).Error; err != nil {
			return err
		}

		userPlatformFee := float64(totalTickets * 30)
		baseTicketRevenue := math.Round((normalizedAmount-userPlatformFee)*100) / 100

		// 2. Platform Wallet: Credit userPlatformFee
		var platformWallet PlatformWalletModel
		if err := tx.Where("id = ?", domain.PlatformWalletID).First(&platformWallet).Error; err != nil {
			return err
		}
		if err := tx.Create(&PlatformWalletTransactionModel{
			ID:            uuid.NewString(),
			WalletID:      domain.PlatformWalletID,
			Type:          domain.WalletTransactionTypeCredit,
			Amount:        userPlatformFee,
			ReferenceType: domain.PlatformRefTypePayment,
			ReferenceID:   bookingID,
			CreatedAt:     now,
		}).Error; err != nil {
			return err
		}
		platformWallet.AvailableBalance = math.Round((platformWallet.AvailableBalance+userPlatformFee)*100) / 100
		platformWallet.TotalCredited = math.Round((platformWallet.TotalCredited+userPlatformFee)*100) / 100
		if err := tx.Model(&PlatformWalletModel{}).Where("id = ?", domain.PlatformWalletID).Updates(map[string]interface{}{
			"available_balance": platformWallet.AvailableBalance,
			"total_credited":    platformWallet.TotalCredited,
			"updated_at":        now,
		}).Error; err != nil {
			return err
		}

		// 3. Organizer Wallet
		var wallet WalletModel
		if err := tx.Where("user_id = ?", organizerID).First(&wallet).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return apiErrors.ErrResourceNotFound
			}
			return err
		}
		if err := tx.Create(&WalletTransactionModel{
			ID:            uuid.NewString(),
			WalletID:      wallet.ID,
			Type:          domain.WalletTransactionTypeCredit,
			Amount:        baseTicketRevenue,
			ReferenceType: domain.WalletReferenceTypeEarning,
			ReferenceID:   bookingID,
			Status:        domain.WalletTransactionStatusCompleted,
			CreatedAt:     now,
		}).Error; err != nil {
			return err
		}
		wallet.PendingBalance = math.Round((wallet.PendingBalance+baseTicketRevenue)*100) / 100
		wallet.ReserveBalance = math.Round((wallet.ReserveBalance-userPlatformFee)*100) / 100
		wallet.TotalCredited = math.Round((wallet.TotalCredited+baseTicketRevenue)*100) / 100
		if err := tx.Model(&WalletModel{}).Where("id = ?", wallet.ID).Updates(map[string]interface{}{
			"pending_balance": wallet.PendingBalance,
			"reserve_balance": wallet.ReserveBalance,
			"total_credited":  wallet.TotalCredited,
			"updated_at":      now,
		}).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *paymentGormRepository) RefundPaymentToWallet(
	userID string,
	paymentID string,
	bookingID string,
	refundAmount float64,
	platformFeeAmount float64,
) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		paymentResult := tx.Model(&PaymentModel{}).
			Where("id = ? AND status = ?", paymentID, domain.PaymentStatusSuccess).
			Update("status", domain.PaymentStatusRefunded)
		if paymentResult.Error != nil {
			return paymentResult.Error
		}
		if paymentResult.RowsAffected == 0 {
			return apiErrors.ErrInvalidStateTransition
		}

		normalizedRefundAmount := math.Round(refundAmount*100) / 100
		totalDebited := math.Round((refundAmount+platformFeeAmount)*100) / 100
		now := time.Now()

		var platformWallet PlatformWalletModel
		if err := tx.Where("id = ?", domain.PlatformWalletID).First(&platformWallet).Error; err != nil {
			return err
		}
		if platformWallet.AvailableBalance < totalDebited {
			return apiErrors.ErrInsufficientBalance
		}
		if err := tx.Create(&PlatformWalletTransactionModel{
			ID:            uuid.NewString(),
			WalletID:      domain.PlatformWalletID,
			Type:          domain.WalletTransactionTypeDebit,
			Amount:        totalDebited,
			ReferenceType: domain.PlatformRefTypeRefund,
			ReferenceID:   bookingID,
			CreatedAt:     now,
		}).Error; err != nil {
			return err
		}
		platformWallet.AvailableBalance = math.Round((platformWallet.AvailableBalance-totalDebited)*100) / 100
		platformWallet.TotalDebited = math.Round((platformWallet.TotalDebited+totalDebited)*100) / 100
		platformWallet.UpdatedAt = now
		if err := tx.Model(&PlatformWalletModel{}).
			Where("id = ?", domain.PlatformWalletID).
			Select("available_balance", "total_debited", "updated_at").
			Updates(PlatformWalletModel{
				AvailableBalance: platformWallet.AvailableBalance,
				TotalDebited:     platformWallet.TotalDebited,
				UpdatedAt:        platformWallet.UpdatedAt,
			}).Error; err != nil {
			return err
		}

		if normalizedRefundAmount <= 0 {
			return nil
		}

		var wallet WalletModel
		if err := tx.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return apiErrors.ErrResourceNotFound
			}
			return err
		}
		if err := tx.Create(&WalletTransactionModel{
			ID:            uuid.NewString(),
			WalletID:      wallet.ID,
			Type:          domain.WalletTransactionTypeCredit,
			Amount:        normalizedRefundAmount,
			ReferenceType: domain.WalletReferenceTypeRefund,
			ReferenceID:   bookingID,
			Status:        domain.WalletTransactionStatusCompleted,
			CreatedAt:     now,
		}).Error; err != nil {
			return err
		}
		wallet.AvailableBalance = math.Round((wallet.AvailableBalance+normalizedRefundAmount)*100) / 100
		wallet.TotalCredited = math.Round((wallet.TotalCredited+normalizedRefundAmount)*100) / 100
		wallet.UpdatedAt = now
		return tx.Model(&WalletModel{}).
			Where("id = ?", wallet.ID).
			Select("available_balance", "total_credited", "updated_at").
			Updates(WalletModel{
				AvailableBalance: wallet.AvailableBalance,
				TotalCredited:    wallet.TotalCredited,
				UpdatedAt:        wallet.UpdatedAt,
			}).Error
	})
}
