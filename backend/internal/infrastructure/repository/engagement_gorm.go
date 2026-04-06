package repository

import (
	"context"
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type VisitorSessionModel struct {
	ID         string  `gorm:"type:uuid;primaryKey"`
	UserID     *string `gorm:"type:uuid"`
	IPAddress  string  `gorm:"type:varchar"`
	UserAgent  string  `gorm:"type:text"`
	CreatedAt  time.Time
	LastSeenAt time.Time
}

func (VisitorSessionModel) TableName() string {
	return "visitor_sessions"
}

type EngagementEventModel struct {
	ID         string  `gorm:"type:uuid;primaryKey"`
	UserID     *string `gorm:"type:uuid"`
	SessionID  string  `gorm:"type:uuid;index"`
	EventID    *string `gorm:"type:uuid;index"`
	EventType  string  `gorm:"index"`
	Metadata   string  `gorm:"type:json"`
	IPAddress  string  `gorm:"type:varchar"`
	UserAgent  string  `gorm:"type:text"`
	CreatedAt  time.Time `gorm:"index"`
}

func (EngagementEventModel) TableName() string {
	return "engagement_events"
}

type EventEngagementDailyModel struct {
	ID                 string    `gorm:"type:uuid;primaryKey"`
	EventID            string    `gorm:"type:uuid;uniqueIndex:idx_event_date"`
	Date               time.Time `gorm:"type:date;uniqueIndex:idx_event_date"`
	Visitors           int
	EventViews         int
	TicketsSelected    int
	CheckoutStarted    int
	SuccessfulBookings int
	CreatedAt          time.Time
}

func (EventEngagementDailyModel) TableName() string {
	return "event_engagement_daily"
}

type engagementGormRepository struct {
	db *gorm.DB
}

func NewEngagementGormRepository(db *gorm.DB) *engagementGormRepository {
	return &engagementGormRepository{db: db}
}

func (r *engagementGormRepository) CreateSession(ctx context.Context, session *domain.VisitorSession) error {
	model := VisitorSessionModel{
		ID:         session.ID,
		UserID:     session.UserID,
		IPAddress:  session.IPAddress,
		UserAgent:  session.UserAgent,
		CreatedAt:  session.CreatedAt,
		LastSeenAt: session.LastSeenAt,
	}
	return r.db.WithContext(ctx).Create(&model).Error
}

func (r *engagementGormRepository) GetSessionByID(ctx context.Context, sessionID string) (*domain.VisitorSession, error) {
	var model VisitorSessionModel
	if err := r.db.WithContext(ctx).Where("id = ?", sessionID).First(&model).Error; err != nil {
		return nil, err
	}
	return &domain.VisitorSession{
		ID:         model.ID,
		UserID:     model.UserID,
		IPAddress:  model.IPAddress,
		UserAgent:  model.UserAgent,
		CreatedAt:  model.CreatedAt,
		LastSeenAt: model.LastSeenAt,
	}, nil
}

func (r *engagementGormRepository) UpdateSessionLastSeen(ctx context.Context, sessionID string, userID *string) error {
	updates := map[string]interface{}{
		"last_seen_at": time.Now(),
	}
	if userID != nil {
		updates["user_id"] = *userID
	}
	return r.db.WithContext(ctx).Model(&VisitorSessionModel{}).Where("id = ?", sessionID).Updates(updates).Error
}

func (r *engagementGormRepository) LogEvent(ctx context.Context, event *domain.EngagementEvent) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Log the individual event
		evtModel := EngagementEventModel{
			ID:        event.ID,
			UserID:    event.UserID,
			SessionID: event.SessionID,
			EventID:   event.EventID,
			EventType: string(event.EventType),
			Metadata:  event.Metadata,
			IPAddress: event.IPAddress,
			UserAgent: event.UserAgent,
			CreatedAt: event.CreatedAt,
		}
		if err := tx.Create(&evtModel).Error; err != nil {
			return err
		}

		// Update daily aggregates if it's related to an event
		if event.EventID != nil {
			dateStr := event.CreatedAt.Format("2006-01-02")
			dateParsed, _ := time.Parse("2006-01-02", dateStr)

			upsertModel := EventEngagementDailyModel{
				ID:        uuid.NewString(),
				EventID:   *event.EventID,
				Date:      dateParsed,
				CreatedAt: time.Now(),
			}

			// Based on the event type, prepare the UPSERT clause
			var updateCol string
			switch event.EventType {
			case domain.EngagementEventEventView:
				upsertModel.EventViews = 1
				updateCol = "event_views"
			case domain.EngagementEventTicketSelected:
				upsertModel.TicketsSelected = 1
				updateCol = "tickets_selected"
			case domain.EngagementEventCheckoutStarted:
				upsertModel.CheckoutStarted = 1
				updateCol = "checkout_started"
			case domain.EngagementEventPageView:
				upsertModel.Visitors = 1
				updateCol = "visitors"
			}

			// If it's an event we track in the daily table
			if updateCol != "" {
				err := tx.Clauses(clause.OnConflict{
					Columns:   []clause.Column{{Name: "event_id"}, {Name: "date"}},
					DoUpdates: clause.Assignments(map[string]interface{}{
						updateCol: gorm.Expr("event_engagement_daily." + updateCol + " + 1"),
					}),
				}).Create(&upsertModel).Error
			
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (r *engagementGormRepository) IncrementSuccessfulBookings(ctx context.Context, eventID string) error {
	dateParsed, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	
	upsertModel := EventEngagementDailyModel{
		ID:                 uuid.NewString(),
		EventID:            eventID,
		Date:               dateParsed,
		SuccessfulBookings: 1,
		CreatedAt:          time.Now(),
	}

	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "event_id"}, {Name: "date"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"successful_bookings": gorm.Expr("event_engagement_daily.successful_bookings + ?", 1),
		}),
	}).Create(&upsertModel).Error
}

