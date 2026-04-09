package repository

import (
	"context"
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
)

type EngagementRepository interface {
	CreateSession(ctx context.Context, session *domain.VisitorSession) error
	GetSessionByID(ctx context.Context, sessionID string) (*domain.VisitorSession, error)
	UpdateSessionLastSeen(ctx context.Context, sessionID string, userID *string) error
	
	LogEvent(ctx context.Context, event *domain.EngagementEvent) error
	
	GetDailyAggregates(ctx context.Context, eventID string, startDate, endDate time.Time) ([]domain.EventEngagementDaily, error)
	GetEngagementReport(ctx context.Context, eventIDs []string, startDate, endDate time.Time) (*domain.EngagementReportStats, error)
	
	
	IncrementSuccessfulBookings(ctx context.Context, eventID string) error
}
