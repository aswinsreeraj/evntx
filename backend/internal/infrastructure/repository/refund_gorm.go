package repository

import (
	"context"
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"gorm.io/gorm"
)

type RefundRequestModel struct {
	ID          string     `gorm:"type:uuid;primaryKey"`
	UserID      string     `gorm:"type:uuid;index;not null"`
	BookingID   string     `gorm:"type:uuid;index;not null"`
	Amount      float64    `gorm:"type:numeric(18,2);not null"`
	Status      string     `gorm:"size:50;not null"`
	RequestedAt time.Time  `gorm:"not null"`
	ProcessedAt *time.Time
	AdminID     *string    `gorm:"type:uuid"`
	CreatedAt   time.Time  `gorm:"not null"`
}

func (RefundRequestModel) TableName() string {
	return "refund_requests"
}

type refundGormRepository struct {
	db *gorm.DB
}

func NewRefundGormRepository(db *gorm.DB) *refundGormRepository {
	return &refundGormRepository{db: db}
}

func (r *refundGormRepository) CreateRefundRequest(ctx context.Context, req *domain.RefundRequest) error {
	model := RefundRequestModel{
		ID:          req.ID,
		UserID:      req.UserID,
		BookingID:   req.BookingID,
		Amount:      req.Amount,
		Status:      string(req.Status),
		RequestedAt: req.RequestedAt,
		ProcessedAt: req.ProcessedAt,
		AdminID:     req.AdminID,
		CreatedAt:   req.CreatedAt,
	}
	return r.db.WithContext(ctx).Create(&model).Error
}

func (r *refundGormRepository) GetRefundRequestByID(ctx context.Context, refundID string) (*domain.RefundRequest, error) {
	var model RefundRequestModel
	if err := r.db.WithContext(ctx).Where("id = ?", refundID).First(&model).Error; err != nil {
		return nil, err
	}
	return refundRequestModelToDomain(model), nil
}

func (r *refundGormRepository) UpdateRefundRequestStatus(ctx context.Context, refundID string, newStatus domain.RefundStatus, adminID *string) error {
	updates := map[string]interface{}{
		"status": string(newStatus),
	}
	if adminID != nil {
		updates["admin_id"] = *adminID
		t := time.Now()
		updates["processed_at"] = &t
	}
	return r.db.WithContext(ctx).Model(&RefundRequestModel{}).Where("id = ?", refundID).Updates(updates).Error
}

func (r *refundGormRepository) AdminGetRefundRequests(ctx context.Context, status string, page, limit int) ([]domain.AdminRefundDetail, int64, error) {
	type adminRefundRow struct {
		RefundRequestModel
		UserName  string `gorm:"column:user_name"`
		UserEmail string `gorm:"column:user_email"`
	}

	var total int64
	var rows []adminRefundRow

	query := r.db.WithContext(ctx).Table("refund_requests rr").
		Select("rr.*, u.name as user_name, u.email as user_email").
		Joins("LEFT JOIN users u ON rr.user_id = u.id")

	if status != "" {
		query = query.Where("rr.status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Order("rr.requested_at DESC").Limit(limit).Offset(offset).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	results := make([]domain.AdminRefundDetail, 0, len(rows))
	for _, row := range rows {
		results = append(results, domain.AdminRefundDetail{
			RefundRequest: *refundRequestModelToDomain(row.RefundRequestModel),
			UserName:      row.UserName,
			UserEmail:     row.UserEmail,
		})
	}

	return results, total, nil
}

func refundRequestModelToDomain(model RefundRequestModel) *domain.RefundRequest {
	return &domain.RefundRequest{
		ID:          model.ID,
		UserID:      model.UserID,
		BookingID:   model.BookingID,
		Amount:      model.Amount,
		Status:      domain.RefundStatus(model.Status),
		RequestedAt: model.RequestedAt,
		ProcessedAt: model.ProcessedAt,
		AdminID:     model.AdminID,
		CreatedAt:   model.CreatedAt,
	}
}