func (r *engagementGormRepository) GetDailyAggregates(ctx context.Context, eventID string, startDate, endDate time.Time) ([]domain.EventEngagementDaily, error) {
	var models []EventEngagementDailyModel
	
	query := r.db.WithContext(ctx).Where("event_id = ?", eventID)
	
	if !startDate.IsZero() {
		query = query.Where("date >= ?", startDate.Format("2006-01-02"))
	}
	if !endDate.IsZero() {
		query = query.Where("date <= ?", endDate.Format("2006-01-02"))
	}
	
	if err := query.Order("date ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	
	results := make([]domain.EventEngagementDaily, 0, len(models))
	for _, m := range models {
		results = append(results, domain.EventEngagementDaily{
			ID:                 m.ID,
			EventID:            m.EventID,
			Date:               m.Date,
			Visitors:           m.Visitors,
			EventViews:         m.EventViews,
			TicketsSelected:    m.TicketsSelected,
			CheckoutStarted:    m.CheckoutStarted,
			SuccessfulBookings: m.SuccessfulBookings,
			CreatedAt:          m.CreatedAt,
		})
	}
	return results, nil
}

func (r *engagementGormRepository) GetEngagementReport(ctx context.Context, eventIDs []string, startDate, endDate time.Time) (*domain.EngagementReportStats, error) {
	var stats domain.EngagementReportStats

	if len(eventIDs) == 0 {
		return &stats, nil
	}

	// Fetch daily aggregates for all organizer events in the date range
	var models []EventEngagementDailyModel
	query := r.db.WithContext(ctx).Where("event_id IN ?", eventIDs)
	if !startDate.IsZero() {
		query = query.Where("date >= ?", startDate.Format("2006-01-02"))
	}
	if !endDate.IsZero() {
		query = query.Where("date <= ?", endDate.Format("2006-01-02"))
	}
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	// Total aggregates for the current period
	var totalVisitors, totalEventViews, totalTicketsSelected, totalCheckout, totalBookings int
	// Day-of-week buckets: 0=Sun,1=Mon,...,6=Sat
	viewingByDow := make([]int, 7)
	checkoutByDow := make([]int, 7)

	for _, m := range models {
		totalVisitors += m.Visitors
		totalEventViews += m.EventViews
		totalTicketsSelected += m.TicketsSelected
		totalCheckout += m.CheckoutStarted
		totalBookings += m.SuccessfulBookings

		dow := int(m.Date.Weekday()) // 0=Sun
		viewingByDow[dow] += m.EventViews
		checkoutByDow[dow] += m.CheckoutStarted
	}

	// --- Previous period for percentage deltas ---
	duration := endDate.Sub(startDate)
	prevEnd := startDate
	prevStart := prevEnd.Add(-duration)

	var prevModels []EventEngagementDailyModel
	r.db.WithContext(ctx).
		Where("event_id IN ? AND date >= ? AND date <= ?", eventIDs, prevStart.Format("2006-01-02"), prevEnd.Format("2006-01-02")).
		Find(&prevModels)

	var prevVisitors, prevCheckout int
	for _, m := range prevModels {
		prevVisitors += m.Visitors
		prevCheckout += m.CheckoutStarted
	}

	// Page Views stat card
	pvPct := 0.0
	if prevVisitors > 0 {
		pvPct = float64(totalVisitors-prevVisitors) / float64(prevVisitors) * 100
	}
	stats.PageViews = domain.StatCard{Value: float64(totalVisitors), Percentage: pvPct}

	// Conversion Rate = (Bookings / Visitors) * 100
	convRate := 0.0
	if totalVisitors > 0 {
		convRate = float64(totalBookings) / float64(totalVisitors) * 100
	}
	prevConvRate := 0.0
	if prevVisitors > 0 {
		prevConvRate = float64(prevCheckout) / float64(prevVisitors) * 100
	}
	convPct := 0.0
	if prevConvRate > 0 {
		convPct = (convRate - prevConvRate) / prevConvRate * 100
	}
	stats.ConversionRate = domain.StatCard{Value: convRate, Percentage: convPct}

	// Funnel steps
	pct := func(numerator, denominator int) float64 {
		if denominator == 0 {
			return 0
		}
		return float64(numerator) / float64(denominator) * 100
	}
	stats.UserJourney = []domain.FunnelStep{
		{Label: "Visitors", Count: totalVisitors, Percentage: 100},
		{Label: "Event Page Views", Count: totalEventViews, Percentage: pct(totalEventViews, totalVisitors)},
		{Label: "Selected Tickets", Count: totalTicketsSelected, Percentage: pct(totalTicketsSelected, totalVisitors)},
		{Label: "Checkout", Count: totalCheckout, Percentage: pct(totalCheckout, totalVisitors)},
	}

	// Day-of-week peak usage (Mon-Sun ordering for display)
	dowLabels := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	dowOrder := []int{1, 2, 3, 4, 5, 6, 0} // Go weekday: 0=Sun, 1=Mon...
	for i, dowIdx := range dowOrder {
		stats.PeakUsage = append(stats.PeakUsage, domain.PeakUsagePoint{
			Label:    dowLabels[i],
			Viewing:  viewingByDow[dowIdx],
			Checkout: checkoutByDow[dowIdx],
		})
	}

	return &stats, nil
}
