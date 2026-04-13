package repository

import (
	"context"
	"math"
	"strings"
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
	apiErrors "github.com/aswinsreeraj/evntx/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventModel struct {
	ID                string
	OrganizerID       string
	Title             string
	Slug              string
	City              string
	VenueName         string
	Category          string
	StartTime         int64
	EndTime           int64
	Tags              string
	Status            string
	CoverImageURL     string
	MinPrice          float64 `gorm:"->"`
	AvailableCapacity int     `gorm:"->"`
	Settled           bool
	RejectionReason   string `gorm:"->"`
	CancellationRequestReason string `gorm:"->"`
	CreatedAt         int64
	UpdatedAt         int64
}

type EventDetailsModel struct {
	EventID            string
	Description        string
	VenueAddress       string
	MapURL             string
	TotalCapacity      int
	AvailableCapacity  int
	Rating             float64
	TotalReviews       int
	TermsAndConditions string
	CreatedAt          int64
	UpdatedAt          int64
}

type EventPersonnelModel struct {
	ID          string
	EventID     string
	Name        string
	Role        string
	Image       string
	ProfileLink string
}

type TicketTypeModel struct {
	ID                string
	EventID           string
	Name              string
	Price             float64
	TotalQuantity     int
	AvailableQuantity int
	Version           int
	CreatedAt         int64
	UpdatedAt         int64
}

type EventModerationLogModel struct {
	ID        string
	EventID   string
	AdminID   string
	Action    string
	Reason    string
	CreatedAt int64
}

type eventGormRepository struct {
	db *gorm.DB
}

func NewEventGormRepository(db *gorm.DB) *eventGormRepository {
	return &eventGormRepository{db: db}
}

func (r *eventGormRepository) ListLiveEvents(city, category, search, sortBy, minPrice, maxPrice, startDate, endDate string, page int, limit int) ([]domain.Event, int64, float64, float64, error) {

	var models []EventModel
	var total int64
	var globalMin, globalMax float64

	r.db.Model(&TicketTypeModel{}).
		Joins("JOIN event_models ON event_models.id::uuid = ticket_type_models.event_id::uuid").
		Where("event_models.status IN ?", []string{"live", "approved"}).
		Select("COALESCE(MIN(ticket_type_models.price), 0), COALESCE(MAX(ticket_type_models.price), 0)").
		Row().Scan(&globalMin, &globalMax)

	query := r.db.Model(&EventModel{}).
		Select("event_models.*, COALESCE((SELECT MIN(price) FROM ticket_type_models WHERE event_id = event_models.id), 0) as min_price, COALESCE((SELECT SUM(available_quantity) FROM ticket_type_models WHERE event_id = event_models.id), 0) as available_capacity").
		Where("status IN ?", []string{"live", "approved"})

	if city != "" {
		cities := strings.Split(city, ",")
		for i := range cities {
			cities[i] = strings.TrimSpace(cities[i])
		}
		query = query.Where("city IN ?", cities)
	}
	if category != "" && category != "All" && category != "all" {
		categories := strings.Split(category, ",")
		for i := range categories {
			categories[i] = strings.TrimSpace(categories[i])
		}
		query = query.Where("category IN ?", categories)
	}
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("title ILIKE ? OR venue_name ILIKE ? OR tags ILIKE ? OR city ILIKE ? OR category ILIKE ?", searchPattern, searchPattern, searchPattern, searchPattern, searchPattern)
	}

	if minPrice != "" || maxPrice != "" {
		subquery := r.db.Model(&TicketTypeModel{}).Select("event_id")
		if minPrice != "" {
			subquery = subquery.Where("price >= ?", minPrice)
		}
		if maxPrice != "" {
			subquery = subquery.Where("price <= ?", maxPrice)
		}
		query = query.Where("id IN (?)", subquery)
	}

	if startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("start_time >= ?", t.Unix())
		} else if t, err := time.Parse(time.RFC3339, startDate); err == nil {
			query = query.Where("start_time >= ?", t.Unix())
		}
	}
	if endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			query = query.Where("start_time <= ?", t.Add(24*time.Hour-time.Second).Unix())
		} else if t, err := time.Parse(time.RFC3339, endDate); err == nil {
			query = query.Where("start_time <= ?", t.Unix())
		}
	}

	query.Count(&total)

	if total == 0 && city != "" {
		query = r.db.Model(&EventModel{}).Where("status IN ?", []string{"live", "approved"})

		if category != "" && category != "All" && category != "all" {
			categories := strings.Split(category, ",")
			for i := range categories {
				categories[i] = strings.TrimSpace(categories[i])
			}
			query = query.Where("category IN ?", categories)
		}

		if search != "" {
			searchPattern := "%" + search + "%"
			query = query.Where("title ILIKE ? OR venue_name ILIKE ? OR tags ILIKE ? OR city ILIKE ? OR category ILIKE ?", searchPattern, searchPattern, searchPattern, searchPattern, searchPattern)
		}

		if minPrice != "" || maxPrice != "" {
			subquery := r.db.Model(&TicketTypeModel{}).Select("event_id")
			if minPrice != "" {
				subquery = subquery.Where("price >= ?", minPrice)
			}
			if maxPrice != "" {
				subquery = subquery.Where("price <= ?", maxPrice)
			}
			query = query.Where("id IN (?)", subquery)
		}

		if startDate != "" {
			if t, err := time.Parse("2006-01-02", startDate); err == nil {
				query = query.Where("start_time >= ?", t.Unix())
			} else if t, err := time.Parse(time.RFC3339, startDate); err == nil {
				query = query.Where("start_time >= ?", t.Unix())
			}
		}
		if endDate != "" {
			if t, err := time.Parse("2006-01-02", endDate); err == nil {
				query = query.Where("start_time <= ?", t.Add(24*time.Hour-time.Second).Unix())
			} else if t, err := time.Parse(time.RFC3339, endDate); err == nil {
				query = query.Where("start_time <= ?", t.Unix())
			}
		}

		query.Count(&total)
	}

	offset := (page - 1) * limit

	orderStr := "start_time ASC"
	switch sortBy {
	case "date_asc":
		orderStr = "start_time ASC"
	case "date_desc":
		orderStr = "start_time DESC"
	case "created_at_desc":
		orderStr = "created_at DESC"
	}

	err := query.
		Order(orderStr).
		Limit(limit).
		Offset(offset).
		Find(&models).Error

	if err != nil {
		return nil, 0, 0, 0, err
	}

	events := make([]domain.Event, 0)

	for _, m := range models {
		events = append(events, domain.Event{
			ID:                m.ID,
			Title:             m.Title,
			Slug:              m.Slug,
			City:              m.City,
			VenueName:         m.VenueName,
			Category:          m.Category,
			StartTime:         time.Unix(m.StartTime, 0),
			EndTime:           time.Unix(m.EndTime, 0),
			Tags:              m.Tags,
			CoverImageURL:     m.CoverImageURL,
			MinPrice:          m.MinPrice,
			AvailableCapacity: m.AvailableCapacity,
			Settled:           m.Settled,
			CreatedAt:         time.Unix(m.CreatedAt, 0),
			UpdatedAt:         time.Unix(m.UpdatedAt, 0),
		})
	}

	return events, total, globalMin, globalMax, nil
}

