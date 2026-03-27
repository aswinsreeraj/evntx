package repository

import (
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"gorm.io/gorm"
)

type WalletModel struct {
	ID               string  `gorm:"type:uuid;primaryKey"`
	UserID           string  `gorm:"type:uuid;uniqueIndex;not null"`
	AvailableBalance float64 `gorm:"type:numeric(18,2);default:0;not null"`
	PendingBalance   float64 `gorm:"type:numeric(18,2);default:0;not null"`
	TotalCredited    float64 `gorm:"type:numeric(18,2);default:0;not null"`
	TotalDebited     float64 `gorm:"type:numeric(18,2);default:0;not null"`
	UpdatedAt        time.Time
}

type walletGormRepository struct {
	db *gorm.DB
}

func NewWalletGormRepository(db *gorm.DB) *walletGormRepository {
	return &walletGormRepository{db: db}
}

func (r *walletGormRepository) CreateWallet(wallet *domain.Wallet) error {
	model := walletDomainToModel(wallet)
	return r.db.Create(&model).Error
}

func (r *walletGormRepository) GetWalletByUserID(userID string) (*domain.Wallet, error) {
	var model WalletModel

	if err := r.db.Where("user_id = ?", userID).First(&model).Error; err != nil {
		return nil, err
	}

	return walletModelToDomain(model), nil
}

func walletDomainToModel(wallet *domain.Wallet) WalletModel {
	return WalletModel{
		ID:               wallet.ID,
		UserID:           wallet.UserID,
		AvailableBalance: wallet.AvailableBalance,
		PendingBalance:   wallet.PendingBalance,
		TotalCredited:    wallet.TotalCredited,
		TotalDebited:     wallet.TotalDebited,
		UpdatedAt:        wallet.UpdatedAt,
	}
}

func walletModelToDomain(model WalletModel) *domain.Wallet {
	return &domain.Wallet{
		ID:               model.ID,
		UserID:           model.UserID,
		AvailableBalance: model.AvailableBalance,
		PendingBalance:   model.PendingBalance,
		TotalCredited:    model.TotalCredited,
		TotalDebited:     model.TotalDebited,
		UpdatedAt:        model.UpdatedAt,
	}
}
