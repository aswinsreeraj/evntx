package repository

import (
	"context"
	"fmt"
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
	ID        string    `gorm:"type:uuid;primaryKey"`
	UserID    *string   `gorm:"type:uuid"`
	SessionID string    `gorm:"type:uuid;index"`
	EventID   *string   `gorm:"type:uuid;index"`
	EventType string    `gorm:"index"`
	Metadata  string    `gorm:"type:json"`
	IPAddress string    `gorm:"type:varchar"`
	UserAgent string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"index"`
}

func (EngagementEventModel) TableName() string {
	return "engagement_events"
}

type EventEngagementDailyModel struct {
	ID                 string    `gorm:"type:uuid;primaryKey"`
	EventID            string    `gorm:"type:uuid;uniqueIndex:idx_event_date"`
	Date               time.Time `gorm:"type:date;uniqueIndex:idx_event_date"`
	Visitors           int
	PageViews          int
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

		dateStr := event.CreatedAt.Format("2006-01-02")
		dateParsed, _ := time.Parse("2006-01-02", dateStr)

		isUniquePlatform := r.isFirstInteractionToday(tx, event.SessionID, domain.PlatformEventID, dateParsed)
		if err := r.upsertDaily(tx, domain.PlatformEventID, dateParsed, event.EventType, isUniquePlatform); err != nil {
			return err
		}

		if event.EventID != nil && *event.EventID != "" {
			isUniqueEvent := r.isFirstInteractionToday(tx, event.SessionID, *event.EventID, dateParsed)
			if err := r.upsertDaily(tx, *event.EventID, dateParsed, event.EventType, isUniqueEvent); err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *engagementGormRepository) isFirstInteractionToday(tx *gorm.DB, sessionID, eventID string, date time.Time) bool {
	var count int64
	startOfDay := date
	endOfDay := date.Add(24 * time.Hour)

	query := tx.Model(&EngagementEventModel{}).
		Where("session_id = ? AND created_at >= ? AND created_at < ?", sessionID, startOfDay, endOfDay)

	if eventID != domain.PlatformEventID {
		query = query.Where("event_id = ?", eventID)
	}

	query.Count(&count)
	return count <= 1
}

func (r *engagementGormRepository) upsertDaily(tx *gorm.DB, eventID string, date time.Time, eventType domain.EngagementEventType, isUnique bool) error {
	upsertModel := EventEngagementDailyModel{
		ID:        uuid.NewString(),
		EventID:   eventID,
		Date:      date,
		CreatedAt: time.Now(),
	}

	updates := make(map[string]interface{})
	incExpr := func(col string) {
		updates[col] = gorm.Expr(fmt.Sprintf("event_engagement_daily.%s + 1", col))
	}

	switch eventType {
	case domain.EngagementEventEventView:
		upsertModel.EventViews = 1
		upsertModel.PageViews = 1
		incExpr("event_views")
		incExpr("page_views")
	case domain.EngagementEventTicketSelected:
		upsertModel.TicketsSelected = 1
		incExpr("tickets_selected")
	case domain.EngagementEventCheckoutStarted:
		upsertModel.CheckoutStarted = 1
		incExpr("checkout_started")
	case domain.EngagementEventPageView:
		upsertModel.PageViews = 1
		incExpr("page_views")
	}

	if isUnique {
		upsertModel.Visitors = 1
		incExpr("visitors")
	}

	if len(updates) > 0 {
		return tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "event_id"}, {Name: "date"}},
			DoUpdates: clause.Assignments(updates),
		}).Create(&upsertModel).Error
	}
	return nil
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
	query := r.db.WithContext(ctx).Model(&EventEngagementDailyModel{}).Where("event_id = ?", eventID)

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
			PageViews:          m.PageViews,
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

	loc, _ := time.LoadLocation("Asia/Calcutta")
	var models []EventEngagementDailyModel
	query := r.db.WithContext(ctx).Model(&EventEngagementDailyModel{}).Where("event_id IN (?)", eventIDs)
	if !startDate.IsZero() {
		query = query.Where("date >= ?", startDate.In(loc).Format("2006-01-02"))
	}
	if !endDate.IsZero() {
		query = query.Where("date <= ?", endDate.In(loc).Format("2006-01-02"))
	}
	
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	var totalVisitors, totalPageViews, totalEventViews, totalTicketsSelected, totalCheckout, totalBookings int
	viewingByDow := make([]int, 7)
	checkoutByDow := make([]int, 7)
	bookingsByDow := make([]int, 7)

	
	hasPlatformID := false
	for _, id := range eventIDs {
		if id == domain.PlatformEventID {
			hasPlatformID = true
			break
		}
	}

	for _, m := range models {
		if hasPlatformID && m.EventID == domain.PlatformEventID {
			
			totalVisitors = m.Visitors
			totalPageViews = m.PageViews
			totalBookings = m.SuccessfulBookings
		} else if !hasPlatformID {
			
			totalVisitors += m.Visitors
			totalPageViews += m.PageViews
			totalBookings += m.SuccessfulBookings
		}

		
		if m.EventID != domain.PlatformEventID {
			totalEventViews += m.EventViews
			totalTicketsSelected += m.TicketsSelected
			totalCheckout += m.CheckoutStarted

			dow := int(m.Date.Weekday())
			viewingByDow[dow] += m.EventViews
			checkoutByDow[dow] += m.CheckoutStarted
			bookingsByDow[dow] += m.SuccessfulBookings
		}
	}

	
	duration := endDate.Sub(startDate)
	prevEnd := startDate
	prevStart := prevEnd.Add(-duration)

	var prevModels []EventEngagementDailyModel
	r.db.WithContext(ctx).Model(&EventEngagementDailyModel{}).
		Where("event_id IN ? AND date >= ? AND date <= ?", eventIDs, prevStart.In(loc).Format("2006-01-02"), prevEnd.In(loc).Format("2006-01-02")).
		Find(&prevModels)

	var prevVisitors, prevBookings, prevPageViews int
	for _, m := range prevModels {
		if hasPlatformID && m.EventID == domain.PlatformEventID {
			prevVisitors = m.Visitors
			prevBookings = m.SuccessfulBookings
			prevPageViews = m.PageViews
		} else if !hasPlatformID {
			prevVisitors += m.Visitors
			prevBookings += m.SuccessfulBookings
			prevPageViews += m.PageViews
		}
	}

	
	pvPct := 0.0
	if prevPageViews > 0 {
		pvPct = float64(totalPageViews-prevPageViews) / float64(prevPageViews) * 100 
	} else if totalPageViews > 0 {
		pvPct = 100.0
	}
	stats.PageViews = domain.StatCard{Value: float64(totalPageViews), Percentage: pvPct}

	
	convRate := 0.0
	if totalVisitors > 0 {
		convRate = float64(totalBookings) / float64(totalVisitors) * 100
	}
	prevConvRate := 0.0
	if prevVisitors > 0 {
		prevConvRate = float64(prevBookings) / float64(prevVisitors) * 100
	}
	convPct := 0.0
	if prevConvRate > 0 {
		convPct = (convRate - prevConvRate) / prevConvRate * 100
	}
	stats.ConversionRate = domain.StatCard{Value: convRate, Percentage: convPct}

	
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
		{Label: "Checkout Started", Count: totalCheckout, Percentage: pct(totalCheckout, totalVisitors)},
		{Label: "Successful Bookings", Count: totalBookings, Percentage: pct(totalBookings, totalVisitors)},
	}

	
	dowLabels := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	dowOrder := []int{1, 2, 3, 4, 5, 6, 0} 
	for i, dowIdx := range dowOrder {
		stats.PeakUsage = append(stats.PeakUsage, domain.PeakUsagePoint{
			Label:    dowLabels[i],
			Viewing:  viewingByDow[dowIdx],
			Bookings: bookingsByDow[dowIdx],
		})
	}

	return &stats, nil
}