func (r *eventGormRepository) AdminSearchEvents(
	search string,
	status string,
	page int,
	limit int,
) ([]domain.AdminEventDetails, int64, error) {

	var models []struct {
		EventModel
		OrganizerName string
		TicketsSold   int64
		Revenue       float64
	}
	var total int64

	query := r.db.Table("event_models").
		Select(`
			event_models.*,
			user_models.name AS organizer_name,
			COALESCE((
				SELECT SUM(bkt.quantity)
				FROM booking_models bk
				JOIN booking_ticket_models bkt ON bkt.booking_id = bk.id
				WHERE bk.event_id = event_models.id AND bk.status IN ('paid', 'confirmed')
			), 0) AS tickets_sold,
			COALESCE((
				SELECT SUM(bk.total_amount)
				FROM booking_models bk
				WHERE bk.event_id = event_models.id AND bk.status IN ('paid', 'confirmed')
			), 0) AS revenue,
			COALESCE((
				SELECT SUM(available_quantity)
				FROM ticket_type_models
				WHERE event_id = event_models.id
			), 0) AS available_capacity,
			(
				SELECT reason
				FROM event_moderation_log_models
				WHERE event_id = event_models.id AND action = 'cancellation_requested'
				ORDER BY created_at DESC
				LIMIT 1
			) AS cancellation_request_reason
		`).
		Joins("LEFT JOIN user_models ON user_models.id::uuid = event_models.organizer_id::uuid")

	if search != "" {
		query = query.Where(
			"event_models.title ILIKE ? OR user_models.name ILIKE ?",
			"%"+search+"%",
			"%"+search+"%",
		)
	}

	if status != "" && status != "all" {
		query = query.Where("event_models.status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * limit
	err := query.
		Order("event_models.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&models).Error

	if err != nil {
		return nil, 0, err
	}

	details := make([]domain.AdminEventDetails, 0, len(models))
	for _, m := range models {
		details = append(details, domain.AdminEventDetails{
			Event: domain.Event{
				ID:        m.ID,
				Title:     m.Title,
				City:      m.City,
				StartTime: time.Unix(m.StartTime, 0),
				Status:    m.Status,
				Settled:   m.Settled,
				CancellationRequestReason: m.CancellationRequestReason,
				CreatedAt: time.Unix(m.CreatedAt, 0),
				UpdatedAt: time.Unix(m.UpdatedAt, 0),
			},
			OrganizerName: m.OrganizerName,
			TicketsSold:   m.TicketsSold,
			Revenue:       int64(m.Revenue),
		})
	}

	return details, total, nil
}

func (r *eventGormRepository) GetEventBySlug(slug string) (*domain.Event, error) {

	var model EventModel

	err := r.db.Model(&EventModel{}).
		Select("event_models.*, COALESCE((SELECT SUM(available_quantity) FROM ticket_type_models WHERE event_id = event_models.id), 0) as available_capacity").
		Where("slug = ? OR id = ?", slug, slug).
		First(&model).Error

	if err != nil {
		return nil, err
	}

	return &domain.Event{
		ID:                model.ID,
		OrganizerID:       model.OrganizerID,
		Title:             model.Title,
		Slug:              model.Slug,
		Status:            model.Status,
		City:              model.City,
		VenueName:         model.VenueName,
		Category:          model.Category,
		StartTime:         time.Unix(model.StartTime, 0),
		EndTime:           time.Unix(model.EndTime, 0),
		Tags:              model.Tags,
		CoverImageURL:     model.CoverImageURL,
		AvailableCapacity: model.AvailableCapacity,
		Settled:           model.Settled,
		CancellationRequestReason: model.CancellationRequestReason,
		CreatedAt:         time.Unix(model.CreatedAt, 0),
		UpdatedAt:         time.Unix(model.UpdatedAt, 0),
	}, nil
}

func (r *eventGormRepository) GetEventByID(eventID string) (*domain.Event, error) {

	var model EventModel

	err := r.db.Model(&EventModel{}).
		Select("event_models.*, COALESCE((SELECT SUM(available_quantity) FROM ticket_type_models WHERE event_id = event_models.id), 0) as available_capacity").
		Where("id = ?", eventID).
		First(&model).Error

	if err != nil {
		return nil, err
	}

	return &domain.Event{
		ID:                model.ID,
		OrganizerID:       model.OrganizerID,
		Title:             model.Title,
		Slug:              model.Slug,
		Status:            model.Status,
		City:              model.City,
		VenueName:         model.VenueName,
		Category:          model.Category,
		StartTime:         time.Unix(model.StartTime, 0),
		EndTime:           time.Unix(model.EndTime, 0),
		Tags:              model.Tags,
		CoverImageURL:     model.CoverImageURL,
		AvailableCapacity: model.AvailableCapacity,
		Settled:           model.Settled,
		CancellationRequestReason: model.CancellationRequestReason,
		CreatedAt:         time.Unix(model.CreatedAt, 0),
		UpdatedAt:         time.Unix(model.UpdatedAt, 0),
	}, nil
}

func (r *eventGormRepository) GetEventDetails(eventID string) (*domain.EventDetails, error) {

	var model EventDetailsModel

	err := r.db.Model(&EventDetailsModel{}).
		Select("event_details_models.*, COALESCE((SELECT SUM(available_quantity) FROM ticket_type_models WHERE event_id = event_details_models.event_id), 0) as available_capacity").
		Where("event_id = ?", eventID).
		First(&model).Error

	if err != nil {
		return nil, err
	}

	return &domain.EventDetails{
		EventID:            model.EventID,
		Description:        model.Description,
		VenueAddress:       model.VenueAddress,
		MapURL:             model.MapURL,
		TotalCapacity:      model.TotalCapacity,
		AvailableCapacity:  model.AvailableCapacity,
		Rating:             model.Rating,
		TotalReviews:       model.TotalReviews,
		TermsAndConditions: model.TermsAndConditions,
		CreatedAt:          time.Unix(model.CreatedAt, 0),
		UpdatedAt:          time.Unix(model.UpdatedAt, 0),
	}, nil
}

func (r *eventGormRepository) GetEventPersonnels(eventID string) ([]domain.EventPersonnel, error) {

	var models []EventPersonnelModel

	err := r.db.
		Where("event_id = ?", eventID).
		Find(&models).Error

	if err != nil {
		return nil, err
	}

	personnels := make([]domain.EventPersonnel, 0)

	for _, m := range models {
		personnels = append(personnels, domain.EventPersonnel{
			ID:          m.ID,
			EventID:     m.EventID,
			Name:        m.Name,
			Role:        m.Role,
			Image:       m.Image,
			ProfileLink: m.ProfileLink,
		})
	}

	return personnels, nil
}

func (r *eventGormRepository) GetTicketTypesByEventID(eventID string) ([]domain.TicketType, error) {
	var models []TicketTypeModel

	err := r.db.Where("event_id = ?", eventID).Find(&models).Error
	if err != nil {
		return nil, err
	}

	tickets := make([]domain.TicketType, 0, len(models))
	for _, m := range models {
		tickets = append(tickets, domain.TicketType{
			ID:                m.ID,
			EventID:           m.EventID,
			Name:              m.Name,
			Price:             m.Price,
			TotalQuantity:     m.TotalQuantity,
			AvailableQuantity: m.AvailableQuantity,
			Version:           m.Version,
		})
	}

	return tickets, nil
}

func (r *eventGormRepository) CreateEvent(ctx context.Context, event *domain.Event, details *domain.EventDetails, tickets []domain.TicketType, personnels []domain.EventPersonnel) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		eventModel := EventModel{
			ID:            event.ID,
			OrganizerID:   event.OrganizerID,
			Title:         event.Title,
			Slug:          event.Slug,
			City:          event.City,
			VenueName:     event.VenueName,
			Category:      event.Category,
			StartTime:     event.StartTime.Unix(),
			EndTime:       event.EndTime.Unix(),
			Tags:          event.Tags,
			Status:        event.Status,
			Settled:       event.Settled,
			CoverImageURL: event.CoverImageURL,
			CreatedAt:     event.CreatedAt.Unix(),
			UpdatedAt:     event.UpdatedAt.Unix(),
		}

		if err := tx.Create(&eventModel).Error; err != nil {
			return err
		}

		detailsModel := EventDetailsModel{
			EventID:            details.EventID,
			Description:        details.Description,
			VenueAddress:       details.VenueAddress,
			MapURL:             details.MapURL,
			TotalCapacity:      details.TotalCapacity,
			AvailableCapacity:  details.AvailableCapacity,
			Rating:             details.Rating,
			TotalReviews:       details.TotalReviews,
			TermsAndConditions: details.TermsAndConditions,
			CreatedAt:          details.CreatedAt.Unix(),
			UpdatedAt:          details.UpdatedAt.Unix(),
		}

		if err := tx.Create(&detailsModel).Error; err != nil {
			return err
		}

		for _, ticket := range tickets {
			ticketModel := TicketTypeModel{
				ID:                ticket.ID,
				EventID:           ticket.EventID,
				Name:              ticket.Name,
				Price:             ticket.Price,
				TotalQuantity:     ticket.TotalQuantity,
				AvailableQuantity: ticket.AvailableQuantity,
				Version:           ticket.Version,
				CreatedAt:         ticket.CreatedAt.Unix(),
				UpdatedAt:         ticket.UpdatedAt.Unix(),
			}

			if err := tx.Create(&ticketModel).Error; err != nil {
				return err
			}
		}

		for _, p := range personnels {
			pModel := EventPersonnelModel{
				ID:          p.ID,
				EventID:     p.EventID,
				Name:        p.Name,
				Role:        p.Role,
				Image:       p.Image,
				ProfileLink: p.ProfileLink,
			}
			if err := tx.Create(&pModel).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *eventGormRepository) UpdateEvent(ctx context.Context, eventID string, eventUpdates map[string]interface{}, detailUpdates map[string]interface{}, ticketUpdates []domain.TicketType, personnelUpdates []domain.EventPersonnel) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if updatedAt, ok := eventUpdates["updated_at"].(time.Time); ok {
			eventUpdates["updated_at"] = updatedAt.Unix()
		}
		if updatedAt, ok := detailUpdates["updated_at"].(time.Time); ok {
			detailUpdates["updated_at"] = updatedAt.Unix()
		}

		if len(eventUpdates) > 0 {
			if err := tx.Model(&EventModel{}).Where("id = ?", eventID).Updates(eventUpdates).Error; err != nil {
				return err
			}
		}

		if len(detailUpdates) > 0 {
			if err := tx.Model(&EventDetailsModel{}).Where("event_id = ?", eventID).Updates(detailUpdates).Error; err != nil {
				return err
			}
		}

		for _, ticket := range ticketUpdates {
			model := TicketTypeModel{
				ID:                ticket.ID,
				EventID:           ticket.EventID,
				Name:              ticket.Name,
				Price:             ticket.Price,
				TotalQuantity:     ticket.TotalQuantity,
				AvailableQuantity: ticket.AvailableQuantity,
				Version:           ticket.Version,
				CreatedAt:         ticket.CreatedAt.Unix(),
				UpdatedAt:         ticket.UpdatedAt.Unix(),
			}

			if err := tx.Save(&model).Error; err != nil {
				return err
			}
		}

		if personnelUpdates != nil {
			if err := tx.Where("event_id = ?", eventID).Delete(&EventPersonnelModel{}).Error; err != nil {
				return err
			}
			for _, p := range personnelUpdates {
				model := EventPersonnelModel{
					ID:          p.ID,
					EventID:     p.EventID,
					Name:        p.Name,
					Role:        p.Role,
					Image:       p.Image,
					ProfileLink: p.ProfileLink,
				}
				if err := tx.Create(&model).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (r *eventGormRepository) UpdateEventStatus(ctx context.Context, eventID string, status string) error {
	return r.db.WithContext(ctx).Model(&EventModel{}).
		Where("id = ?", eventID).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": gorm.Expr("EXTRACT(EPOCH FROM NOW())"),
		}).Error
}

func (r *eventGormRepository) SettleEventEarnings(
	ctx context.Context,
	eventID string,
	organizerID string,
	totalAmountParam float64,
) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		eventResult := tx.Model(&EventModel{}).
			Where("id = ? AND status = ? AND settled = ?", eventID, "completed", false).
			Updates(map[string]interface{}{
				"settled":    true,
				"updated_at": gorm.Expr("EXTRACT(EPOCH FROM NOW())"),
			})
		if eventResult.Error != nil {
			return eventResult.Error
		}
		if eventResult.RowsAffected == 0 {
			return apiErrors.ErrDuplicateResource
		}

		var bookingStats struct {
			TotalAmount  float64
			TotalTickets int64
		}

		if err := tx.Table("booking_models").
			Select("COALESCE(SUM(booking_models.total_amount), 0) as total_amount, COALESCE(SUM(booking_ticket_models.quantity), 0) as total_tickets").
			Joins("LEFT JOIN booking_ticket_models ON booking_ticket_models.booking_id::uuid = booking_models.id::uuid").
			Where("booking_models.event_id = ? AND booking_models.status IN ('paid', 'confirmed')", eventID).
			Scan(&bookingStats).Error; err != nil {
			return err
		}

		if bookingStats.TotalAmount <= 0 {
			return nil
		}

		baseT := bookingStats.TotalAmount - float64(bookingStats.TotalTickets*30)
		organizerFee := math.Round((baseT*0.05)*100) / 100
		organizerRevenue := math.Round((baseT-organizerFee)*100) / 100
		reserveRelease := float64(bookingStats.TotalTickets * 30)
		now := time.Now()

		var wallet WalletModel
		if err := tx.Where("user_id = ?", organizerID).First(&wallet).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return apiErrors.ErrResourceNotFound
			}
			return err
		}

		if wallet.PendingBalance < baseT {
			return apiErrors.ErrInsufficientBalance
		}

		if err := tx.Create(&WalletTransactionModel{
			ID:            uuid.NewString(),
			WalletID:      wallet.ID,
			Type:          domain.WalletTransactionTypeCredit,
			Amount:        organizerRevenue,
			ReferenceType: domain.WalletReferenceTypeSettlement,
			ReferenceID:   eventID,
			Status:        domain.WalletTransactionStatusCompleted,
			CreatedAt:     now,
		}).Error; err != nil {
			return err
		}

		wallet.PendingBalance = math.Round((wallet.PendingBalance-baseT)*100) / 100
		wallet.AvailableBalance = math.Round((wallet.AvailableBalance+organizerRevenue)*100) / 100
		wallet.ReserveBalance = math.Round((wallet.ReserveBalance+reserveRelease)*100) / 100
		wallet.UpdatedAt = now

		if err := tx.Model(&WalletModel{}).
			Where("id = ?", wallet.ID).
			Select("pending_balance", "available_balance", "reserve_balance", "updated_at").
			Updates(WalletModel{
				PendingBalance:   wallet.PendingBalance,
				AvailableBalance: wallet.AvailableBalance,
				ReserveBalance:   wallet.ReserveBalance,
				UpdatedAt:        wallet.UpdatedAt,
			}).Error; err != nil {
			return err
		}

		var platformWallet PlatformWalletModel
		if err := tx.Where("id = ?", domain.PlatformWalletID).First(&platformWallet).Error; err != nil {
			return err
		}

		if err := tx.Create(&PlatformWalletTransactionModel{
			ID:            uuid.NewString(),
			WalletID:      domain.PlatformWalletID,
			Type:          domain.WalletTransactionTypeCredit,
			Amount:        organizerFee,
			ReferenceType: domain.PlatformRefTypeEarning,
			ReferenceID:   eventID,
			CreatedAt:     now,
		}).Error; err != nil {
			return err
		}

		platformWallet.AvailableBalance = math.Round((platformWallet.AvailableBalance+organizerFee)*100) / 100
		platformWallet.TotalCredited = math.Round((platformWallet.TotalCredited+organizerFee)*100) / 100
		platformWallet.UpdatedAt = now

		if err := tx.Model(&PlatformWalletModel{}).
			Where("id = ?", domain.PlatformWalletID).
			Select("available_balance", "total_credited", "updated_at").
			Updates(PlatformWalletModel{
				AvailableBalance: platformWallet.AvailableBalance,
				TotalCredited:    platformWallet.TotalCredited,
				UpdatedAt:        platformWallet.UpdatedAt,
			}).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *eventGormRepository) ApproveEvent(ctx context.Context, eventID string) error {
	return r.db.WithContext(ctx).Model(&EventModel{}).
		Where("id = ?", eventID).
		Updates(map[string]interface{}{
			"status":     "approved",
			"updated_at": gorm.Expr("EXTRACT(EPOCH FROM NOW())"),
		}).Error
}

