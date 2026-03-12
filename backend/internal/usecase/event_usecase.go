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

func (u *EventUsecase) UpdateEvent(
	ctx context.Context,
	organizerID string,
	eventID string,
	eventUpdates map[string]interface{},
	detailsUpdates map[string]interface{},
	ticketUpdates []domain.TicketType,
) error {

	event, err := u.repo.GetEventByID(eventID)
	if err != nil {
		return errors.New("event not found")
	}

	if event.OrganizerID != organizerID {
		return errors.New("EVT_004: Forbidden action")
	}

	if event.Status != "draft" && event.Status != "rejected" {
		return errors.New("EVT_006: Event cannot be updated in current state")
	}

	details, err := u.repo.GetEventDetails(eventID)
	if err != nil {
		return errors.New("event details not found")
	}

	// Calculate ticket bounds correctly using old vs new capacity maps
	capacity := details.TotalCapacity
	if newCap, exists := detailsUpdates["total_capacity"]; exists {
		if c, ok := newCap.(int); ok {
			capacity = c
		} else if c, ok := newCap.(float64); ok {
			capacity = int(c)
		}
	}

	totalTickets := 0
	for _, t := range ticketUpdates {
		if t.Price < 0 {
			return errors.New("ticket price must be >= 0")
		}
		if t.TotalQuantity < 0 {
			return errors.New("ticket quantity must be >= 0")
		}
		totalTickets += t.TotalQuantity
	}

	// Need to check if total tickets exceed capacity, but this rule might be tricky if they don't submit all tickets.
	// Since we assume ticketUpdates replaces or upserts, let's just make sure the sum of incoming ones doesn't overflow.
	// To be perfectly safe, we'll enforce that the total updated sum doesn't exceed the capacity at face value.
	if totalTickets > capacity {
		return errors.New("sum of ticket quantities must not exceed total capacity")
	}

	now := time.Now()
	eventUpdates["updated_at"] = now
	detailsUpdates["updated_at"] = now

	for i := range ticketUpdates {
		if ticketUpdates[i].ID == "" {
			ticketUpdates[i].ID = uuid.NewString()
			ticketUpdates[i].CreatedAt = now
		}
		ticketUpdates[i].EventID = eventID
		ticketUpdates[i].UpdatedAt = now
		
		// If AvailableQuantity somehow wasn't passed, sync it to total. Note: Partial updates from frontend might be tricky for tickets
		// The requirement demands we upsert rows. We sync available natively.
		ticketUpdates[i].AvailableQuantity = ticketUpdates[i].TotalQuantity 
	}

	err = u.repo.UpdateEvent(ctx, eventID, eventUpdates, detailsUpdates, ticketUpdates)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to update event in database")
		return err
	}

	logger.Log.Info().
		Str("event_id", eventID).
		Str("organizer_id", organizerID).
		Time("timestamp", now).
		Msg("event_updated")

	return nil
}

func (u *EventUsecase) SubmitEventForApproval(ctx context.Context, organizerID string, eventID string) error {
	event, err := u.repo.GetEventByID(eventID)
	if err != nil {
		return errors.New("event not found")
	}

	if event.OrganizerID != organizerID {
		return errors.New("EVT_004: Forbidden action")
	}

	if event.Status != "draft" {
		return errors.New("EVT_006: Event cannot be submitted in current state")
	}

	err = u.repo.UpdateEventStatus(ctx, eventID, "pending")
	if err != nil {
		return err
	}

	logger.Log.Info().
		Str("event_id", eventID).
		Str("organizer_id", organizerID).
		Str("previous_status", event.Status).
		Str("new_status", "pending").
		Time("timestamp", time.Now()).
		Msg("event_submitted_for_approval")

	return nil
}

