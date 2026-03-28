package repository

import (
	"context"
	"errors"
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
	apiErrors "github.com/aswinsreeraj/evntx/pkg/errors"
	"github.com/aswinsreeraj/evntx/pkg/logger"
	"github.com/google/uuid"
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

type TicketModel struct {
	ID           string
	BookingID    string
	TicketTypeID string
	TicketCode   string
	QRPayload    string
	Status       string
	CheckedInAt  *int64
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

				if err := tx.Where("id = ?", reqTicket.TicketTypeID).First(&ticketModel).Error; err != nil {
					return err
				}

				if ticketModel.AvailableQuantity < reqTicket.Quantity {
					return apiErrors.ErrTicketSoldOut
				}

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

			bookingModel := BookingModel{
				ID:          booking.ID,
				UserID:      booking.UserID,
				EventID:     booking.EventID,
				Status:      booking.Status,
				TotalAmount: booking.TotalAmount,
				ExpiresAt:   booking.ExpiresAt.Unix(),
				CreatedAt:   booking.CreatedAt.Unix(),
			}

			if err := tx.Create(&bookingModel).Error; err != nil {
				return err
			}

			for _, ticket := range tickets {
				btModel := BookingTicketModel{
					BookingID:    ticket.BookingID,
					TicketTypeID: ticket.TicketTypeID,
					Quantity:     ticket.Quantity,
				}
				if err := tx.Create(&btModel).Error; err != nil {
					return err
				}

				for i := 0; i < ticket.Quantity; i++ {
					newTktId := uuid.New().String()
					tktModel := TicketModel{
						ID:           newTktId,
						BookingID:    ticket.BookingID,
						TicketTypeID: ticket.TicketTypeID,
						TicketCode:   "TKT-" + newTktId[:8],
						QRPayload:    newTktId,
						Status:       "valid",
					}
					if err := tx.Create(&tktModel).Error; err != nil {
						return err
					}
				}
			}

			return nil
		})

		if err == nil {
			return nil
		}

		if errors.Is(err, apiErrors.ErrTicketSoldOut) {
			return err
		}

		if err.Error() == "conflict" {
			continue
		}

		return err
	}

	return apiErrors.ErrTicketSoldOut
}

func (r *bookingGormRepository) FindByID(ctx context.Context, bookingID string) (*domain.Booking, error) {
	var model BookingModel

	if err := r.db.WithContext(ctx).Where("id = ?", bookingID).First(&model).Error; err != nil {
		return nil, err
	}

	return &domain.Booking{
		ID:          model.ID,
		UserID:      model.UserID,
		EventID:     model.EventID,
		Status:      model.Status,
		TotalAmount: model.TotalAmount,
		ExpiresAt:   time.Unix(model.ExpiresAt, 0),
		CreatedAt:   time.Unix(model.CreatedAt, 0),
	}, nil
}