func (r *eventGormRepository) RejectEvent(ctx context.Context, eventID string, adminID string, reason string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&EventModel{}).
			Where("id = ?", eventID).
			Updates(map[string]interface{}{
				"status":     "rejected",
				"updated_at": gorm.Expr("EXTRACT(EPOCH FROM NOW())"),
			}).Error; err != nil {
			return err
		}

		logModel := EventModerationLogModel{
			ID:        uuid.New().String(),
			EventID:   eventID,
			AdminID:   adminID,
			Action:    "rejected",
			Reason:    reason,
			CreatedAt: time.Now().Unix(),
		}

		return tx.Create(&logModel).Error
	})
}

func (r *eventGormRepository) SuspendLiveEvent(ctx context.Context, eventID string, adminID string, reason string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&EventModel{}).
			Where("id = ?", eventID).
			Updates(map[string]interface{}{
				"status":     "suspended",
				"updated_at": gorm.Expr("EXTRACT(EPOCH FROM NOW())"),
			}).Error; err != nil {
			return err
		}

		logModel := EventModerationLogModel{
			ID:        uuid.New().String(),
			EventID:   eventID,
			AdminID:   adminID,
			Action:    "suspended",
			Reason:    reason,
			CreatedAt: time.Now().Unix(),
		}

		return tx.Create(&logModel).Error
	})
}

