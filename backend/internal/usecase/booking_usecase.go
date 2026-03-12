package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/internal/repository"
	"github.com/aswinsreeraj/evntx/pkg/logger"
	"github.com/google/uuid"
)

type BookingUsecase struct {
	bookingRepo repository.BookingRepository
	eventRepo   repository.EventRepository
}

func NewBookingUsecase(bookingRepo repository.BookingRepository, eventRepo repository.EventRepository) *BookingUsecase {
	return &BookingUsecase{
		bookingRepo: bookingRepo,
		eventRepo:   eventRepo,
	}
}

func (u *BookingUsecase) ReserveTickets(ctx context.Context, userID string, eventID string, requests []domain.TicketRequest) (*domain.Booking, error) {
	// 1. Fetch Event and validate state
	event, err := u.eventRepo.GetEventByID(eventID)
	if err != nil {
		return nil, errors.New("event not found")
	}

	if event.Status != "approved" && event.Status != "live" {
		return nil, errors.New("EVT_012: Event not live")
	}

	// 2. Fetch ticket types to calculate total amount
	ticketTypes, err := u.eventRepo.GetTicketTypesByEventID(eventID)
	if err != nil {
		return nil, err
	}

	// Create a map for quick lookup
	ticketMap := make(map[string]domain.TicketType)
	for _, tt := range ticketTypes {
		ticketMap[tt.ID] = tt
	}

	var totalAmount float64
	var bookingTickets []domain.BookingTicket
	bookingID := uuid.New().String()

	for _, req := range requests {
		tt, exists := ticketMap[req.TicketTypeID]
		if !exists {
			return nil, errors.New("invalid ticket type for this event")
		}

		totalAmount += tt.Price * float64(req.Quantity)

		bookingTickets = append(bookingTickets, domain.BookingTicket{
			BookingID:    bookingID,
			TicketTypeID: req.TicketTypeID,
			Quantity:     req.Quantity,
		})
	}

	now := time.Now()
	// Booking expiration: 10 minutes
	expiresAt := now.Add(10 * time.Minute)

	booking := &domain.Booking{
		ID:          bookingID,
		UserID:      userID,
		EventID:     eventID,
		Status:      "reserved",
		TotalAmount: totalAmount,
		ExpiresAt:   expiresAt,
		CreatedAt:   now,
	}

	// 3. Trigger Repository Transaction
	err = u.bookingRepo.ReserveTickets(ctx, booking, bookingTickets)
	if err != nil {
		if err.Error() == "EVT_009: Ticket sold out" {
			return nil, errors.New("EVT_009: Tickets sold out") // Match required message exactly
		}
		return nil, err
	}

	logger.Log.Info().
		Str("booking_id", booking.ID).
		Str("user_id", booking.UserID).
		Str("event_id", booking.EventID).
		Float64("total_amount", booking.TotalAmount).
		Time("expires_at", booking.ExpiresAt).
		Time("timestamp", now).
		Msg("booking_reserved")

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

	logger.Log.Info().
		Str("user_id", userID).
		Int("page", page).
		Int("limit", limit).
		Time("timestamp", time.Now()).
		Msg("user_bookings_fetched")

	return bookings, total, nil
}
