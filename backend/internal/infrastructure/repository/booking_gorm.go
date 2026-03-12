package repository

import (
	"context"
	"errors"
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/pkg/logger"
	"gorm.io/gorm"
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
	maxRetries := 3

	for attempt := 0; attempt < maxRetries; attempt++ {
		err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			for _, reqTicket := range tickets {
				var ticketModel TicketTypeModel

				// 1. Fetch current version and inventory without pessimistic lock
				if err := tx.Where("id = ?", reqTicket.TicketTypeID).First(&ticketModel).Error; err != nil {
					return err
				}

				// Validate inventory
				if ticketModel.AvailableQuantity < reqTicket.Quantity {
					return errors.New("EVT_009: Ticket sold out")
				}

				// 2. Optimistic update with version check
				res := tx.Model(&TicketTypeModel{}).
					Where("id = ? AND version = ? AND available_quantity >= ?", ticketModel.ID, ticketModel.Version, reqTicket.Quantity).
					Updates(map[string]interface{}{
						"available_quantity": gorm.Expr("available_quantity - ?", reqTicket.Quantity),
						"version":            gorm.Expr("version + 1"),
					})

				if res.Error != nil {
					return res.Error
				}

				if res.RowsAffected == 0 {
					logger.Log.Warn().
						Str("ticket_type_id", ticketModel.ID).
						Int("requested_quantity", reqTicket.Quantity).
						Time("timestamp", time.Now()).
						Msg("inventory_conflict")

					return errors.New("conflict")
				}
			}

			// 3. Create the booking record
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

			// 4. Create the booking ticket associations
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

		if err == nil {
			return nil // Success
		}

		// If it's specifically sold out, don't retry, fail fast
		if errors.Is(err, errors.New("EVT_009: Ticket sold out")) || err.Error() == "EVT_009: Ticket sold out" {
			return err
		}

		// If it's a conflict, retry on next loop iteration
		if err.Error() == "conflict" {
			continue
		}

		// Fallback native error pass-out
		return err
	}

	return errors.New("EVT_009: Ticket sold out")
}


func (r *bookingGormRepository) ExpireBookings(ctx context.Context) ([]domain.Booking, error) {
	var expiredBookings []BookingModel
	var returnedBookings []domain.Booking

	// Fetch up to 100 expired bookings
	err := r.db.WithContext(ctx).Where("status = ? AND expires_at < ?", "reserved", time.Now().Unix()).
		Limit(100).
		Find(&expiredBookings).Error

	if err != nil {
		return nil, err
	}

	for _, bm := range expiredBookings {
		err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			// Update status to 'expired'
			if err := tx.Model(&BookingModel{}).Where("id = ?", bm.ID).Update("status", "expired").Error; err != nil {
				return err
			}

			// Fetch ticket relations
			var bTickets []BookingTicketModel
			if err := tx.Where("booking_id = ?", bm.ID).Find(&bTickets).Error; err != nil {
				return err
			}

			// Restore ticket quantities
			for _, bt := range bTickets {
				if err := tx.Model(&TicketTypeModel{}).
					Where("id = ?", bt.TicketTypeID).
					Update("available_quantity", gorm.Expr("available_quantity + ?", bt.Quantity)).Error; err != nil {
					return err
				}
			}

			// Map for return
			returnedBookings = append(returnedBookings, domain.Booking{
				ID:          bm.ID,
				UserID:      bm.UserID,
				EventID:     bm.EventID,
				Status:      "expired",
				TotalAmount: bm.TotalAmount,
				ExpiresAt:   time.Unix(bm.ExpiresAt, 0),
				CreatedAt:   time.Unix(bm.CreatedAt, 0),
			})

			return nil
		})

		if err != nil {
			// Log this or continue, returning what we have processed so far. We'll proceed with processing others.
			continue
		}
	}

	return returnedBookings, nil
}