func (r *eventGormRepository) GetEventsByOrganizerID(organizerID string, status string) ([]domain.Event, error) {
	var models []EventModel
	query := r.db.Model(&EventModel{}).
		Select("event_models.*, "+
			"(SELECT reason FROM event_moderation_log_models WHERE event_id = event_models.id AND action IN ('rejected', 'suspended', 'cancellation_rejected') ORDER BY created_at DESC LIMIT 1) as rejection_reason, "+
			"(SELECT reason FROM event_moderation_log_models WHERE event_id = event_models.id AND action = 'cancellation_requested' ORDER BY created_at DESC LIMIT 1) as cancellation_request_reason, "+
			"COALESCE((SELECT SUM(available_quantity) FROM ticket_type_models WHERE event_id = event_models.id), 0) as available_capacity").
		Where("organizer_id = ?", organizerID)

	if status != "" && status != "All" && status != "all" {
		query = query.Where("status = ?", strings.ToLower(status))
	}

	err := query.Order("start_time DESC").Find(&models).Error
	if err != nil {
		return nil, err
	}

	events := make([]domain.Event, 0)
	for _, m := range models {
		events = append(events, domain.Event{
			ID:                m.ID,
			OrganizerID:       m.OrganizerID,
			Title:             m.Title,
			Slug:              m.Slug,
			Status:            m.Status,
			City:              m.City,
			VenueName:         m.VenueName,
			Category:          m.Category,
			StartTime:         time.Unix(m.StartTime, 0),
			EndTime:           time.Unix(m.EndTime, 0),
			Tags:              m.Tags,
			CoverImageURL:     m.CoverImageURL,
			AvailableCapacity: m.AvailableCapacity,
			Settled:           m.Settled,
			RejectionReason:   m.RejectionReason,
			CancellationRequestReason: m.CancellationRequestReason,
			CreatedAt:         time.Unix(m.CreatedAt, 0),
			UpdatedAt:         time.Unix(m.UpdatedAt, 0),
		})
	}
	return events, nil
}

func (r *eventGormRepository) DeleteEvent(ctx context.Context, eventID string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("event_id = ?", eventID).Delete(&EventPersonnelModel{}).Error; err != nil {
			return err
		}
		if err := tx.Where("event_id = ?", eventID).Delete(&TicketTypeModel{}).Error; err != nil {
			return err
		}
		if err := tx.Where("event_id = ?", eventID).Delete(&EventDetailsModel{}).Error; err != nil {
			return err
		}
		if err := tx.Where("id = ?", eventID).Delete(&EventModel{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *eventGormRepository) CancelLiveEvent(ctx context.Context, eventID string, organizerID string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		eventResult := tx.Model(&EventModel{}).
			Where("id = ? AND organizer_id = ? AND status IN ('live', 'approved', 'cancellation_pending')", eventID, organizerID).
			Updates(map[string]interface{}{
				"status":     "cancelled",
				"updated_at": gorm.Expr("EXTRACT(EPOCH FROM NOW())"),
			})
		if eventResult.Error != nil {
			return eventResult.Error
		}
		if eventResult.RowsAffected == 0 {
			return apiErrors.ErrInvalidStateTransition
		}

		var bookings []BookingModel
		if err := tx.Where("event_id = ? AND status IN ('paid', 'confirmed')", eventID).Find(&bookings).Error; err != nil {
			return err
		}

		if len(bookings) == 0 {
			return nil
		}

		now := time.Now()

		for _, bm := range bookings {
			var bookingTickets []BookingTicketModel
			if err := tx.Where("booking_id = ?", bm.ID).Find(&bookingTickets).Error; err != nil {
				return err
			}

			var totalTicketsCancelled int
			for _, bt := range bookingTickets {
				totalTicketsCancelled += bt.Quantity
			}

			baseRefund := bm.TotalAmount - float64(totalTicketsCancelled*30)
			totalRefundToUser := bm.TotalAmount

			if totalRefundToUser <= 0 {
				continue
			}

			if err := tx.Model(&BookingModel{}).Where("id = ?", bm.ID).Update("status", "cancelled").Error; err != nil {
				return err
			}
			if err := tx.Model(&TicketModel{}).Where("booking_id = ?", bm.ID).Update("status", "cancelled").Error; err != nil {
				return err
			}

			var userWallet WalletModel
			if err := tx.Where("user_id = ?", bm.UserID).First(&userWallet).Error; err != nil {
				return err
			}

			if err := tx.Create(&WalletTransactionModel{
				ID:            uuid.NewString(),
				WalletID:      userWallet.ID,
				Type:          domain.WalletTransactionTypeCredit,
				Amount:        totalRefundToUser,
				ReferenceType: domain.WalletReferenceTypeOrganizerCancellation,
				ReferenceID:   bm.ID,
				Status:        domain.WalletTransactionStatusCompleted,
				CreatedAt:     now,
			}).Error; err != nil {
				return err
			}

			userWallet.AvailableBalance = math.Round((userWallet.AvailableBalance+totalRefundToUser)*100) / 100
			if err := tx.Model(&WalletModel{}).Where("id = ?", userWallet.ID).Updates(map[string]interface{}{
				"available_balance": userWallet.AvailableBalance,
				"updated_at":        now,
			}).Error; err != nil {
				return err
			}

			var orgWallet WalletModel
			if err := tx.Where("user_id = ?", organizerID).First(&orgWallet).Error; err != nil {
				return err
			}

			platformFee := float64(totalTicketsCancelled * 30)

			if err := tx.Create(&WalletTransactionModel{
				ID:            uuid.NewString(),
				WalletID:      orgWallet.ID,
				Type:          domain.WalletTransactionTypeDebit,
				Amount:        totalRefundToUser,
				ReferenceType: domain.WalletReferenceTypeOrganizerCancellation,
				ReferenceID:   bm.ID,
				Status:        domain.WalletTransactionStatusCompleted,
				CreatedAt:     now,
			}).Error; err != nil {
				return err
			}

			orgWallet.PendingBalance = math.Round((orgWallet.PendingBalance-baseRefund)*100) / 100
			orgWallet.AvailableBalance = math.Round((orgWallet.AvailableBalance-platformFee)*100) / 100
			orgWallet.ReserveBalance = math.Round((orgWallet.ReserveBalance+platformFee)*100) / 100

			if err := tx.Model(&WalletModel{}).Where("id = ?", orgWallet.ID).Updates(map[string]interface{}{
				"pending_balance":   orgWallet.PendingBalance,
				"available_balance": orgWallet.AvailableBalance,
				"reserve_balance":   orgWallet.ReserveBalance,
				"updated_at":        now,
			}).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *eventGormRepository) RequestEventCancellation(ctx context.Context, eventID string, organizerID string, reason string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		eventResult := tx.Model(&EventModel{}).
			Where("id = ? AND organizer_id = ? AND status = 'live'", eventID, organizerID).
			Updates(map[string]interface{}{
				"status":     "cancellation_pending",
				"updated_at": gorm.Expr("EXTRACT(EPOCH FROM NOW())"),
			})
		if eventResult.Error != nil {
			return eventResult.Error
		}
		if eventResult.RowsAffected == 0 {
			return apiErrors.ErrInvalidStateTransition
		}

		logModel := EventModerationLogModel{
			ID:        uuid.New().String(),
			EventID:   eventID,
			AdminID:   organizerID,
			Action:    "cancellation_requested",
			Reason:    reason,
			CreatedAt: time.Now().Unix(),
		}
		return tx.Create(&logModel).Error
	})
}

func (r *eventGormRepository) RejectEventCancellation(ctx context.Context, eventID string, adminID string, reason string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		eventResult := tx.Model(&EventModel{}).
			Where("id = ? AND status = 'cancellation_pending'", eventID).
			Updates(map[string]interface{}{
				"status":     "live",
				"updated_at": gorm.Expr("EXTRACT(EPOCH FROM NOW())"),
			})
		if eventResult.Error != nil {
			return eventResult.Error
		}
		if eventResult.RowsAffected == 0 {
			return apiErrors.ErrInvalidStateTransition
		}

		logModel := EventModerationLogModel{
			ID:        uuid.New().String(),
			EventID:   eventID,
			AdminID:   adminID,
			Action:    "cancellation_rejected",
			Reason:    reason,
			CreatedAt: time.Now().Unix(),
		}
		return tx.Create(&logModel).Error
	})
}

