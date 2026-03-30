package repository

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
	apiErrors "github.com/aswinsreeraj/evntx/pkg/errors"
	"github.com/aswinsreeraj/evntx/pkg/logger"
	"github.com/google/uuid"
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

func (r *bookingGormRepository) CancelBooking(ctx context.Context, bookingID string, userID string, items []domain.TicketCancelRequest, isRefundable bool) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var bm BookingModel
		if err := tx.Where("id = ? AND user_id = ?", bookingID, userID).First(&bm).Error; err != nil {
			return apiErrors.ErrResourceNotFound
		}

		if bm.Status != "paid" && bm.Status != "confirmed" && bm.Status != "reserved" {
			return apiErrors.ErrInvalidStateTransition
		}

		var totalRefund float64
		var totalTicketsCancelled int

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
			totalTicketsCancelled += item.Quantity
		}

		if (bm.Status == "paid" || bm.Status == "confirmed") && totalRefund > 0 && isRefundable {
			now := time.Now()

			var userWallet WalletModel
			if err := tx.Where("user_id = ?", bm.UserID).First(&userWallet).Error; err != nil {
				return err
			}

			if err := tx.Create(&WalletTransactionModel{
				ID:            uuid.NewString(),
				WalletID:      userWallet.ID,
				Type:          domain.WalletTransactionTypeCredit,
				Amount:        totalRefund,
				ReferenceType: domain.WalletReferenceTypeUserCancellation,
				ReferenceID:   bookingID,
				Status:        domain.WalletTransactionStatusCompleted,
				CreatedAt:     now,
			}).Error; err != nil {
				return err
			}

			userWallet.AvailableBalance = math.Round((userWallet.AvailableBalance+totalRefund)*100) / 100
			if err := tx.Model(&WalletModel{}).Where("id = ?", userWallet.ID).Updates(map[string]interface{}{
				"available_balance": userWallet.AvailableBalance,
				"updated_at":        now,
			}).Error; err != nil {
				return err
			}

			var event EventModel
			if err := tx.Where("id = ?", bm.EventID).First(&event).Error; err != nil {
				return err
			}

			var orgWallet WalletModel
			if err := tx.Where("user_id = ?", event.OrganizerID).First(&orgWallet).Error; err != nil {
				return err
			}

			if err := tx.Create(&WalletTransactionModel{
				ID:            uuid.NewString(),
				WalletID:      orgWallet.ID,
				Type:          domain.WalletTransactionTypeDebit,
				Amount:        totalRefund,
				ReferenceType: domain.WalletReferenceTypeUserCancellation,
				ReferenceID:   bookingID,
				Status:        domain.WalletTransactionStatusCompleted,
				CreatedAt:     now,
			}).Error; err != nil {
				return err
			}

			orgWallet.PendingBalance = math.Round((orgWallet.PendingBalance-totalRefund)*100) / 100
			orgWallet.ReserveBalance = math.Round((orgWallet.ReserveBalance+float64(totalTicketsCancelled*30))*100) / 100
			if err := tx.Model(&WalletModel{}).Where("id = ?", orgWallet.ID).Updates(map[string]interface{}{
				"pending_balance": orgWallet.PendingBalance,
				"reserve_balance": orgWallet.ReserveBalance,
				"updated_at":      now,
			}).Error; err != nil {
				return err
			}
		}

		if totalRefund > 0 && isRefundable {
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
		EventStatus    string
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
		event_models.status AS event_status,
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
			EventStatus:    r.EventStatus,
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

func (r *bookingGormRepository) CheckInTicket(
	ctx context.Context,
	eventID string,
	ticketCode string,
) (*domain.TicketCheckIn, error) {
	var checkedInTicket *domain.TicketCheckIn

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var result struct {
			TicketID    string
			TicketCode  string
			EventID     string
			Status      string
			CheckedInAt *int64
		}

		if err := tx.Table("ticket_models").
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Select(`
				ticket_models.id AS ticket_id,
				ticket_models.ticket_code,
				booking_models.event_id,
				ticket_models.status,
				ticket_models.checked_in_at
			`).
			Joins("JOIN booking_models ON booking_models.id = ticket_models.booking_id").
			Where("ticket_models.ticket_code = ?", ticketCode).
			Take(&result).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apiErrors.New(404, apiErrors.ResourceNotFound, "Ticket not found")
			}

			return err
		}

		if result.EventID != eventID {
			return apiErrors.New(404, apiErrors.ResourceNotFound, "Ticket not found")
		}

		if result.Status == "used" {
			return apiErrors.New(409, apiErrors.InvalidStateTransition, "Ticket already used")
		}

		if result.Status != "valid" {
			return apiErrors.New(400, apiErrors.InvalidStateTransition, "Ticket is not valid for check-in")
		}

		now := time.Now()
		checkedInAt := now.Unix()

		updateResult := tx.Model(&TicketModel{}).
			Where("id = ? AND status = ?", result.TicketID, "valid").
			Updates(map[string]interface{}{
				"status":        "used",
				"checked_in_at": checkedInAt,
			})
		if updateResult.Error != nil {
			return updateResult.Error
		}
		if updateResult.RowsAffected == 0 {
			return apiErrors.New(409, apiErrors.InvalidStateTransition, "Ticket already used")
		}

		checkedInTicket = &domain.TicketCheckIn{
			TicketID:    result.TicketID,
			TicketCode:  result.TicketCode,
			Status:      "used",
			CheckedInAt: now,
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return checkedInTicket, nil
}

