package repository

import (
	"math"
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
	apiErrors "github.com/aswinsreeraj/evntx/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PlatformWalletModel struct {
	ID               string    `gorm:"type:uuid;primaryKey"`
	AvailableBalance float64   `gorm:"type:numeric(18,2);default:0;not null"`
	PendingBalance   float64   `gorm:"type:numeric(18,2);default:0;not null"`
	RefundReserve    float64   `gorm:"type:numeric(18,2);default:0;not null"`
	TotalCredited    float64   `gorm:"type:numeric(18,2);default:0;not null"`
	TotalDebited     float64   `gorm:"type:numeric(18,2);default:0;not null"`
	UpdatedAt        time.Time `gorm:"not null"`
}

type PlatformWalletTransactionModel struct {
	ID            string    `gorm:"type:uuid;primaryKey"`
	WalletID      string    `gorm:"type:uuid;index;not null"`
	Type          string    `gorm:"size:2;not null"`
	Amount        float64   `gorm:"type:numeric(18,2);not null"`
	ReferenceType string    `gorm:"not null"`
	ReferenceID   string    `gorm:"not null"`
	CreatedAt     time.Time `gorm:"not null"`
}

type platformWalletGormRepository struct {
	db *gorm.DB
}

func NewPlatformWalletGormRepository(db *gorm.DB) *platformWalletGormRepository {
	return &platformWalletGormRepository{db: db}
}

func (r *platformWalletGormRepository) EnsureExists() error {
	wallet := PlatformWalletModel{
		ID:               domain.PlatformWalletID,
		AvailableBalance: 0,
		PendingBalance:   0,
		RefundReserve:    0,
		TotalCredited:    0,
		TotalDebited:     0,
		UpdatedAt:        time.Now(),
	}
	return r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&wallet).Error
}

func (r *platformWalletGormRepository) GetPlatformWallet() (*domain.PlatformWallet, error) {
	var model PlatformWalletModel
	if err := r.db.Where("id = ?", domain.PlatformWalletID).First(&model).Error; err != nil {
		return nil, err
	}
	return &domain.PlatformWallet{
		ID:               model.ID,
		AvailableBalance: model.AvailableBalance,
		PendingBalance:   model.PendingBalance,
		RefundReserve:    model.RefundReserve,
		TotalCredited:    model.TotalCredited,
		TotalDebited:     model.TotalDebited,
		UpdatedAt:        model.UpdatedAt,
	}, nil
}

func (r *platformWalletGormRepository) ApplyPlatformTransaction(
	txnType string,
	amount float64,
	referenceType string,
	referenceID string,
) error {
	if txnType != domain.WalletTransactionTypeCredit && txnType != domain.WalletTransactionTypeDebit {
		return apiErrors.New(400, apiErrors.InvalidRequestBody, "invalid platform transaction type")
	}

	normalized := math.Round(amount*100) / 100
	now := time.Now()

	return r.db.Transaction(func(tx *gorm.DB) error {
		var wallet PlatformWalletModel
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", domain.PlatformWalletID).
			First(&wallet).Error; err != nil {
			return err
		}

		if err := tx.Create(&PlatformWalletTransactionModel{
			ID:            uuid.NewString(),
			WalletID:      domain.PlatformWalletID,
			Type:          txnType,
			Amount:        normalized,
			ReferenceType: referenceType,
			ReferenceID:   referenceID,
			CreatedAt:     now,
		}).Error; err != nil {
			return err
		}

		switch txnType {
		case domain.WalletTransactionTypeCredit:
			wallet.AvailableBalance = math.Round((wallet.AvailableBalance+normalized)*100) / 100
			wallet.TotalCredited = math.Round((wallet.TotalCredited+normalized)*100) / 100
		case domain.WalletTransactionTypeDebit:
			if wallet.AvailableBalance < normalized {
				return apiErrors.ErrInsufficientBalance
			}
			wallet.AvailableBalance = math.Round((wallet.AvailableBalance-normalized)*100) / 100
			wallet.TotalDebited = math.Round((wallet.TotalDebited+normalized)*100) / 100
		}
		wallet.UpdatedAt = now

		return tx.Model(&PlatformWalletModel{}).
			Where("id = ?", domain.PlatformWalletID).
			Select("available_balance", "pending_balance", "refund_reserve", "total_credited", "total_debited", "updated_at").
			Updates(PlatformWalletModel{
				AvailableBalance: wallet.AvailableBalance,
				PendingBalance:   wallet.PendingBalance,
				RefundReserve:    wallet.RefundReserve,
				TotalCredited:    wallet.TotalCredited,
				TotalDebited:     wallet.TotalDebited,
				UpdatedAt:        wallet.UpdatedAt,
			}).Error
	})
}
