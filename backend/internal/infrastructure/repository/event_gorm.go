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
		Joins("JOIN event_models ON event_models.id = ticket_type_models.event_id").
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
			), 0) AS available_capacity
		`).
		Joins("LEFT JOIN user_models ON user_models.id::text = event_models.organizer_id")

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
			Joins("LEFT JOIN booking_ticket_models ON booking_ticket_models.booking_id = booking_models.id").
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
			"(SELECT reason FROM event_moderation_log_models WHERE event_id = event_models.id AND action IN ('rejected', 'suspended') ORDER BY created_at DESC LIMIT 1) as rejection_reason, "+
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
			Where("id = ? AND organizer_id = ? AND status IN ('live', 'approved')", eventID, organizerID).
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
