package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/internal/repository"
	apiErrors "github.com/aswinsreeraj/evntx/pkg/errors"
	"github.com/aswinsreeraj/evntx/pkg/logger"
	"github.com/google/uuid"
)

type BookingUsecase struct {
	bookingRepo         repository.BookingRepository
	eventRepo           repository.EventRepository
	notificationUsecase *NotificationUsecase
	roleRepo            repository.UserRoleRepository
	settingsRepo        repository.SettingsRepository
}

func NewBookingUsecase(
	bookingRepo repository.BookingRepository,
	eventRepo repository.EventRepository,
	roleRepo repository.UserRoleRepository,
	notificationUsecase *NotificationUsecase,
	settingsRepo repository.SettingsRepository,
) *BookingUsecase {
	return &BookingUsecase{
		bookingRepo:         bookingRepo,
		eventRepo:           eventRepo,
		roleRepo:            roleRepo,
		notificationUsecase: notificationUsecase,
		settingsRepo:        settingsRepo,
	}
}

func (u *BookingUsecase) ReserveTickets(ctx context.Context, userID string, eventID string, requests []domain.TicketRequest) (*domain.Booking, error) {
	event, err := u.eventRepo.GetEventByID(eventID)
	if err != nil {
		return nil, errors.New("event not found")
	}

	if event.Status != "approved" && event.Status != "live" {
		return nil, errors.New("EVT_012: Event not live")
	}

	ticketTypes, err := u.eventRepo.GetTicketTypesByEventID(eventID)
	if err != nil {
		return nil, err
	}

	ticketMap := make(map[string]domain.TicketType)
	for _, tt := range ticketTypes {
		ticketMap[tt.ID] = tt
	}

	var baseTotal float64
	var bookingTickets []domain.BookingTicket
	bookingID := uuid.New().String()

	for _, req := range requests {
		tt, exists := ticketMap[req.TicketTypeID]
		if !exists {
			return nil, apiErrors.ErrInvalidRequestBody
		}

		baseTotal += tt.Price * float64(req.Quantity)

		bookingTickets = append(bookingTickets, domain.BookingTicket{
			BookingID:    bookingID,
			TicketTypeID: req.TicketTypeID,
			Quantity:     req.Quantity,
		})
	}

	now := time.Now()
	expiresAt := now.Add(1 * time.Minute)

	var totalTickets int
	for _, req := range requests {
		totalTickets += req.Quantity
	}

	userFee := 0.0
	if baseTotal > 0 {
		userFee = float64(30 * totalTickets)
	}
	totalAmount := baseTotal + userFee

	booking := &domain.Booking{
		ID:          bookingID,
		UserID:      userID,
		EventID:     eventID,
		Status:      "reserved",
		TotalAmount: totalAmount,
		ExpiresAt:   expiresAt,
		CreatedAt:   now,
	}

	err = u.bookingRepo.ReserveTickets(ctx, booking, bookingTickets)
	if err != nil {
		return nil, err
	}

	if event.Status == "approved" {
		_ = u.eventRepo.UpdateEventStatus(ctx, eventID, "live")
		logger.Log.Info().
			Str("event", "event_state_changed").
			Str("entity", "event").
			Str("entity_id", eventID).
			Str("from", "approved").
			Str("to", "live").
			Str("actor_id", userID).
			Msg("")
	}

	logger.Log.Info().
		Str("event", "booking_reserved").
		Str("booking_id", booking.ID).
		Str("user_id", booking.UserID).
		Str("event_id", booking.EventID).
		Interface("tickets", requests).
		Float64("total_amount", booking.TotalAmount).
		Time("expires_at", booking.ExpiresAt).
		Msg("")

	if u.notificationUsecase != nil {
		if notifyErr := u.notificationUsecase.SendNotification(
			userID,
			domain.NotificationTypeBookingReserved,
			"Booking reserved",
			"Booking reserved. Complete payment before expiry.",
			map[string]interface{}{
				"booking_id":  booking.ID,
				"event_id":    booking.EventID,
				"event_title": event.Title,
				"expires_at":  booking.ExpiresAt,
			},
		); notifyErr != nil {
			logger.Log.Warn().
				Err(notifyErr).
				Str("user_id", userID).
				Str("booking_id", booking.ID).
				Msg("notification_send_failed")
		}
	}

	return booking, nil
}

func (u *BookingUsecase) ProcessExpiredBookings(ctx context.Context) error {
	expiredBookings, err := u.bookingRepo.ExpireBookings(ctx)
	if err != nil {
		return err
	}

	for _, b := range expiredBookings {
		logger.Log.Info().
			Str("booking_id", b.ID).
			Str("event_id", b.EventID).
			Str("user_id", b.UserID).
			Time("expired_at", b.ExpiresAt).
			Time("timestamp", time.Now()).
			Msg("booking_expired")
	}

	return nil
}

