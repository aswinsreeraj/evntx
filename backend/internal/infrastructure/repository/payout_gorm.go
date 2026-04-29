package repository

import (
	"context"
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
	repositoryContract "github.com/aswinsreeraj/evntx/internal/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PayoutRequestModel struct {
	ID		string		`gorm:"type:uuid;primaryKey"`
	UserID		string		`gorm:"type:uuid;index;not null"`
	Amount		float64		`gorm:"type:numeric(18,2);not null"`
	Status		string		`gorm:"size:50;not null"`
	RequestedAt	time.Time	`gorm:"not null"`
	ReviewedAt	*time.Time
	ProcessedAt	*time.Time
	AdminID		*string		`gorm:"type:uuid"`
	FailureReason	*string		`gorm:"type:text"`
	CreatedAt	time.Time	`gorm:"not null"`
}

type PayoutCredentialModel struct {
	ID			string	`gorm:"type:uuid;primaryKey"`
	UserID			string	`gorm:"type:uuid;uniqueIndex;not null"`
	AccountHolderName	string	`gorm:"not null"`
	AccountNumberEncrypted	string	`gorm:"not null"`
	IFSCCodeEncrypted	string	`gorm:"not null"`
	UPIIDEncrypted		*string
	CreatedAt		time.Time	`gorm:"not null"`
	UpdatedAt		time.Time	`gorm:"not null"`
}

func (PayoutRequestModel) TableName() string {
	return "payout_requests"
}

func (PayoutCredentialModel) TableName() string {
	return "payout_credentials"
}

type payoutGormRepository struct {
	db *gorm.DB
}

func NewPayoutGormRepository(db *gorm.DB) *payoutGormRepository {
	return &payoutGormRepository{db: db}
}

func (r *payoutGormRepository) UpsertCredential(ctx context.Context, cred *domain.PayoutCredential) error {
	model := PayoutCredentialModel{
		ID:			cred.ID,
		UserID:			cred.UserID,
		AccountHolderName:	cred.AccountHolderName,
		AccountNumberEncrypted:	cred.AccountNumberEncrypted,
		IFSCCodeEncrypted:	cred.IFSCCodeEncrypted,
		UPIIDEncrypted:		cred.UPIIDEncrypted,
		CreatedAt:		cred.CreatedAt,
		UpdatedAt:		cred.UpdatedAt,
	}

	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:	[]clause.Column{{Name: "user_id"}},
		DoUpdates:	clause.AssignmentColumns([]string{"account_holder_name", "account_number_encrypted", "ifsc_code_encrypted", "upi_id_encrypted", "updated_at"}),
	}).Create(&model).Error
}

func (r *payoutGormRepository) GetCredentialByUserID(ctx context.Context, userID string) (*domain.PayoutCredential, error) {
	var model PayoutCredentialModel
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&model).Error; err != nil {
		return nil, err
	}
	return &domain.PayoutCredential{
		ID:			model.ID,
		UserID:			model.UserID,
		AccountHolderName:	model.AccountHolderName,
		AccountNumberEncrypted:	model.AccountNumberEncrypted,
		IFSCCodeEncrypted:	model.IFSCCodeEncrypted,
		UPIIDEncrypted:		model.UPIIDEncrypted,
		CreatedAt:		model.CreatedAt,
		UpdatedAt:		model.UpdatedAt,
	}, nil
}

func (r *payoutGormRepository) CreatePayoutRequest(ctx context.Context, req *domain.PayoutRequest) error {
	model := PayoutRequestModel{
		ID:		req.ID,
		UserID:		req.UserID,
		Amount:		req.Amount,
		Status:		string(req.Status),
		RequestedAt:	req.RequestedAt,
		ReviewedAt:	req.ReviewedAt,
		ProcessedAt:	req.ProcessedAt,
		AdminID:	req.AdminID,
		FailureReason:	req.FailureReason,
		CreatedAt:	req.CreatedAt,
	}
	return r.db.WithContext(ctx).Create(&model).Error
}

func (r *payoutGormRepository) GetPayoutRequestByID(ctx context.Context, payoutID string) (*domain.PayoutRequest, error) {
	var model PayoutRequestModel
	if err := r.db.WithContext(ctx).Where("id = ?", payoutID).First(&model).Error; err != nil {
		return nil, err
	}
	return payoutRequestModelToDomain(model), nil
}

func (r *payoutGormRepository) UpdatePayoutRequestStatus(ctx context.Context, payoutID string, newStatus domain.PayoutStatus, adminID *string, failureReason *string) error {
	updates := map[string]interface{}{
		"status": string(newStatus),
	}
	if adminID != nil {
		updates["admin_id"] = *adminID
		t := time.Now()
		updates["reviewed_at"] = &t
	}
	if failureReason != nil {
		updates["failure_reason"] = *failureReason
	}
	return r.db.WithContext(ctx).Model(&PayoutRequestModel{}).Where("id = ?", payoutID).Updates(updates).Error
}

func (r *payoutGormRepository) GetPayoutRequestsByUserID(ctx context.Context, userID string, page, limit int) ([]domain.PayoutRequest, int64, error) {
	var models []PayoutRequestModel
	var total int64

	query := r.db.WithContext(ctx).Model(&PayoutRequestModel{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Order("requested_at DESC").Limit(limit).Offset(offset).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	var results []domain.PayoutRequest
	for _, m := range models {
		results = append(results, *payoutRequestModelToDomain(m))
	}
	return results, total, nil
}

func (r *payoutGormRepository) AdminGetPayoutRequests(ctx context.Context, status string, page, limit int) ([]domain.AdminPayoutDetail, int64, error) {
	type adminPayoutRow struct {
		PayoutRequestModel
		UserName	string	`gorm:"column:user_name"`
		UserEmail	string	`gorm:"column:user_email"`
	}

	var total int64
	var rows []adminPayoutRow

	query := r.db.WithContext(ctx).Table("payout_requests pr").
		Select("pr.*, u.name as user_name, u.email as user_email").
		Joins("LEFT JOIN user_models u ON pr.user_id::uuid = u.id::uuid")

	if status != "" {
		query = query.Where("pr.status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Order("pr.requested_at DESC").Limit(limit).Offset(offset).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	results := make([]domain.AdminPayoutDetail, 0, len(rows))
	for _, row := range rows {
		results = append(results, domain.AdminPayoutDetail{
			PayoutRequest:	*payoutRequestModelToDomain(row.PayoutRequestModel),
			UserName:	row.UserName,
			UserEmail:	row.UserEmail,
		})
	}

	return results, total, nil
}

func (r *payoutGormRepository) GetTotalPayoutsSum(ctx context.Context) (float64, error) {
	var total float64
	err := r.db.WithContext(ctx).Model(&PayoutRequestModel{}).Select("COALESCE(SUM(amount), 0)").Scan(&total).Error
	return total, err
}

func (r *payoutGormRepository) WithTransaction(fn func(repo repositoryContract.PayoutRepository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return fn(NewPayoutGormRepository(tx))
	})
}

func payoutRequestModelToDomain(model PayoutRequestModel) *domain.PayoutRequest {
	return &domain.PayoutRequest{
		ID:		model.ID,
		UserID:		model.UserID,
		Amount:		model.Amount,
		Status:		domain.PayoutStatus(model.Status),
		RequestedAt:	model.RequestedAt,
		ReviewedAt:	model.ReviewedAt,
		ProcessedAt:	model.ProcessedAt,
		AdminID:	model.AdminID,
		FailureReason:	model.FailureReason,
		CreatedAt:	model.CreatedAt,
	}
}
