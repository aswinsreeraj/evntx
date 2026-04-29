package repository

import (
	"context"

	"github.com/aswinsreeraj/evntx/internal/domain"
)

type PayoutRepository interface {
	UpsertCredential(ctx context.Context, cred *domain.PayoutCredential) error
	GetCredentialByUserID(ctx context.Context, userID string) (*domain.PayoutCredential, error)
	CreatePayoutRequest(ctx context.Context, req *domain.PayoutRequest) error
	GetPayoutRequestByID(ctx context.Context, payoutID string) (*domain.PayoutRequest, error)
	UpdatePayoutRequestStatus(ctx context.Context, payoutID string, newStatus domain.PayoutStatus, adminID *string, failureReason *string) error
	GetPayoutRequestsByUserID(ctx context.Context, userID string, page, limit int) ([]domain.PayoutRequest, int64, error)
	AdminGetPayoutRequests(ctx context.Context, status string, page, limit int) ([]domain.AdminPayoutDetail, int64, error)
	GetTotalPayoutsSum(ctx context.Context) (float64, error)
	WithTransaction(fn func(repo PayoutRepository) error) error
}
