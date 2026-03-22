package repository

import (
	"context"

	"github.com/aswinsreeraj/evntx/internal/domain"
)

type EventRepository interface {
	ListLiveEvents(city, category, search, sortBy, minPrice, maxPrice, startDate, endDate string, page int, limit int) ([]domain.Event, int64, float64, float64, error)
	AdminSearchEvents(search string, status string, page int, limit int) ([]domain.AdminEventDetails, int64, error)
	GetEventBySlug(slug string) (*domain.Event, error)
	GetEventDetails(eventID string) (*domain.EventDetails, error)
	GetEventPersonnels(eventID string) ([]domain.EventPersonnel, error)
	GetEventsByOrganizerID(organizerID string, status string) ([]domain.Event, error)
	GetTicketTypesByEventID(eventID string) ([]domain.TicketType, error)
	GetEventByID(eventID string) (*domain.Event, error)
	CreateEvent(ctx context.Context, event *domain.Event, details *domain.EventDetails, tickets []domain.TicketType, personnels []domain.EventPersonnel) error
	UpdateEvent(ctx context.Context, eventID string, eventUpdates map[string]interface{}, detailUpdates map[string]interface{}, ticketUpdates []domain.TicketType, personnelUpdates []domain.EventPersonnel) error
	UpdateEventStatus(ctx context.Context, eventID string, status string) error
	ApproveEvent(ctx context.Context, eventID string) error
	RejectEvent(ctx context.Context, eventID string, adminID string, reason string) error
	DeleteEvent(ctx context.Context, eventID string) error
}