func (r *eventGormRepository) FindPastLiveEvents(ctx context.Context, now time.Time) ([]domain.Event, error) {
	var models []EventModel
	nowUnix := now.Unix()

	if err := r.db.WithContext(ctx).
		Where("status = ? AND end_time < ?", "live", nowUnix).
		Find(&models).Error; err != nil {
		return nil, err
	}

	events := make([]domain.Event, len(models))
	for i, model := range models {
		events[i] = domain.Event{
			ID:            model.ID,
			OrganizerID:   model.OrganizerID,
			Title:         model.Title,
			Slug:          model.Slug,
			Status:        model.Status,
			City:          model.City,
			VenueName:     model.VenueName,
			Category:      model.Category,
			StartTime:     time.Unix(model.StartTime, 0),
			EndTime:       time.Unix(model.EndTime, 0),
			Tags:          model.Tags,
			CoverImageURL: model.CoverImageURL,
			Settled:       model.Settled,
			CreatedAt:     time.Unix(model.CreatedAt, 0),
			UpdatedAt:     time.Unix(model.UpdatedAt, 0),
		}
	}

	return events, nil
}

func (r *eventGormRepository) GetDashboardStats(organizerID string) (*domain.OrganizerDashboardStats, error) {
	var events []EventModel
	if err := r.db.Where("organizer_id = ?", organizerID).Find(&events).Error; err != nil {
		return nil, err
	}

	eventIDs := make([]string, 0, len(events))
	var activeEvents, pendingEvents int
	var activeEventsPrevMonth int
	var eventTitleMap = make(map[string]string)

	now := time.Now()
	firstDayThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	firstDayPrevMonth := firstDayThisMonth.AddDate(0, -1, 0)

	for _, e := range events {
		eventIDs = append(eventIDs, e.ID)
		eventTitleMap[e.ID] = e.Title
		
		if e.Status == "live" {
			activeEvents++
			if time.Unix(e.CreatedAt, 0).Before(firstDayThisMonth) {
			    activeEventsPrevMonth++
			}
		} else if e.Status == "pending" {
			pendingEvents++
		}
	}

	var stats domain.OrganizerDashboardStats
	if len(eventIDs) == 0 {
		return &stats, nil
	}

	var bookings []struct {
		EventID     string
		TotalAmount float64
		CreatedAt   int64 `gorm:"column:created_at"` 
		Tickets     int
	}

	r.db.Table("booking_models").
	    Select("booking_models.event_id, booking_models.total_amount, booking_models.created_at, (SELECT COALESCE(SUM(quantity), 0) FROM booking_ticket_models WHERE booking_ticket_models.booking_id = booking_models.id) as tickets").
		Where("booking_models.event_id IN ? AND booking_models.status IN ?", eventIDs, []string{"paid", "completed"}).
		Find(&bookings)

	var totalRev, totalRevPrev float64
	var totalTkt, totalTktPrev int

	var revenueByMonth = make(map[string]float64)
	for i := 11; i >= 0; i-- {
		m := now.AddDate(0, -i, 0)
		revenueByMonth[m.Format("Jan")] = 0
	}

	var revenueByEvent = make(map[string]float64)

	for _, b := range bookings {
		totalRev += b.TotalAmount
		totalTkt += b.Tickets
		
		bCreatedAt := time.Unix(b.CreatedAt, 0)
		if bCreatedAt.After(firstDayPrevMonth) && bCreatedAt.Before(firstDayThisMonth) {
			totalRevPrev += b.TotalAmount
			totalTktPrev += b.Tickets
		}

		revenueByEvent[b.EventID] += b.TotalAmount

		monthStr := bCreatedAt.Format("Jan")
		if _, exists := revenueByMonth[monthStr]; exists {
			revenueByMonth[monthStr] += b.TotalAmount
		}
	}
    
	var revPercentage, tktPercentage, activePercentage float64
    if totalRevPrev > 0 {
        revPercentage = ((totalRev - totalRevPrev) / totalRevPrev) * 100
    }
	if totalTktPrev > 0 {
        tktPercentage = float64(totalTkt - totalTktPrev) / float64(totalTktPrev) * 100
    }
	if activeEventsPrevMonth > 0 {
        activePercentage = float64(activeEvents - activeEventsPrevMonth) / float64(activeEventsPrevMonth) * 100
    }

	stats.TotalRevenue = domain.StatCard{Value: totalRev, Percentage: revPercentage}
	stats.TicketsSold = domain.StatCard{Value: float64(totalTkt), Percentage: tktPercentage}
	stats.ActiveEvents = domain.StatCard{Value: float64(activeEvents), Percentage: activePercentage}
	stats.PendingEvents = domain.StatCard{Value: float64(pendingEvents), Percentage: 0}

	for i := 11; i >= 0; i-- {
		m := now.AddDate(0, -i, 0).Format("Jan")
		stats.RevenueOverview = append(stats.RevenueOverview, domain.RevenuePoint{
			Date:   m,
			Amount: revenueByMonth[m],
		})
	}

	for evtID, val := range revenueByEvent {
        if val > 0 && eventTitleMap[evtID] != "" {
		    stats.SalesBreakdown = append(stats.SalesBreakdown, domain.EventSalesBreakdown{
			    EventName: eventTitleMap[evtID],
			    Revenue:   val,
		    })
        }
	}

	return &stats, nil
}

