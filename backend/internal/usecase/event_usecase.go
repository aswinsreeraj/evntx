package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/internal/repository"
	apiErrors "github.com/aswinsreeraj/evntx/pkg/errors"
	"github.com/aswinsreeraj/evntx/pkg/logger"
	"github.com/google/uuid"
)

type EventUsecase struct {
	repo			repository.EventRepository
	bookingRepo		repository.BookingRepository
	notificationUsecase	*NotificationUsecase
	settingsRepo		repository.SettingsRepository
}

func NewEventUsecase(
	repo repository.EventRepository,
	bookingRepo repository.BookingRepository,
	notificationUsecase *NotificationUsecase,
	settingsRepo repository.SettingsRepository,
) *EventUsecase {
	return &EventUsecase{repo: repo, bookingRepo: bookingRepo, notificationUsecase: notificationUsecase, settingsRepo: settingsRepo}
}

func (u *EventUsecase) ListEvents(city, category, search, sortBy, minPrice, maxPrice, startDate, endDate string, page int, limit int) (interface{}, int64, float64, float64, error) {
	return u.repo.ListLiveEvents(city, category, search, sortBy, minPrice, maxPrice, startDate, endDate, page, limit)
}

func (u *EventUsecase) AdminSearchEvents(search string, status string, page int, limit int) ([]domain.AdminEventDetails, int64, error) {
	return u.repo.AdminSearchEvents(search, status, page, limit)
}

func (u *EventUsecase) GetOrganizerDashboardStats(organizerID string) (*domain.OrganizerDashboardStats, error) {
	return u.repo.GetDashboardStats(organizerID)
}

func (u *EventUsecase) GetSalesReport(organizerID string, eventID string, startDate string, endDate string) (*domain.SalesReportStats, error) {
	return u.repo.GetSalesReport(organizerID, eventID, startDate, endDate)
}

func (u *EventUsecase) GetAdminDashboardStats(span string, groupBy string) (*domain.AdminDashboardStats, error) {
	return u.repo.GetAdminDashboardStats(span, groupBy)
}

func (u *EventUsecase) GetAdminRevenueReport(startDate, endDate time.Time, groupBy string) (*domain.AdminRevenueReport, error) {
	return u.repo.GetAdminRevenueReport(startDate, endDate, groupBy)
}

func (u *EventUsecase) GetEventByID(id string) (*domain.Event, error) {
	return u.repo.GetEventByID(id)
}

