package repository

import (
	"context"
	"errors"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BookingModel struct {
	ID          string
	UserID      string
	EventID     string
	Status      string
	TotalAmount float64
	ExpiresAt   int64
	CreatedAt   int64
}

type BookingTicketModel struct {
	BookingID    string
	TicketTypeID string
	Quantity     int
}

type bookingGormRepository struct {
	db *gorm.DB
}

func NewBookingGormRepository(db *gorm.DB) *bookingGormRepository {
	return &bookingGormRepository{db: db}
}

func (r *bookingGormRepository) ReserveTickets(ctx context.Context, booking *domain.Booking, tickets []domain.BookingTicket) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Lock and fetch ticket types to prevent race conditions during inventory deduction
		for _, reqTicket := range tickets {
			var ticketModel TicketTypeModel

			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("id = ?", reqTicket.TicketTypeID).
				First(&ticketModel).Error; err != nil {
				return err
			}

			// Validate inventory
			if ticketModel.AvailableQuantity < reqTicket.Quantity {
				return errors.New("EVT_009: Ticket sold out")
			}

			// Deduct inventory
			if err := tx.Model(&TicketTypeModel{}).
				Where("id = ?", reqTicket.TicketTypeID).
				Update("available_quantity", gorm.Expr("available_quantity - ?", reqTicket.Quantity)).Error; err != nil {
				return err
			}
		}

		// 2. Create the booking record
		bookingModel := BookingModel{
			ID:          booking.ID,
			UserID:      booking.UserID,
			EventID:     booking.EventID,
			Status:      booking.Status, // should be "reserved"
			TotalAmount: booking.TotalAmount,
			ExpiresAt:   booking.ExpiresAt.Unix(),
			CreatedAt:   booking.CreatedAt.Unix(),
		}

		if err := tx.Create(&bookingModel).Error; err != nil {
			return err
		}

		// 3. Create the booking ticket associations
		for _, ticket := range tickets {
			btModel := BookingTicketModel{
				BookingID:    ticket.BookingID,
				TicketTypeID: ticket.TicketTypeID,
				Quantity:     ticket.Quantity,
			}
			if err := tx.Create(&btModel).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