func (r *eventGormRepository) GetSalesReport(organizerID string, eventID string, startDate string, endDate string) (*domain.SalesReportStats, error) {
	var events []EventModel
	query := r.db.Where("organizer_id = ?", organizerID)
	if eventID != "" && eventID != "all" {
		query = query.Where("id = ?", eventID)
	}

	if err := query.Find(&events).Error; err != nil {
		return nil, err
	}

	eventIDs := make([]string, 0, len(events))
	var eventTitleMap = make(map[string]string)
	for _, e := range events {
		eventIDs = append(eventIDs, e.ID)
		eventTitleMap[e.ID] = e.Title
	}

	var stats domain.SalesReportStats
	if len(eventIDs) == 0 {
		return &stats, nil
	}

	var startTime, endTime time.Time
	var validRange bool
	if startDate != "" && endDate != "" {
		st, err1 := time.Parse(time.RFC3339, startDate)
		et, err2 := time.Parse(time.RFC3339, endDate)
		if err1 == nil && err2 == nil {
			startTime = st
			endTime = et
			validRange = true
		}
	}

	if !validRange { 
		endTime = time.Now()
		startTime = endTime.AddDate(0, 0, -30)
	}

	duration := endTime.Sub(startTime)
	prevStartTime := startTime.Add(-duration)

	var bookings []struct {
		EventID     string
		TotalAmount float64
		CreatedAt   int64 `gorm:"column:created_at"`
		Tickets     int
	}

	
	r.db.Table("booking_models").
		Select("booking_models.event_id, booking_models.total_amount, booking_models.created_at, (SELECT COALESCE(SUM(quantity), 0) FROM booking_ticket_models WHERE booking_ticket_models.booking_id = booking_models.id) as tickets").
		Where("booking_models.event_id IN ? AND booking_models.status IN ?", eventIDs, []string{"paid", "completed"}).
		Find(&bookings)

	var totalRev, totalRevPrev float64
	var totalTkt, totalTktPrev int

	
	useDays := duration.Hours() <= 31*24
	var revenueMap = make(map[string]float64)
	
	
	if useDays {
		days := int(duration.Hours() / 24)
		for i := 0; i <= days; i++ {
			d := startTime.AddDate(0, 0, i)
			revenueMap[d.Format("Jan 02")] = 0
		}
	} else {
		months := int(duration.Hours() / (24 * 30))
		for i := 0; i <= months; i++ {
			m := startTime.AddDate(0, i, 0)
			revenueMap[m.Format("Jan 2006")] = 0
		}
	}

	var ticketsByEvent = make(map[string]int)

	for _, b := range bookings {
		bCreatedAt := time.Unix(b.CreatedAt, 0)

		
		if bCreatedAt.After(startTime) && bCreatedAt.Before(endTime) {
			totalRev += b.TotalAmount
			totalTkt += b.Tickets
			ticketsByEvent[b.EventID] += b.Tickets
			
			var timeKey string
			if useDays {
				timeKey = bCreatedAt.Format("Jan 02")
			} else {
				timeKey = bCreatedAt.Format("Jan 2006")
			}
			if _, exists := revenueMap[timeKey]; exists {
				revenueMap[timeKey] += b.TotalAmount
			}
		}

		
		if bCreatedAt.After(prevStartTime) && bCreatedAt.Before(startTime) {
			totalRevPrev += b.TotalAmount
			totalTktPrev += b.Tickets
		}
	}

	var revPercentage, tktPercentage float64
	if totalRevPrev > 0 {
		revPercentage = ((totalRev - totalRevPrev) / totalRevPrev) * 100
	}
	if totalTktPrev > 0 {
		tktPercentage = float64(totalTkt - totalTktPrev) / float64(totalTktPrev) * 100
	}

	stats.TotalRevenue = domain.StatCard{Value: totalRev, Percentage: revPercentage}
	stats.TicketsSold = domain.StatCard{Value: float64(totalTkt), Percentage: tktPercentage}

	
	if useDays {
		days := int(duration.Hours() / 24)
		for i := 0; i <= days; i++ {
			d := startTime.AddDate(0, 0, i).Format("Jan 02")
			stats.RevenueOverTime = append(stats.RevenueOverTime, domain.RevenuePoint{
				Date:   d,
				Amount: revenueMap[d],
			})
		}
	} else {
		months := int(duration.Hours() / (24 * 30))
		for i := 0; i <= months; i++ {
			m := startTime.AddDate(0, i, 0).Format("Jan 2006")
			stats.RevenueOverTime = append(stats.RevenueOverTime, domain.RevenuePoint{
				Date:   m,
				Amount: revenueMap[m],
			})
		}
	}

	
	for evtID, tkts := range ticketsByEvent {
		if tkts > 0 && eventTitleMap[evtID] != "" {
			percentage := 0.0
			if totalTkt > 0 {
				percentage = float64(tkts) / float64(totalTkt) * 100
			}
			stats.TicketsPerEvent = append(stats.TicketsPerEvent, domain.TicketSalesProportion{
				EventName:       eventTitleMap[evtID],
				TicketsSold:     tkts,
				PercentageTotal: percentage,
			})
		}
	}

	return &stats, nil
}