func (r *bookingGormRepository) CancelBooking(ctx context.Context, bookingID string, userID string, items []domain.TicketCancelRequest) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var bm BookingModel
		if err := tx.Where("id = ? AND user_id = ?", bookingID, userID).First(&bm).Error; err != nil {
			return apiErrors.ErrResourceNotFound
		}

		if bm.Status != "paid" && bm.Status != "confirmed" && bm.Status != "reserved" {
			return apiErrors.ErrInvalidStateTransition
		}

		var totalRefund float64

		for _, item := range items {
			if item.Quantity <= 0 {
				continue
			}

			var tt TicketTypeModel
			if err := tx.Where("name = ? AND event_id = ?", item.TicketType, bm.EventID).First(&tt).Error; err != nil {
				return apiErrors.ErrInvalidRequestBody
			}

			var tM []TicketModel
			if err := tx.Where("booking_id = ? AND ticket_type_id = ? AND status != 'cancelled'", bookingID, tt.ID).
				Limit(item.Quantity).
				Find(&tM).Error; err != nil {
				return err
			}

			if len(tM) < item.Quantity {
				return apiErrors.ErrInvalidRequestBody
			}

			var idsToCancel []string
			for _, t := range tM {
				idsToCancel = append(idsToCancel, t.ID)
			}

			if err := tx.Model(&TicketModel{}).Where("id IN ?", idsToCancel).Update("status", "cancelled").Error; err != nil {
				return err
			}

			if err := tx.Model(&BookingTicketModel{}).
				Where("booking_id = ? AND ticket_type_id = ?", bookingID, tt.ID).
				Update("quantity", gorm.Expr("quantity - ?", item.Quantity)).Error; err != nil {
				return err
			}

			if err := tx.Model(&TicketTypeModel{}).
				Where("id = ?", tt.ID).
				Update("available_quantity", gorm.Expr("available_quantity + ?", item.Quantity)).Error; err != nil {
				return err
			}

			totalRefund += tt.Price * float64(item.Quantity)
		}

		if totalRefund > 0 {
			if err := tx.Model(&BookingModel{}).Where("id = ?", bookingID).
				Update("total_amount", gorm.Expr("total_amount - ?", totalRefund)).Error; err != nil {
				return err
			}
		}

		var remainingTickets int64
		if err := tx.Model(&TicketModel{}).Where("booking_id = ? AND status != 'cancelled'", bookingID).Count(&remainingTickets).Error; err != nil {
			return err
		}

		if remainingTickets == 0 {
			if err := tx.Model(&BookingModel{}).Where("id = ?", bookingID).Update("status", "cancelled").Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *bookingGormRepository) ExpireBookings(ctx context.Context) ([]domain.Booking, error) {
	var expiredBookings []BookingModel
	var returnedBookings []domain.Booking

	err := r.db.WithContext(ctx).Where("status = ? AND expires_at < ?", "reserved", time.Now().Unix()).
		Limit(100).
		Find(&expiredBookings).Error

	if err != nil {
		return nil, err
	}

	for _, bm := range expiredBookings {
		err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := tx.Model(&BookingModel{}).Where("id = ?", bm.ID).Update("status", "expired").Error; err != nil {
				return err
			}

			var bTickets []BookingTicketModel
			if err := tx.Where("booking_id = ?", bm.ID).Find(&bTickets).Error; err != nil {
				return err
			}

			for _, bt := range bTickets {
				if err := tx.Model(&TicketTypeModel{}).
					Where("id = ?", bt.TicketTypeID).
					Update("available_quantity", gorm.Expr("available_quantity + ?", bt.Quantity)).Error; err != nil {
					return err
				}
			}

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
			continue
		}
	}

	return returnedBookings, nil
}

func (r *bookingGormRepository) GetPaidBookingsByEventID(ctx context.Context, eventID string) ([]domain.Booking, error) {
	var models []BookingModel

	if err := r.db.WithContext(ctx).
		Where("event_id = ? AND status = ?", eventID, "paid").
		Find(&models).Error; err != nil {
		return nil, err
	}

	bookings := make([]domain.Booking, 0, len(models))
	for _, model := range models {
		bookings = append(bookings, domain.Booking{
			ID:          model.ID,
			UserID:      model.UserID,
			EventID:     model.EventID,
			Status:      model.Status,
			TotalAmount: model.TotalAmount,
			ExpiresAt:   time.Unix(model.ExpiresAt, 0),
			CreatedAt:   time.Unix(model.CreatedAt, 0),
		})
	}

	return bookings, nil
}

