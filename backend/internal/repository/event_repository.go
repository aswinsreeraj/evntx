package repository

import (
	"context"

	"github.com/aswinsreeraj/evntx/internal/domain"
)

type EventRepository interface {
	ListLiveEvents(city string, page int, limit int) ([]domain.Event, int64, error)
	GetEventBySlug(slug string) (*domain.Event, error)
	GetEventDetails(eventID string) (*domain.EventDetails, error)
	GetEventPersonnels(eventID string) ([]domain.EventPersonnel, error)
	GetEventByID(eventID string) (*domain.Event, error)
	CreateEvent(ctx context.Context, event *domain.Event, details *domain.EventDetails, tickets []domain.TicketType) error
	UpdateEvent(ctx context.Context, eventID string, eventUpdates map[string]interface{}, detailUpdates map[string]interface{}, ticketUpdates []domain.TicketType) error
}