func (u *EventUsecase) GetEvent(slug string) (interface{}, interface{}, interface{}, interface{}, error) {

	event, err := u.repo.GetEventBySlug(slug)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	if event.Status != "approved" && event.Status != "live" && event.Status != "completed" {
		return nil, nil, nil, nil, apiErrors.ErrResourceNotFound
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
		return nil, nil, nil, nil, apiErrors.ErrResourceNotFound
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
		return apiErrors.ErrResourceNotFound
	}

	if event.OrganizerID != organizerID {
		return apiErrors.ErrForbiddenAction
	}

	if event.Status != "draft" && event.Status != "rejected" && event.Status != "approved" {
		return apiErrors.ErrInvalidStateTransition
	}

	details, err := u.repo.GetEventDetails(eventID)
	if err != nil {
		return apiErrors.ErrResourceNotFound
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

	return nil
}

func (u *EventUsecase) SubmitEventForApproval(ctx context.Context, organizerID string, eventID string) error {
	event, err := u.repo.GetEventByID(eventID)
	if err != nil {
		return apiErrors.ErrResourceNotFound
	}

	if event.OrganizerID != organizerID {
		return apiErrors.ErrForbiddenAction
	}

	if event.Status != "draft" && event.Status != "rejected" {
		return apiErrors.ErrInvalidStateTransition
	}

	nextStatus := "pending"
	if u.settingsRepo != nil {
		if settings, settingsErr := u.settingsRepo.GetPlatformSettings(); settingsErr == nil && !settings.RequireAdminApprovalForEvents {
			nextStatus = "approved"
		}
	}

	err = u.repo.UpdateEventStatus(ctx, eventID, nextStatus)
	if err != nil {
		return err
	}

	logger.Log.Info().
		Str("event", "event_state_changed").
		Str("entity", "event").
		Str("entity_id", eventID).
		Str("from", event.Status).
		Str("to", nextStatus).
		Str("actor_id", organizerID).
		Msg("")

	return nil
}

func (u *EventUsecase) ApproveEvent(ctx context.Context, adminID string, eventID string) error {
	event, err := u.repo.GetEventByID(eventID)
	if err != nil {
		return apiErrors.ErrResourceNotFound
	}

	if event.Status != "pending" {
		return apiErrors.ErrInvalidStateTransition
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
		return apiErrors.ErrResourceNotFound
	}

	if event.Status != "pending" {
		return apiErrors.ErrInvalidStateTransition
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

func (u *EventUsecase) SuspendLiveEvent(ctx context.Context, adminID string, eventID string, reason string) error {
	event, err := u.repo.GetEventByID(eventID)
	if err != nil {
		return apiErrors.ErrResourceNotFound
	}

	if event.Status != "live" {
		return apiErrors.ErrInvalidStateTransition
	}

	err = u.repo.SuspendLiveEvent(ctx, eventID, adminID, reason)
	if err != nil {
		return err
	}

	if u.notificationUsecase != nil {
		err = u.notificationUsecase.SendNotification(
			event.OrganizerID,
			"event_suspended",
			"Event Suspended",
			"Your event '"+event.Title+"' has been suspended by the admin. Reason: "+reason,
			nil,
		)
		if err != nil {
			logger.Log.Error().Err(err).Msg("Failed to send suspension notification")
		}
	}

	logger.Log.Info().
		Str("event", "event_state_changed").
		Str("entity", "event").
		Str("entity_id", eventID).
		Str("from", event.Status).
		Str("to", "suspended").
		Str("actor_id", adminID).
		Msg("")

	return nil
}

func (u *EventUsecase) CompleteEvent(ctx context.Context, adminID string, eventID string) error {
	event, err := u.repo.GetEventByID(eventID)
	if err != nil {
		return apiErrors.ErrResourceNotFound
	}

	if event.Status != "live" {
		return apiErrors.ErrInvalidStateTransition
	}

	if err := u.repo.UpdateEventStatus(ctx, eventID, "completed"); err != nil {
		return err
	}

	if err := u.SettleEventEarnings(ctx, eventID); err != nil {
		return err
	}

	logger.Log.Info().
		Str("event", "event_state_changed").
		Str("entity", "event").
		Str("entity_id", eventID).
		Str("from", event.Status).
		Str("to", "completed").
		Str("actor_id", adminID).
		Msg("")

	return nil
}

func (u *EventUsecase) SettleEventEarnings(ctx context.Context, eventID string) error {
	event, err := u.repo.GetEventByID(eventID)
	if err != nil {
		return apiErrors.ErrResourceNotFound
	}

	if event.Status != "completed" {
		return apiErrors.ErrInvalidStateTransition
	}

	if event.Settled {
		return apiErrors.ErrDuplicateResource
	}

	bookings, err := u.bookingRepo.GetPaidBookingsByEventID(ctx, eventID)
	if err != nil {
		return err
	}

	var totalAmount float64
	for _, booking := range bookings {
		totalAmount += booking.TotalAmount
	}

	if err := u.repo.SettleEventEarnings(ctx, eventID, event.OrganizerID, totalAmount); err != nil {
		return err
	}

	logger.Log.Info().
		Str("event", "event_earnings_settled").
		Str("entity", "event").
		Str("entity_id", eventID).
		Str("organizer_id", event.OrganizerID).
		Float64("amount", totalAmount).
		Msg("")

	return nil
}

func (u *EventUsecase) GetOrganizerEvents(ctx context.Context, organizerID string, status string) ([]domain.Event, error) {
	return u.repo.GetEventsByOrganizerID(organizerID, status)
}

func (u *EventUsecase) DeleteEvent(ctx context.Context, organizerID string, eventID string) error {
	event, err := u.repo.GetEventByID(eventID)
	if err != nil {
		return apiErrors.ErrResourceNotFound
	}

	if event.OrganizerID != organizerID {
		return apiErrors.ErrForbiddenAction
	}

	if event.Status == "completed" || event.Status == "live" {
		return apiErrors.New(400, apiErrors.InvalidStateTransition, "Cannot delete a completed or live event")
	}

	err = u.repo.DeleteEvent(ctx, eventID)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to delete event")
		return err
	}

	return nil
}

func (u *EventUsecase) RequestEventCancellation(ctx context.Context, organizerID string, eventID string, reason string) error {
	event, err := u.repo.GetEventByID(eventID)
	if err != nil {
		return apiErrors.ErrResourceNotFound
	}
	if event.OrganizerID != organizerID {
		return apiErrors.ErrForbiddenAction
	}
	if event.Status != "live" {
		return apiErrors.ErrInvalidStateTransition
	}
	if u.settingsRepo != nil {
		if settings, settingsErr := u.settingsRepo.GetPlatformSettings(); settingsErr == nil && !settings.AllowEventCancellation {
			return apiErrors.New(403, apiErrors.ForbiddenAction, "Event cancellation requests are currently disabled by admin")
		}
	}
	if err := u.repo.RequestEventCancellation(ctx, eventID, organizerID, reason); err != nil {
		return err
	}
	logger.Log.Info().
		Str("event", "event_state_changed").
		Str("entity", "event").
		Str("entity_id", eventID).
		Str("from", event.Status).
		Str("to", "cancellation_pending").
		Str("actor_id", organizerID).
		Msg("")
	return nil
}

func (u *EventUsecase) ApproveEventCancellation(ctx context.Context, adminID string, eventID string) error {
	event, err := u.repo.GetEventByID(eventID)
	if err != nil {
		return apiErrors.ErrResourceNotFound
	}
	if event.Status != "cancellation_pending" {
		return apiErrors.ErrInvalidStateTransition
	}
	if err := u.repo.CancelLiveEvent(ctx, eventID, event.OrganizerID); err != nil {
		return err
	}
	if u.notificationUsecase != nil {
		_ = u.notificationUsecase.SendNotification(
			event.OrganizerID,
			"event_cancellation_approved",
			"Event cancellation approved",
			"Your event cancellation request has been approved by admin.",
			map[string]interface{}{"event_id": eventID},
		)
	}
	return nil
}

func (u *EventUsecase) RejectEventCancellation(ctx context.Context, adminID string, eventID string, reason string) error {
	event, err := u.repo.GetEventByID(eventID)
	if err != nil {
		return apiErrors.ErrResourceNotFound
	}
	if event.Status != "cancellation_pending" {
		return apiErrors.ErrInvalidStateTransition
	}
	if err := u.repo.RejectEventCancellation(ctx, eventID, adminID, reason); err != nil {
		return err
	}
	if u.notificationUsecase != nil {
		_ = u.notificationUsecase.SendNotification(
			event.OrganizerID,
			"event_cancellation_rejected",
			"Event cancellation rejected",
			"Your event cancellation request was rejected by admin. Reason: "+reason,
			map[string]interface{}{"event_id": eventID, "reason": reason},
		)
	}
	return nil
}

func (u *EventUsecase) AutoProcessCompletedEvents(ctx context.Context) error {
	now := time.Now()
	events, err := u.repo.FindPastLiveEvents(ctx, now)
	if err != nil {
		return err
	}

	for _, event := range events {
		if err := u.repo.UpdateEventStatus(ctx, event.ID, "completed"); err != nil {
			logger.Log.Error().Err(err).Str("event_id", event.ID).Msg("Failed to auto-complete event status")
			continue
		}

		if err := u.SettleEventEarnings(ctx, event.ID); err != nil {
			logger.Log.Error().Err(err).Str("event_id", event.ID).Msg("Failed to settle event earnings automatically")
			continue
		}

		logger.Log.Info().
			Str("event", "event_auto_completed").
			Str("entity_id", event.ID).
			Msg("Successfully auto-completed event and settled earnings")
	}

	return nil
}
