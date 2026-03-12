package repository

import (
	"github.com/aswinsreeraj/evntx/internal/domain"
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
		Title:         model.Title,
		Slug:          model.Slug,
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
