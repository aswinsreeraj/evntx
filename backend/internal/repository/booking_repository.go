package repository

import (
	"context"

	"github.com/aswinsreeraj/evntx/internal/domain"
)

type BookingRepository interface {
	ReserveTickets(ctx context.Context, booking *domain.Booking, tickets []domain.BookingTicket) error
	CancelBooking(ctx context.Context, bookingID string, userID string, items []domain.TicketCancelRequest) error
	ExpireBookings(ctx context.Context) ([]domain.Booking, error)
	GetUserBookings(ctx context.Context, userID string, page int, limit int, status string) ([]domain.BookingWithEvent, int64, error)
	GetUserTickets(ctx context.Context, userID string, eventID string, bookingID string, status string) ([]domain.TicketWithEvent, error)
}