func (r *eventGormRepository) GetAdminDashboardStats() (*domain.AdminDashboardStats, error) {
	now := time.Now()
	firstDayThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	firstDayPrevMonth := firstDayThisMonth.AddDate(0, -1, 0)
	var stats domain.AdminDashboardStats

	
	var totalRevResult struct{ Total float64 }
	r.db.Table("booking_models").Where("status IN ?", []string{"paid", "completed"}).
		Select("COALESCE(SUM(total_amount), 0) as total").Scan(&totalRevResult)

	var prevRevResult struct{ Total float64 }
	r.db.Table("booking_models").
		Where("status IN ? AND created_at >= ? AND created_at < ?", []string{"paid", "completed"}, firstDayPrevMonth.Unix(), firstDayThisMonth.Unix()).
		Select("COALESCE(SUM(total_amount), 0) as total").Scan(&prevRevResult)

	var thisRevResult struct{ Total float64 }
	r.db.Table("booking_models").
		Where("status IN ? AND created_at >= ?", []string{"paid", "completed"}, firstDayThisMonth.Unix()).
		Select("COALESCE(SUM(total_amount), 0) as total").Scan(&thisRevResult)

	revPct := 0.0
	if prevRevResult.Total > 0 {
		revPct = (thisRevResult.Total - prevRevResult.Total) / prevRevResult.Total * 100
	}
	stats.Revenue = domain.AdminStatCard{Value: totalRevResult.Total, Percentage: revPct}

	
	var totalUsers, prevMonthUsers, thisMonthUsers int64
	r.db.Table("user_models").Count(&totalUsers)
	r.db.Table("user_models").Where("created_at >= ? AND created_at < ?", firstDayPrevMonth.Unix(), firstDayThisMonth.Unix()).Count(&prevMonthUsers)
	r.db.Table("user_models").Where("created_at >= ?", firstDayThisMonth.Unix()).Count(&thisMonthUsers)

	usersPct := 0.0
	if prevMonthUsers > 0 {
		usersPct = float64(thisMonthUsers-prevMonthUsers) / float64(prevMonthUsers) * 100
	}
	stats.TotalUsers = domain.AdminStatCard{Value: float64(totalUsers), Percentage: usersPct}

	
	var totalOrgs, prevOrgs, thisOrgs int64
	r.db.Table("user_role_models").Where("role = ?", "organizer").Count(&totalOrgs)
	r.db.Table("user_role_models").
		Joins("JOIN user_models ON user_models.id::uuid = user_role_models.user_id::uuid").
		Where("user_role_models.role = ? AND user_models.created_at >= ? AND user_models.created_at < ?", "organizer", firstDayPrevMonth.Unix(), firstDayThisMonth.Unix()).
		Count(&prevOrgs)
	r.db.Table("user_role_models").
		Joins("JOIN user_models ON user_models.id::uuid = user_role_models.user_id::uuid").
		Where("user_role_models.role = ? AND user_models.created_at >= ?", "organizer", firstDayThisMonth.Unix()).
		Count(&thisOrgs)

	orgsPct := 0.0
	if prevOrgs > 0 {
		orgsPct = float64(thisOrgs-prevOrgs) / float64(prevOrgs) * 100
	}
	stats.TotalOrganizers = domain.AdminStatCard{Value: float64(totalOrgs), Percentage: orgsPct}

	
	var totalEvents, prevEvents, thisMonthEvents int64
	r.db.Table("event_models").Count(&totalEvents)
	r.db.Table("event_models").Where("created_at >= ? AND created_at < ?", firstDayPrevMonth.Unix(), firstDayThisMonth.Unix()).Count(&prevEvents)
	r.db.Table("event_models").Where("created_at >= ?", firstDayThisMonth.Unix()).Count(&thisMonthEvents)

	eventsPct := 0.0
	if prevEvents > 0 {
		eventsPct = float64(thisMonthEvents-prevEvents) / float64(prevEvents) * 100
	}
	stats.TotalEvents = domain.AdminStatCard{Value: float64(totalEvents), Percentage: eventsPct}

	
	var totalBookings, refundedBookings int64
	r.db.Table("booking_models").Count(&totalBookings)
	r.db.Table("booking_models").Where("status = ?", "refunded").Count(&refundedBookings)
	refundRate := 0.0
	if totalBookings > 0 {
		refundRate = float64(refundedBookings) / float64(totalBookings) * 100
	}
	
	var prevTotal, prevRefunded int64
	r.db.Table("booking_models").Where("created_at >= ? AND created_at < ?", firstDayPrevMonth.Unix(), firstDayThisMonth.Unix()).Count(&prevTotal)
	r.db.Table("booking_models").Where("status = ? AND created_at >= ? AND created_at < ?", "refunded", firstDayPrevMonth.Unix(), firstDayThisMonth.Unix()).Count(&prevRefunded)
	prevRefundRate := 0.0
	if prevTotal > 0 {
		prevRefundRate = float64(prevRefunded) / float64(prevTotal) * 100
	}
	refundRatePct := 0.0
	if prevRefundRate > 0 {
		refundRatePct = (refundRate - prevRefundRate) / prevRefundRate * 100
	}
	stats.RefundRate = domain.AdminStatCard{Value: refundRate, Percentage: refundRatePct}

	
	var prevGrowth int64
	r.db.Table("user_models").Where("created_at >= ? AND created_at < ?", firstDayPrevMonth.AddDate(0, -1, 0).Unix(), firstDayPrevMonth.Unix()).Count(&prevGrowth)
	growthPct := 0.0
	if prevGrowth > 0 {
		growthPct = float64(thisMonthUsers-prevGrowth) / float64(prevGrowth) * 100
	}
	stats.UserGrowth = domain.AdminStatCard{Value: float64(thisMonthUsers), Percentage: growthPct}

	
	var pendingEvents int64
	r.db.Table("event_models").Where("status = ?", "pending").Count(&pendingEvents)
	stats.PendingApprovals = domain.AdminStatCard{Value: float64(pendingEvents)}

	
	var activeEvents, prevActiveEvents int64
	r.db.Table("event_models").Where("status = ?", "live").Count(&activeEvents)
	r.db.Table("event_models").Where("status = ? AND created_at >= ? AND created_at < ?", "live", firstDayPrevMonth.Unix(), firstDayThisMonth.Unix()).Count(&prevActiveEvents)
	activePct := 0.0
	if prevActiveEvents > 0 {
		activePct = float64(activeEvents-prevActiveEvents) / float64(prevActiveEvents) * 100
	}
	stats.ActiveEvents = domain.AdminStatCard{Value: float64(activeEvents), Percentage: activePct}

	
	var bookingRows []struct {
		TotalAmount float64
		CreatedAt   int64 `gorm:"column:created_at"`
	}
	r.db.Table("booking_models").
		Select("total_amount, created_at").
		Where("status IN ? AND created_at >= ?", []string{"paid", "completed"}, now.AddDate(-1, 0, 0).Unix()).
		Find(&bookingRows)

	revenueByMonth := make(map[string]float64)
	for i := 11; i >= 0; i-- {
		m := now.AddDate(0, -i, 0).Format("Jan")
		revenueByMonth[m] = 0
	}
	for _, b := range bookingRows {
		key := time.Unix(b.CreatedAt, 0).Format("Jan")
		if _, ok := revenueByMonth[key]; ok {
			revenueByMonth[key] += b.TotalAmount
		}
	}
	for i := 11; i >= 0; i-- {
		m := now.AddDate(0, -i, 0).Format("Jan")
		stats.RevenueOverview = append(stats.RevenueOverview, domain.RevenuePoint{
			Date:   m,
			Amount: revenueByMonth[m],
		})
	}

	return &stats, nil
}

