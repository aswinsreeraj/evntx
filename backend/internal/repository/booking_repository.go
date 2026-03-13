package repository

import (
	"context"

	"github.com/aswinsreeraj/evntx/internal/domain"
)

type BookingRepository interface {
	ReserveTickets(ctx context.Context, booking *domain.Booking, tickets []domain.BookingTicket) error
	ExpireBookings(ctx context.Context) ([]domain.Booking, error)
	GetUserBookings(ctx context.Context, userID string, page int, limit int, status string) ([]domain.BookingWithEvent, int64, error)
	GetUserTickets(ctx context.Context, userID string, eventID string, status string) ([]domain.TicketWithEvent, error)
}



