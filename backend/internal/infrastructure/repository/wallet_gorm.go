package repository

import (
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
	repositoryContract "github.com/aswinsreeraj/evntx/internal/repository"
	"gorm.io/gorm"
)

type WalletModel struct {
	ID			string	`gorm:"type:uuid;primaryKey"`
	UserID			string	`gorm:"type:uuid;uniqueIndex;not null"`
	AvailableBalance	float64	`gorm:"type:numeric(18,2);default:0;not null"`
	PendingBalance		float64	`gorm:"type:numeric(18,2);default:0;not null"`
	ReserveBalance		float64	`gorm:"type:numeric(18,2);default:0;not null"`
	TotalCredited		float64	`gorm:"type:numeric(18,2);default:0;not null"`
	TotalDebited		float64	`gorm:"type:numeric(18,2);default:0;not null"`
	UpdatedAt		time.Time
}

type WalletTransactionModel struct {
	ID		string		`gorm:"type:uuid;primaryKey"`
	WalletID	string		`gorm:"type:uuid;index;not null"`
	Type		string		`gorm:"size:2;not null"`
	Amount		float64		`gorm:"type:numeric(18,2);not null"`
	ReferenceType	string		`gorm:"not null"`
	ReferenceID	string		`gorm:"not null"`
	Status		string		`gorm:"not null"`
	CreatedAt	time.Time	`gorm:"index;not null"`
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

func (r *walletGormRepository) GetWalletByID(walletID string) (*domain.Wallet, error) {
	var model WalletModel

	if err := r.db.Where("id = ?", walletID).First(&model).Error; err != nil {
		return nil, err
	}

	return walletModelToDomain(model), nil
}

func (r *walletGormRepository) UpdateWallet(wallet *domain.Wallet) error {
	return r.db.Model(&WalletModel{}).
		Where("id = ?", wallet.ID).
		Select(
			"available_balance",
			"pending_balance",
			"reserve_balance",
			"total_credited",
			"total_debited",
			"updated_at",
		).
		Updates(WalletModel{
			AvailableBalance:	wallet.AvailableBalance,
			PendingBalance:		wallet.PendingBalance,
			ReserveBalance:		wallet.ReserveBalance,
			TotalCredited:		wallet.TotalCredited,
			TotalDebited:		wallet.TotalDebited,
			UpdatedAt:		wallet.UpdatedAt,
		}).Error
}

func (r *walletGormRepository) CreateTransaction(txn *domain.WalletTransaction) error {
	model := WalletTransactionModel{
		ID:		txn.ID,
		WalletID:	txn.WalletID,
		Type:		txn.Type,
		Amount:		txn.Amount,
		ReferenceType:	txn.ReferenceType,
		ReferenceID:	txn.ReferenceID,
		Status:		txn.Status,
		CreatedAt:	txn.CreatedAt,
	}

	return r.db.Create(&model).Error
}

func (r *walletGormRepository) GetTransactionsByWalletID(
	walletID string,
	filters domain.WalletTransactionFilter,
	page int,
	limit int,
) ([]domain.WalletTransaction, int64, error) {
	var models []WalletTransactionModel
	var total int64

	query := r.db.Model(&WalletTransactionModel{}).Where("wallet_id = ?", walletID)

	if filters.Type != "" {
		query = query.Where("type = ?", filters.Type)
	}

	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&models).Error; err != nil {
		return nil, 0, err
	}

	transactions := make([]domain.WalletTransaction, 0, len(models))
	for _, model := range models {
		transactions = append(transactions, domain.WalletTransaction{
			ID:		model.ID,
			WalletID:	model.WalletID,
			Type:		model.Type,
			Amount:		model.Amount,
			ReferenceType:	model.ReferenceType,
			ReferenceID:	model.ReferenceID,
			Status:		model.Status,
			CreatedAt:	model.CreatedAt,
		})
	}

	return transactions, total, nil
}

func (r *walletGormRepository) WithTransaction(fn func(repo repositoryContract.WalletRepository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return fn(NewWalletGormRepository(tx))
	})
}

func walletDomainToModel(wallet *domain.Wallet) WalletModel {
	return WalletModel{
		ID:			wallet.ID,
		UserID:			wallet.UserID,
		AvailableBalance:	wallet.AvailableBalance,
		PendingBalance:		wallet.PendingBalance,
		ReserveBalance:		wallet.ReserveBalance,
		TotalCredited:		wallet.TotalCredited,
		TotalDebited:		wallet.TotalDebited,
		UpdatedAt:		wallet.UpdatedAt,
	}
}

func walletModelToDomain(model WalletModel) *domain.Wallet {
	return &domain.Wallet{
		ID:			model.ID,
		UserID:			model.UserID,
		AvailableBalance:	model.AvailableBalance,
		PendingBalance:		model.PendingBalance,
		ReserveBalance:		model.ReserveBalance,
		TotalCredited:		model.TotalCredited,
		TotalDebited:		model.TotalDebited,
		UpdatedAt:		model.UpdatedAt,
	}
}

func (r *walletGormRepository) UpdateTransactionStatusByReference(refType string, refID string, status string) error {
	return r.db.Model(&WalletTransactionModel{}).Where("reference_type = ? AND reference_id = ?", refType, refID).Update("status", status).Error
}
