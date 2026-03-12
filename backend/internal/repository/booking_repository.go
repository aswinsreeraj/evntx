package repository

import (
	"context"

	"github.com/aswinsreeraj/evntx/internal/domain"
)

type BookingRepository interface {
	ReserveTickets(ctx context.Context, booking *domain.Booking, tickets []domain.BookingTicket) error
	ExpireBookings(ctx context.Context) ([]domain.Booking, error)
}

