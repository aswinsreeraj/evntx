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

func (u *EventUsecase) ListEvents(city, category, search, sortBy, minPrice, maxPrice, startDate, endDate string, page int, limit int) (interface{}, int64, float64, float64, error) {
	return u.repo.ListLiveEvents(city, category, search, sortBy, minPrice, maxPrice, startDate, endDate, page, limit)
}

func (u *EventUsecase) AdminSearchEvents(search string, status string, page int, limit int) ([]domain.AdminEventDetails, int64, error) {
	return u.repo.AdminSearchEvents(search, status, page, limit)
}

func (u *EventUsecase) GetEvent(slug string) (interface{}, interface{}, interface{}, interface{}, error) {

	event, err := u.repo.GetEventBySlug(slug)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	if event.Status != "approved" && event.Status != "live" {
		return nil, nil, nil, nil, errors.New("event not found")
	}

	details, err := u.repo.GetEventDetails(event.ID)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	personnels, err := u.repo.GetEventPersonnels(event.ID)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	tickets, err := u.repo.GetTicketTypesByEventID(event.ID)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return event, details, personnels, tickets, nil
}

func (u *EventUsecase) GetOrganizerEvent(slug string, organizerID string) (interface{}, interface{}, interface{}, interface{}, error) {
	event, err := u.repo.GetEventBySlug(slug)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	if event.OrganizerID != organizerID {
		return nil, nil, nil, nil, errors.New("event not found")
	}

	details, err := u.repo.GetEventDetails(event.ID)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	personnels, err := u.repo.GetEventPersonnels(event.ID)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	tickets, err := u.repo.GetTicketTypesByEventID(event.ID)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return event, details, personnels, tickets, nil
}

func (u *EventUsecase) AdminGetEvent(slug string) (interface{}, interface{}, interface{}, interface{}, error) {
	event, err := u.repo.GetEventBySlug(slug)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	details, err := u.repo.GetEventDetails(event.ID)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	personnels, err := u.repo.GetEventPersonnels(event.ID)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	tickets, err := u.repo.GetTicketTypesByEventID(event.ID)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return event, details, personnels, tickets, nil
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
	personnels []domain.EventPersonnel,
) (string, error) {

	if !event.EndTime.IsZero() && !event.StartTime.Before(event.EndTime) {
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

	for i := range personnels {
		personnels[i].ID = uuid.NewString()
		personnels[i].EventID = eventID
	}

	err := u.repo.CreateEvent(ctx, event, details, tickets, personnels)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to create event in database")
		return "", err
	}

	logger.Log.Info().
		Str("event", "event_state_changed").
		Str("entity", "event").
		Str("entity_id", eventID).
		Str("from", "").
		Str("to", "draft").
		Str("actor_id", organizerID).
		Msg("")

	return eventID, nil
}

func (u *EventUsecase) UpdateEvent(
	ctx context.Context,
	organizerID string,
	eventID string,
	eventUpdates map[string]interface{},
	detailsUpdates map[string]interface{},
	ticketUpdates []domain.TicketType,
	personnelUpdates []domain.EventPersonnel,
) error {

	event, err := u.repo.GetEventByID(eventID)
	if err != nil {
		return errors.New("event not found")
	}

	if event.OrganizerID != organizerID {
		return errors.New("EVT_004: Forbidden action")
	}

	if event.Status != "draft" && event.Status != "rejected" && event.Status != "approved" {
		return errors.New("EVT_006: Event cannot be updated in current state")
	}

	details, err := u.repo.GetEventDetails(eventID)
	if err != nil {
		return errors.New("event details not found")
	}

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

	if totalTickets > capacity {
		return errors.New("sum of ticket quantities must not exceed total capacity")
	}

	now := time.Now()
	eventUpdates["updated_at"] = now
	detailsUpdates["updated_at"] = now
 
	if event.Status == "approved" {
		eventUpdates["status"] = "draft"
		logger.Log.Info().
			Str("event", "event_state_changed").
			Str("entity", "event").
			Str("entity_id", eventID).
			Str("from", "approved").
			Str("to", "draft").
			Str("actor_id", organizerID).
			Msg("")
	}

	for i := range ticketUpdates {
		if ticketUpdates[i].ID == "" {
			ticketUpdates[i].ID = uuid.NewString()
			ticketUpdates[i].CreatedAt = now
		}
		ticketUpdates[i].EventID = eventID
		ticketUpdates[i].UpdatedAt = now

		ticketUpdates[i].AvailableQuantity = ticketUpdates[i].TotalQuantity
	}

	for i := range personnelUpdates {
		if personnelUpdates[i].ID == "" {
			personnelUpdates[i].ID = uuid.NewString()
		}
		personnelUpdates[i].EventID = eventID
	}

	err = u.repo.UpdateEvent(ctx, eventID, eventUpdates, detailsUpdates, ticketUpdates, personnelUpdates)
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

	if event.Status != "draft" && event.Status != "rejected" {
		return errors.New("EVT_006: Event cannot be submitted in current state")
	}

	err = u.repo.UpdateEventStatus(ctx, eventID, "pending")
	if err != nil {
		return err
	}

	logger.Log.Info().
		Str("event", "event_state_changed").
		Str("entity", "event").
		Str("entity_id", eventID).
		Str("from", event.Status).
		Str("to", "pending").
		Str("actor_id", organizerID).
		Msg("")

	return nil
}

func (u *EventUsecase) ApproveEvent(ctx context.Context, adminID string, eventID string) error {
	event, err := u.repo.GetEventByID(eventID)
	if err != nil {
		return errors.New("event not found")
	}

	if event.Status == "approved" || event.Status == "live" || event.Status == "completed" {
		return errors.New("EVT_006: Event cannot be approved in current state")
	}

	err = u.repo.ApproveEvent(ctx, eventID)
	if err != nil {
		return err
	}

	logger.Log.Info().
		Str("event", "event_state_changed").
		Str("entity", "event").
		Str("entity_id", eventID).
		Str("from", event.Status).
		Str("to", "approved").
		Str("actor_id", adminID).
		Msg("")

	return nil
}

func (u *EventUsecase) RejectEvent(ctx context.Context, adminID string, eventID string, reason string) error {
	event, err := u.repo.GetEventByID(eventID)
	if err != nil {
		return errors.New("event not found")
	}

	if event.Status == "rejected" || event.Status == "completed" {
		return errors.New("EVT_006: Event cannot be rejected in current state")
	}

	err = u.repo.RejectEvent(ctx, eventID, adminID, reason)
	if err != nil {
		return err
	}

	logger.Log.Info().
		Str("event", "event_state_changed").
		Str("entity", "event").
		Str("entity_id", eventID).
		Str("from", event.Status).
		Str("to", "rejected").
		Str("actor_id", adminID).
		Msg("")

	return nil
}

func (u *EventUsecase) GetOrganizerEvents(ctx context.Context, organizerID string, status string) ([]domain.Event, error) {
	return u.repo.GetEventsByOrganizerID(organizerID, status)
}

func (u *EventUsecase) DeleteEvent(ctx context.Context, organizerID string, eventID string) error {
	event, err := u.repo.GetEventByID(eventID)
	if err != nil {
		return errors.New("event not found")
	}

	if event.OrganizerID != organizerID {
		return errors.New("EVT_004: Forbidden action")
	}

	err = u.repo.DeleteEvent(ctx, eventID)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to delete event")
		return err
	}

	logger.Log.Info().
		Str("event_id", eventID).
		Str("organizer_id", organizerID).
		Time("timestamp", time.Now()).
		Msg("event_deleted")

	return nil
}
