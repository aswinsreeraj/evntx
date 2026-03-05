package repository

import "github.com/aswinsreeraj/evntx/internal/domain"

type EventRepository interface {
	ListLiveEvents(city string, page int, limit int) ([]domain.Event, int64, error)
	GetEventBySlug(slug string) (*domain.Event, error)
	GetEventDetails(eventID string) (*domain.EventDetails, error)
	GetEventPersonnels(eventID string) ([]domain.EventPersonnel, error)
}
