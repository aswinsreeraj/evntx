package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/internal/repository"
	"github.com/aswinsreeraj/evntx/pkg/logger"
	"github.com/google/uuid"
)

type EventUsecase struct {
	repo repository.EventRepository
}

func NewEventUsecase(repo repository.EventRepository) *EventUsecase {
	return &EventUsecase{repo: repo}
}

func (u *EventUsecase) ListEvents(city string, page int, limit int) (interface{}, int64, error) {
	return u.repo.ListLiveEvents(city, page, limit)
}

func (u *EventUsecase) GetEvent(slug string) (interface{}, interface{}, interface{}, error) {

	event, err := u.repo.GetEventBySlug(slug)
	if err != nil {
		return nil, nil, nil, err
	}

	details, err := u.repo.GetEventDetails(event.ID)
	if err != nil {
		return nil, nil, nil, err
	}

	personnels, err := u.repo.GetEventPersonnels(event.ID)
	if err != nil {
		return nil, nil, nil, err
	}

	return event, details, personnels, nil
}

func generateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "'", "")
	return slug
}

func (u *EventUsecase) CreateEvent(
	ctx context.Context,
	organizerID string,
	event *domain.Event,
	details *domain.EventDetails,
	tickets []domain.TicketType,
) (string, error) {

	if !event.StartTime.Before(event.EndTime) {
		return "", errors.New("start_time must be before end_time")
	}

	totalTickets := 0
	for _, t := range tickets {
		if t.Price < 0 {
			return "", errors.New("ticket price must be >= 0")
		}
		if t.TotalQuantity < 0 {
			return "", errors.New("ticket quantity must be >= 0")
		}
		totalTickets += t.TotalQuantity
	}

	if totalTickets > details.TotalCapacity {
		return "", errors.New("total ticket quantities must not exceed event capacity")
	}

	eventID := uuid.NewString()
	now := time.Now()

	event.ID = eventID
	event.OrganizerID = organizerID
	event.Slug = generateSlug(event.Title)
	event.Status = "draft"
	event.CreatedAt = now
	event.UpdatedAt = now

	details.EventID = eventID
	details.AvailableCapacity = details.TotalCapacity
	details.CreatedAt = now
	details.UpdatedAt = now

	for i := range tickets {
		tickets[i].ID = uuid.NewString()
		tickets[i].EventID = eventID
		tickets[i].AvailableQuantity = tickets[i].TotalQuantity
		tickets[i].Version = 1
		tickets[i].CreatedAt = now
		tickets[i].UpdatedAt = now
	}

	err := u.repo.CreateEvent(ctx, event, details, tickets)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to create event in database")
		return "", err
	}

	logger.Log.Info().
		Str("event_id", eventID).
		Str("organizer_id", organizerID).
		Time("timestamp", now).
		Msg("event_created")

	return eventID, nil
}