func (r *bookingGormRepository) GetUserBookings(ctx context.Context, userID string, page int, limit int, status string) ([]domain.BookingWithEvent, int64, error) {
	var total int64
	query := r.db.WithContext(ctx).Table("booking_models").Where("booking_models.user_id = ?", userID)

	if status != "" {
		query = query.Where("booking_models.status = ?", status)
	} else {
		query = query.Where("booking_models.status NOT IN (?)", []string{"reserved", "expired"})
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	var results []struct {
		BookingID      string
		EventID        string
		EventTitle     string
		EventCity      string
		EventStartTime int64
		Status         string
		TotalAmount    float64
		TicketCount    int
		CreatedAt      int64
		CoverImageURL  string
		VenueName      string
		Tags           string
	}

	err := query.Select(`
		booking_models.id AS booking_id, 
		booking_models.event_id, 
		event_models.title AS event_title, 
		event_models.city AS event_city, 
		event_models.start_time AS event_start_time, 
		booking_models.status AS status, 
		booking_models.total_amount, 
		booking_models.created_at, 
		event_models.cover_image_url,
		event_models.venue_name,
		event_models.tags,
		COALESCE(SUM(booking_ticket_models.quantity), 0) AS ticket_count
	`).
		Joins("JOIN event_models ON event_models.id = booking_models.event_id").
		Joins("LEFT JOIN booking_ticket_models ON booking_ticket_models.booking_id = booking_models.id").
		Group("booking_models.id, event_models.id").
		Order("booking_models.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&results).Error

	if err != nil {
		return nil, 0, err
	}

	bookings := make([]domain.BookingWithEvent, 0, len(results))
	for _, r := range results {
		bookings = append(bookings, domain.BookingWithEvent{
			BookingID:      r.BookingID,
			EventID:        r.EventID,
			EventTitle:     r.EventTitle,
			EventCity:      r.EventCity,
			EventStartTime: time.Unix(r.EventStartTime, 0),
			Status:         r.Status,
			TotalAmount:    r.TotalAmount,
			TicketCount:    r.TicketCount,
			CreatedAt:      time.Unix(r.CreatedAt, 0),
			CoverImageURL:  r.CoverImageURL,
			VenueName:      r.VenueName,
			Tags:           r.Tags,
		})
	}

	return bookings, total, nil
}

func (r *bookingGormRepository) GetUserTickets(ctx context.Context, userID string, eventID string, bookingID string, status string) ([]domain.TicketWithEvent, error) {
	query := r.db.WithContext(ctx).Table("ticket_models").
		Joins("JOIN booking_models ON booking_models.id = ticket_models.booking_id").
		Joins("JOIN event_models ON event_models.id = booking_models.event_id").
		Joins("JOIN ticket_type_models ON ticket_type_models.id = ticket_models.ticket_type_id").
		Where("booking_models.user_id = ?", userID)

	if bookingID != "" {
		query = query.Where("booking_models.id = ?", bookingID)
	}
	if eventID != "" {
		query = query.Where("event_models.id = ?", eventID)
	}
	if status != "" {
		query = query.Where("ticket_models.status = ?", status)
	} else {
		query = query.Where("ticket_models.status != ?", "cancelled")
	}

	// Only show tickets for paid or confirmed bookings
	query = query.Where("booking_models.status IN (?)", []string{"paid", "confirmed"})

	var results []struct {
		TicketID    string
		TicketCode  string
		EventID     string
		EventTitle  string
		TicketType  string
		Status      string
		CheckedInAt *int64
		CreatedAt   int64 `gorm:"column:created_at"`
	}

	err := query.Select(`
		ticket_models.id AS ticket_id,
		ticket_models.ticket_code,
		event_models.id AS event_id,
		event_models.title AS event_title,
		ticket_type_models.name AS ticket_type,
		ticket_models.status,
		ticket_models.checked_in_at,
		booking_models.created_at AS created_at
	`).
		Order("booking_models.created_at DESC").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	tickets := make([]domain.TicketWithEvent, 0, len(results))
	for _, r := range results {
		var checkedIn time.Time
		var checkedInPtr *time.Time
		if r.CheckedInAt != nil {
			checkedIn = time.Unix(*r.CheckedInAt, 0)
			checkedInPtr = &checkedIn
		}

		tickets = append(tickets, domain.TicketWithEvent{
			TicketID:    r.TicketID,
			TicketCode:  r.TicketCode,
			EventID:     r.EventID,
			EventTitle:  r.EventTitle,
			TicketType:  r.TicketType,
			Status:      r.Status,
			CheckedInAt: checkedInPtr,
		})
	}

	return tickets, nil
}