func (r *bookingGormRepository) GetTicketCountByBookingID(ctx context.Context, bookingID string) (int, error) {
	var models []BookingTicketModel
	if err := r.db.WithContext(ctx).Where("booking_id = ?", bookingID).Find(&models).Error; err != nil {
		return 0, err
	}
	total := 0
	for _, m := range models {
		total += m.Quantity
	}
	return total, nil
}

func (r *bookingGormRepository) PayWithWallet(ctx context.Context, bookingID string, userID string, amount float64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var booking BookingModel
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND user_id = ? AND status = ?", bookingID, userID, "reserved").First(&booking).Error; err != nil {
			return err
		}

		if math.Abs(booking.TotalAmount-amount) > 0.01 {
			return apiErrors.New(400, apiErrors.InvalidRequestBody, "amount mismatch")
		}

		var wallet WalletModel
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", userID).First(&wallet).Error; err != nil {
			return err
		}

		if wallet.AvailableBalance < amount {
			return apiErrors.ErrInsufficientBalance
		}

		now := time.Now()
		normalizedAmount := math.Round(amount*100) / 100

		// 1. Deduct from User Wallet
		if err := tx.Create(&WalletTransactionModel{
			ID:            uuid.NewString(),
			WalletID:      wallet.ID,
			Type:          domain.WalletTransactionTypeDebit,
			Amount:        normalizedAmount,
			ReferenceType: domain.WalletReferenceTypePurchase,
			ReferenceID:   bookingID,
			Status:        domain.WalletTransactionStatusCompleted,
			CreatedAt:     now,
		}).Error; err != nil {
			return err
		}

		wallet.AvailableBalance = math.Round((wallet.AvailableBalance-normalizedAmount)*100) / 100
		wallet.TotalDebited = math.Round((wallet.TotalDebited+normalizedAmount)*100) / 100
		wallet.UpdatedAt = now

		if err := tx.Model(&WalletModel{}).Where("id = ?", wallet.ID).Updates(map[string]interface{}{
			"available_balance": wallet.AvailableBalance,
			"total_debited":     wallet.TotalDebited,
			"updated_at":        wallet.UpdatedAt,
		}).Error; err != nil {
			return err
		}

		// 2. Update Booking Status
		if err := tx.Model(&BookingModel{}).Where("id = ?", bookingID).Update("status", "paid").Error; err != nil {
			return err
		}

		// 3. Credit Platform & Organizer
		var totalTickets int64
		if err := tx.Model(&BookingTicketModel{}).Where("booking_id = ?", bookingID).Select("COALESCE(SUM(quantity), 0)").Scan(&totalTickets).Error; err != nil {
			return err
		}

		userPlatformFee := float64(totalTickets * 30)
		baseTicketRevenue := math.Round((normalizedAmount-userPlatformFee)*100) / 100

		// Platform Wallet
		var platformWallet PlatformWalletModel
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", domain.PlatformWalletID).First(&platformWallet).Error; err != nil {
			return err
		}
		if err := tx.Create(&PlatformWalletTransactionModel{
			ID:            uuid.NewString(),
			WalletID:      domain.PlatformWalletID,
			Type:          domain.WalletTransactionTypeCredit,
			Amount:        userPlatformFee,
			ReferenceType: domain.PlatformRefTypePayment,
			ReferenceID:   bookingID,
			CreatedAt:     now,
		}).Error; err != nil {
			return err
		}
		platformWallet.AvailableBalance = math.Round((platformWallet.AvailableBalance+userPlatformFee)*100) / 100
		platformWallet.TotalCredited = math.Round((platformWallet.TotalCredited+userPlatformFee)*100) / 100
		platformWallet.UpdatedAt = now
		if err := tx.Model(&PlatformWalletModel{}).Where("id = ?", domain.PlatformWalletID).Updates(map[string]interface{}{
			"available_balance": platformWallet.AvailableBalance,
			"total_credited":    platformWallet.TotalCredited,
			"updated_at":        platformWallet.UpdatedAt,
		}).Error; err != nil {
			return err
		}

		// Organizer Wallet
		var event EventModel
		if err := tx.Where("id = ?", booking.EventID).First(&event).Error; err != nil {
			return err
		}

		var orgWallet WalletModel
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", event.OrganizerID).First(&orgWallet).Error; err != nil {
			return err
		}
		if err := tx.Create(&WalletTransactionModel{
			ID:            uuid.NewString(),
			WalletID:      orgWallet.ID,
			Type:          domain.WalletTransactionTypeCredit,
			Amount:        baseTicketRevenue,
			ReferenceType: domain.WalletReferenceTypeEarning,
			ReferenceID:   bookingID,
			Status:        domain.WalletTransactionStatusCompleted,
			CreatedAt:     now,
		}).Error; err != nil {
			return err
		}
		orgWallet.PendingBalance = math.Round((orgWallet.PendingBalance+baseTicketRevenue)*100) / 100
		orgWallet.ReserveBalance = math.Round((orgWallet.ReserveBalance-userPlatformFee)*100) / 100
		orgWallet.TotalCredited = math.Round((orgWallet.TotalCredited+baseTicketRevenue)*100) / 100
		orgWallet.UpdatedAt = now
		if err := tx.Model(&WalletModel{}).Where("id = ?", orgWallet.ID).Updates(map[string]interface{}{
			"pending_balance": orgWallet.PendingBalance,
			"reserve_balance": orgWallet.ReserveBalance,
			"total_credited":  orgWallet.TotalCredited,
			"updated_at":      orgWallet.UpdatedAt,
		}).Error; err != nil {
			return err
		}

		return nil
	})
}

