package repository

import (
	"context"
	"github.com/aswinsreeraj/evntx/internal/domain"
)

type RefundRepository interface {
	CreateRefundRequest(ctx context.Context, req *domain.RefundRequest) error
	GetRefundRequestByID(ctx context.Context, refundID string) (*domain.RefundRequest, error)
	AdminGetRefundRequests(ctx context.Context, status string, page int, limit int) ([]domain.AdminRefundDetail, int64, error)
	UpdateRefundRequestStatus(ctx context.Context, refundID string, newStatus domain.RefundStatus, adminID *string) error
}
