package usecase

import (
	"context"
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/internal/repository"
	"github.com/google/uuid"
)

type EngagementUsecase struct {
	repo repository.EngagementRepository
}

func NewEngagementUsecase(repo repository.EngagementRepository) *EngagementUsecase {
	return &EngagementUsecase{repo: repo}
}

func (u *EngagementUsecase) InitializeSession(ctx context.Context, userID *string, ipAddress, userAgent string) (*domain.VisitorSession, error) {
	session := &domain.VisitorSession{
		ID:		uuid.NewString(),
		UserID:		userID,
		IPAddress:	ipAddress,
		UserAgent:	userAgent,
		CreatedAt:	time.Now(),
		LastSeenAt:	time.Now(),
	}

	err := u.repo.CreateSession(ctx, session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (u *EngagementUsecase) TrackEvent(
	ctx context.Context,
	sessionID string,
	userID *string,
	eventType domain.EngagementEventType,
	eventID *string,
	metadata, ipAddress, userAgent string,
) error {

	_ = u.repo.UpdateSessionLastSeen(ctx, sessionID, userID)

	if eventID != nil && *eventID != "" {
		if _, err := uuid.Parse(*eventID); err != nil {

			eventID = nil
		}
	}

	if metadata == "" {
		metadata = "{}"
	}

	evt := &domain.EngagementEvent{
		ID:		uuid.NewString(),
		UserID:		userID,
		SessionID:	sessionID,
		EventID:	eventID,
		EventType:	eventType,
		Metadata:	metadata,
		IPAddress:	ipAddress,
		UserAgent:	userAgent,
		CreatedAt:	time.Now(),
	}

	return u.repo.LogEvent(ctx, evt)
}

func (u *EngagementUsecase) GetDailyReport(ctx context.Context, eventID string, startDate, endDate time.Time) ([]domain.EventEngagementDaily, error) {
	return u.repo.GetDailyAggregates(ctx, eventID, startDate, endDate)
}

func (u *EngagementUsecase) GetEngagementReport(ctx context.Context, eventIDs []string, startDate, endDate time.Time) (*domain.EngagementReportStats, error) {
	return u.repo.GetEngagementReport(ctx, eventIDs, startDate, endDate)
}
