package repository

import (
	"encoding/json"
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
	apiErrors "github.com/aswinsreeraj/evntx/pkg/errors"
	"gorm.io/gorm"
)

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

func (r *paymentGormRepository) UpdateStatus(paymentID string, status string) error {
	return r.db.Model(&PaymentModel{}).
		Where("id = ?", paymentID).
		Update("status", status).Error
}

func (r *paymentGormRepository) MarkPaymentSuccess(paymentID string, bookingID string) error {
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

		return nil
	})
}
