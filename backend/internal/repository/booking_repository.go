package repository

import (
	"context"

	"github.com/aswinsreeraj/evntx/internal/domain"
)

type BookingRepository interface {
	ReserveTickets(ctx context.Context, booking *domain.Booking, tickets []domain.BookingTicket) error
	FindByID(ctx context.Context, bookingID string) (*domain.Booking, error)
	CancelBooking(ctx context.Context, bookingID string, userID string, items []domain.TicketCancelRequest, isRefundable bool) error
	ExpireBookings(ctx context.Context) ([]domain.Booking, error)
	GetPaidBookingsByEventID(ctx context.Context, eventID string) ([]domain.Booking, error)
	GetUserBookings(ctx context.Context, userID string, page int, limit int, status string) ([]domain.BookingWithEvent, int64, error)
	GetUserTickets(ctx context.Context, userID string, eventID string, bookingID string, status string) ([]domain.TicketWithEvent, error)
	CheckInTicket(ctx context.Context, eventID string, ticketCode string) (*domain.TicketCheckIn, error)
	GetTicketCountByBookingID(ctx context.Context, bookingID string) (int, error)
	GetBookingContextsByIDs(ctx context.Context, bookingIDs []string) (map[string]domain.BookingContextDetails, error)
	PayWithWallet(ctx context.Context, bookingID string, userID string, amount float64) error
}