func (r *eventGormRepository) GetAdminRevenueReport(startDate, endDate time.Time) (*domain.AdminRevenueReport, error) {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	firstDayThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	firstDayPrevMonth := firstDayThisMonth.AddDate(0, -1, 0)

	if startDate.IsZero() {
		startDate = now.AddDate(-1, 0, 0)
	}
	if endDate.IsZero() {
		endDate = now
	}

	var report domain.AdminRevenueReport

	
	var todayRev struct{ Total float64 }
	r.db.Table("booking_models").
		Where("status IN ? AND created_at >= ?", []string{"paid", "completed"}, todayStart.Unix()).
		Select("COALESCE(SUM(total_amount), 0) as total").Scan(&todayRev)

	var yesterdayRev struct{ Total float64 }
	r.db.Table("booking_models").
		Where("status IN ? AND created_at >= ? AND created_at < ?", []string{"paid", "completed"},
			todayStart.AddDate(0, 0, -1).Unix(), todayStart.Unix()).
		Select("COALESCE(SUM(total_amount), 0) as total").Scan(&yesterdayRev)

	todayPct := 0.0
	if yesterdayRev.Total > 0 {
		todayPct = (todayRev.Total - yesterdayRev.Total) / yesterdayRev.Total * 100
	}
	report.RevenueToday = domain.AdminStatCard{Value: todayRev.Total, Percentage: todayPct}

	
	var thisMonthRev struct{ Total float64 }
	r.db.Table("booking_models").
		Where("status IN ? AND created_at >= ?", []string{"paid", "completed"}, firstDayThisMonth.Unix()).
		Select("COALESCE(SUM(total_amount), 0) as total").Scan(&thisMonthRev)

	var prevMonthRev struct{ Total float64 }
	r.db.Table("booking_models").
		Where("status IN ? AND created_at >= ? AND created_at < ?", []string{"paid", "completed"},
			firstDayPrevMonth.Unix(), firstDayThisMonth.Unix()).
		Select("COALESCE(SUM(total_amount), 0) as total").Scan(&prevMonthRev)

	monthPct := 0.0
	if prevMonthRev.Total > 0 {
		monthPct = (thisMonthRev.Total - prevMonthRev.Total) / prevMonthRev.Total * 100
	}
	report.RevenueThisMonth = domain.AdminStatCard{Value: thisMonthRev.Total, Percentage: monthPct}

	
	var totalRev struct{ Total float64 }
	r.db.Table("booking_models").Where("status IN ?", []string{"paid", "completed"}).
		Select("COALESCE(SUM(total_amount), 0) as total").Scan(&totalRev)

	
	var prevYearRev struct{ Total float64 }
	r.db.Table("booking_models").
		Where("status IN ? AND created_at >= ? AND created_at < ?", []string{"paid", "completed"},
			now.AddDate(-2, 0, 0).Unix(), now.AddDate(-1, 0, 0).Unix()).
		Select("COALESCE(SUM(total_amount), 0) as total").Scan(&prevYearRev)

	var thisYearRev struct{ Total float64 }
	r.db.Table("booking_models").
		Where("status IN ? AND created_at >= ?", []string{"paid", "completed"}, now.AddDate(-1, 0, 0).Unix()).
		Select("COALESCE(SUM(total_amount), 0) as total").Scan(&thisYearRev)

	totalPct := 0.0
	if prevYearRev.Total > 0 {
		totalPct = (thisYearRev.Total - prevYearRev.Total) / prevYearRev.Total * 100
	}
	report.TotalRevenue = domain.AdminStatCard{Value: totalRev.Total, Percentage: totalPct}

	
	growthRate := monthPct
	report.GrowthRate = domain.AdminStatCard{Value: growthRate, Percentage: monthPct}

	
	var bookingRows []struct {
		TotalAmount float64
		CreatedAt   int64 `gorm:"column:created_at"`
	}
	r.db.Table("booking_models").
		Select("total_amount, created_at").
		Where("status IN ? AND created_at >= ? AND created_at <= ?", []string{"paid", "completed"}, startDate.Unix(), endDate.Unix()).
		Find(&bookingRows)

	revenueByMonth := make(map[string]float64)
	months := int(endDate.Sub(startDate).Hours() / (24 * 30))
	if months < 1 {
		months = 1
	}
	for i := months; i >= 0; i-- {
		m := startDate.AddDate(0, i, 0)
		if m.After(endDate) {
			continue
		}
		key := startDate.AddDate(0, months-i, 0).Format("Jan")
		revenueByMonth[key] = 0
	}
	
	revenueByMonth = make(map[string]float64)
	cur := time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, startDate.Location())
	for !cur.After(endDate) {
		revenueByMonth[cur.Format("Jan")] = 0
		cur = cur.AddDate(0, 1, 0)
	}

	for _, b := range bookingRows {
		key := time.Unix(b.CreatedAt, 0).Format("Jan")
		if _, ok := revenueByMonth[key]; ok {
			revenueByMonth[key] += b.TotalAmount
		}
	}

	cur = time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, startDate.Location())
	for !cur.After(endDate) {
		key := cur.Format("Jan")
		report.RevenueOverTime = append(report.RevenueOverTime, domain.RevenuePoint{
			Date:   key,
			Amount: revenueByMonth[key],
		})
		cur = cur.AddDate(0, 1, 0)
	}

	
	var catRows []struct {
		Category string
		Revenue  float64
	}
	r.db.Table("booking_models").
		Joins("JOIN event_models ON event_models.id::uuid = booking_models.event_id::uuid").
		Select("event_models.category as category, COALESCE(SUM(booking_models.total_amount), 0) as revenue").
		Where("booking_models.status IN ?", []string{"paid", "completed"}).
		Group("event_models.category").
		Order("revenue DESC").
		Scan(&catRows)

	for _, row := range catRows {
		cat := row.Category
		if cat == "" {
			cat = "Others"
		}
		report.CategoryBreakdown = append(report.CategoryBreakdown, domain.CategoryRevenue{
			Category: cat,
			Revenue:  row.Revenue,
		})
	}

	
	var refundRows []struct {
		TotalAmount float64
		CreatedAt   int64 `gorm:"column:created_at"`
	}
	r.db.Table("booking_models").
		Select("total_amount, created_at").
		Where("status = ? AND created_at >= ?", "refunded", now.AddDate(0, -6, 0).Unix()).
		Find(&refundRows)

	refundByMonth := make(map[string]float64)
	for i := 5; i >= 0; i-- {
		m := now.AddDate(0, -i, 0).Format("Jan")
		refundByMonth[m] = 0
	}
	for _, r2 := range refundRows {
		key := time.Unix(r2.CreatedAt, 0).Format("Jan")
		if _, ok := refundByMonth[key]; ok {
			refundByMonth[key] += r2.TotalAmount
		}
	}

	
	var prevRefundTotal float64
	for k, v := range refundByMonth {
		if k == firstDayPrevMonth.Format("Jan") {
			prevRefundTotal = v
		}
	}
	thisMonthRefund := refundByMonth[now.Format("Jan")]
	refundPct := 0.0
	if prevRefundTotal > 0 {
		refundPct = (thisMonthRefund - prevRefundTotal) / prevRefundTotal * 100
	}

	var totalRefundAmount float64
	for _, v := range refundByMonth {
		totalRefundAmount += v
	}
	report.RefundTotal = domain.AdminStatCard{Value: totalRefundAmount, Percentage: refundPct}

	for i := 5; i >= 0; i-- {
		m := now.AddDate(0, -i, 0)
		report.RefundAnalytics = append(report.RefundAnalytics, domain.RefundDataPoint{
			Month:  m.Format("Jan"),
			Amount: refundByMonth[m.Format("Jan")],
		})
	}

	
	var orgRows []struct {
		OrganizerID   string
		Name          string
		TotalRevenue  float64
		ActiveEvents  int
		PendingEvents int
	}
	r.db.Table("booking_models").
		Joins("JOIN event_models ON event_models.id::uuid = booking_models.event_id::uuid").
		Joins("JOIN user_models ON user_models.id::uuid = event_models.organizer_id::uuid").
		Select(`event_models.organizer_id,
			user_models.name,
			COALESCE(SUM(booking_models.total_amount), 0) as total_revenue,
			COUNT(DISTINCT CASE WHEN event_models.status = 'live' THEN event_models.id END) as active_events,
			COUNT(DISTINCT CASE WHEN event_models.status = 'pending' THEN event_models.id END) as pending_events`).
		Where("booking_models.status IN ?", []string{"paid", "completed"}).
		Group("event_models.organizer_id, user_models.name").
		Order("total_revenue DESC").
		Limit(5).
		Scan(&orgRows)

	for _, o := range orgRows {
		var avgRating struct{ Avg float64 }
		r.db.Table("event_details_models").
			Joins("JOIN event_models ON event_models.id::uuid = event_details_models.event_id::uuid").
			Where("event_models.organizer_id = ?", o.OrganizerID).
			Select("COALESCE(AVG(event_details_models.rating), 0) as avg").
			Scan(&avgRating)

		report.TopOrganizers = append(report.TopOrganizers, domain.TopOrganizerEntry{
			Name:           o.Name,
			Revenue:        o.TotalRevenue,
			ActiveEvents:   o.ActiveEvents,
			PendingEvents:  o.PendingEvents,
			AvgEventRating: math.Round(avgRating.Avg*10) / 10,
		})
	}

	
	var userRows []struct {
		UserID        string
		Name          string
		TotalSpent    float64
		BookingsCount int
	}
	r.db.Table("booking_models").
		Joins("JOIN user_models ON user_models.id::uuid = booking_models.user_id::uuid").
		Select(`booking_models.user_id,
			user_models.name,
			COALESCE(SUM(booking_models.total_amount), 0) as total_spent,
			COUNT(*) as bookings_count`).
		Where("booking_models.status IN ?", []string{"paid", "completed"}).
		Group("booking_models.user_id, user_models.name").
		Order("total_spent DESC").
		Limit(5).
		Scan(&userRows)

	for _, u := range userRows {
		report.TopUsers = append(report.TopUsers, domain.TopUserEntry{
			Name:           u.Name,
			EventsAttended: u.BookingsCount,
			TotalSpent:     u.TotalSpent,
		})
	}

	return &report, nil
}
