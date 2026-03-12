package repository

import (
	"context"
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventModel struct {
	ID            string
	OrganizerID   string
	Title         string
	Slug          string
	City          string
	VenueName     string
	Category      string
	StartTime     int64
	EndTime       int64
	Tags          string
	Status        string
	CoverImageURL string
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

func (r *eventGormRepository) ListLiveEvents(city string, page int, limit int) ([]domain.Event, int64, error) {

	var models []EventModel
	var total int64

	query := r.db.Model(&EventModel{}).Where("status = ?", "live")

	if city != "" {
		query = query.Where("city ILIKE ?", "%"+city+"%")
	}

	query.Count(&total)

	offset := (page - 1) * limit

	err := query.
		Order("start_time ASC").
		Limit(limit).
		Offset(offset).
		Find(&models).Error

	if err != nil {
		return nil, 0, err
	}

	events := make([]domain.Event, 0)

	for _, m := range models {
		events = append(events, domain.Event{
			ID:            m.ID,
			Title:         m.Title,
			Slug:          m.Slug,
			City:          m.City,
			VenueName:     m.VenueName,
			Category:      m.Category,
			CoverImageURL: m.CoverImageURL,
		})
	}

	return events, total, nil
}

func (r *eventGormRepository) GetEventBySlug(slug string) (*domain.Event, error) {

	var model EventModel

	err := r.db.
		Where("slug = ? AND status = ?", slug, "live").
		First(&model).Error

	if err != nil {
		return nil, err
	}

	return &domain.Event{
		ID:            model.ID,
		OrganizerID:   model.OrganizerID,
		Title:         model.Title,
		Slug:          model.Slug,
		Status:        model.Status,
		City:          model.City,
		VenueName:     model.VenueName,
		Category:      model.Category,
		CoverImageURL: model.CoverImageURL,
	}, nil
}

func (r *eventGormRepository) GetEventByID(eventID string) (*domain.Event, error) {

	var model EventModel

	err := r.db.
		Where("id = ?", eventID).
		First(&model).Error

	if err != nil {
		return nil, err
	}

	return &domain.Event{
		ID:            model.ID,
		OrganizerID:   model.OrganizerID,
		Title:         model.Title,
		Slug:          model.Slug,
		Status:        model.Status,
		City:          model.City,
		VenueName:     model.VenueName,
		Category:      model.Category,
		CoverImageURL: model.CoverImageURL,
	}, nil
}

func (r *eventGormRepository) GetEventDetails(eventID string) (*domain.EventDetails, error) {

	var model EventDetailsModel

	err := r.db.
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

func (r *eventGormRepository) CreateEvent(ctx context.Context, event *domain.Event, details *domain.EventDetails, tickets []domain.TicketType) error {
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
			CoverImageURL: event.CoverImageURL,
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

		return nil
	})
}

func (r *eventGormRepository) UpdateEvent(ctx context.Context, eventID string, eventUpdates map[string]interface{}, detailUpdates map[string]interface{}, ticketUpdates []domain.TicketType) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

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

			// Save handles both Insert and Update depending on whether Primary Key (ID) explicitly exists
			if err := tx.Save(&model).Error; err != nil {
				return err
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