func (u *BookingUsecase) GetUserBookings(ctx context.Context, userID string, page int, limit int, status string) ([]domain.BookingWithEvent, int64, error) {
	bookings, total, err := u.bookingRepo.GetUserBookings(ctx, userID, page, limit, status)
	if err != nil {
		return nil, 0, err
	}

	return bookings, total, nil
}

func (u *BookingUsecase) GetUserTickets(ctx context.Context, userID string, eventID string, bookingID string, status string) ([]domain.TicketWithEvent, error) {
	tickets, err := u.bookingRepo.GetUserTickets(ctx, userID, eventID, bookingID, status)
	if err != nil {
		return nil, err
	}

	return tickets, nil
}

func (u *BookingUsecase) CancelBooking(ctx context.Context, bookingID string, userID string, items []domain.TicketCancelRequest) error {
	booking, err := u.bookingRepo.FindByID(ctx, bookingID)
	if err != nil {
		return err
	}

	event, err := u.eventRepo.GetEventByID(booking.EventID)
	if err != nil {
		return err
	}

	settings, settingsErr := u.settingsRepo.GetPlatformSettings()
	if settingsErr == nil && !settings.AllowEventCancellation {
		return apiErrors.New(403, apiErrors.ForbiddenAction, "Event cancellation is currently disabled by admin")
	}

	refundWindowDays := 1
	if settingsErr == nil && settings.RefundWindowDays >= 0 {
		refundWindowDays = settings.RefundWindowDays
	}

	refundWindowDuration := time.Duration(refundWindowDays) * 24 * time.Hour
	timeUntilEvent := time.Until(event.StartTime)
	isRefundable := timeUntilEvent >= refundWindowDuration

	err = u.bookingRepo.CancelBooking(ctx, bookingID, userID, items, isRefundable)
	if err != nil {
		return err
	}

	logger.Log.Info().
		Str("booking_id", bookingID).
		Str("user_id", userID).
		Int("items_count", len(items)).
		Time("timestamp", time.Now()).
		Msg("booking_cancelled_partially")

	return nil
}

func (u *BookingUsecase) CheckInTicket(
	ctx context.Context,
	eventID string,
	actorID string,
	ticketCode string,
) (*domain.TicketCheckIn, error) {
	roles, err := u.roleRepo.GetRolesByUserID(actorID)
	if err != nil {
		return nil, err
	}

	isAdmin := false
	isOrganizer := false
	for _, role := range roles {
		switch role {
		case domain.RoleAdmin:
			isAdmin = true
		case domain.RoleOrganizer:
			isOrganizer = true
		}
	}

	if !isAdmin {
		if !isOrganizer {
			return nil, apiErrors.ErrForbiddenAction
		}

		event, err := u.eventRepo.GetEventByID(eventID)
		if err != nil {
			return nil, err
		}

		if event.OrganizerID != actorID {
			return nil, apiErrors.ErrForbiddenAction
		}
	}

	return u.bookingRepo.CheckInTicket(ctx, eventID, ticketCode)
}

func (u *BookingUsecase) PayWithWallet(ctx context.Context, bookingID string, userID string) error {
	booking, err := u.bookingRepo.FindByID(ctx, bookingID)
	if err != nil {
		return err
	}

	if booking.UserID != userID {
		return apiErrors.ErrForbiddenAction
	}

	walletSetting, err := u.settingsRepo.GetPaymentProviderConfig("wallet")
	if err == nil && !walletSetting.IsEnabled {
		return apiErrors.New(400, "PAYMENT_DISABLED", "Wallet payment is disabled by admin")
	}

	if booking.Status != "reserved" {
		return apiErrors.ErrInvalidStateTransition
	}

	if time.Now().After(booking.ExpiresAt) {
		return apiErrors.ErrBookingExpired
	}

	err = u.bookingRepo.PayWithWallet(ctx, bookingID, userID, booking.TotalAmount)
	if err != nil {
		return err
	}

	if u.notificationUsecase != nil {
		event, _ := u.eventRepo.GetEventByID(booking.EventID)
		_ = u.notificationUsecase.SendNotification(
			userID,
			domain.NotificationTypePaymentSuccess,
			"Payment successful",
			"Payment successful via wallet. Tickets confirmed.",
			map[string]interface{}{
				"booking_id":  booking.ID,
				"event_id":    booking.EventID,
				"event_title": event.Title,
				"amount":      booking.TotalAmount,
			},
		)

		_ = u.notificationUsecase.SendNotification(
			userID,
			domain.NotificationTypeTicketGenerated,
			"Your tickets are generated",
			"Your tickets are generated",
			map[string]interface{}{
				"booking_id":  booking.ID,
				"event_id":    booking.EventID,
				"event_title": event.Title,
			},
		)
	}

	return nil
}
