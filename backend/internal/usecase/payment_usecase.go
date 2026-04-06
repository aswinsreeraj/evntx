package usecase

import (
	"context"
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/internal/repository"
	apiErrors "github.com/aswinsreeraj/evntx/pkg/errors"
	"github.com/aswinsreeraj/evntx/pkg/logger"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentUsecase struct {
	bookingRepo         repository.BookingRepository
	eventRepo           repository.EventRepository
	paymentRepo         repository.PaymentRepository
	razorpayService     repository.RazorpayService
	notificationUsecase *NotificationUsecase
	engagementRepo      repository.EngagementRepository
}

func NewPaymentUsecase(
	bookingRepo repository.BookingRepository,
	eventRepo repository.EventRepository,
	paymentRepo repository.PaymentRepository,
	razorpayService repository.RazorpayService,
	notificationUsecase *NotificationUsecase,
	engagementRepo repository.EngagementRepository,
) *PaymentUsecase {
	return &PaymentUsecase{
		bookingRepo:         bookingRepo,
		eventRepo:           eventRepo,
		paymentRepo:         paymentRepo,
		razorpayService:     razorpayService,
		notificationUsecase: notificationUsecase,
		engagementRepo:      engagementRepo,
	}
}

func (u *PaymentUsecase) CreatePaymentOrder(ctx context.Context, bookingID string, userID string) (*domain.PaymentOrderResponse, error) {
	booking, err := u.bookingRepo.FindByID(ctx, bookingID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apiErrors.ErrResourceNotFound
		}

		return nil, err
	}

	if booking.UserID != userID {
		return nil, apiErrors.ErrForbiddenAction
	}

	if booking.Status != "reserved" {
		return nil, apiErrors.ErrInvalidStateTransition
	}

	if booking.TotalAmount == 0 {
		payment := &domain.Payment{
			ID:                uuid.NewString(),
			BookingID:         booking.ID,
			Provider:          "free",
			ProviderReference: "FREE_" + booking.ID,
			Amount:            0,
			Status:            domain.PaymentStatusSuccess,
			RawResponse:       []byte(`{}`),
			CreatedAt:         time.Now(),
		}

		if err := u.paymentRepo.CreatePayment(payment); err != nil {
			return nil, err
		}

		if err := u.paymentRepo.MarkPaymentSuccess(payment.ID, booking.ID, booking.EventID, 0); err != nil {
			return nil, err
		}

		if u.engagementRepo != nil {
			_ = u.engagementRepo.IncrementSuccessfulBookings(ctx, booking.EventID)
		}

		logger.Log.Info().
			Str("booking_id", booking.ID).
			Str("user_id", userID).
			Str("payment_id", payment.ID).
			Msg("free_payment_order_created_and_successful")

		return &domain.PaymentOrderResponse{
			OrderID:       "FREE_" + booking.ID,
			Amount:        0,
			Currency:      "INR",
			RazorpayKey:   "",
			IsFreeBooking: true,
		}, nil
	}

	amountInPaise := int64(math.Round(booking.TotalAmount * 100))
	order, err := u.razorpayService.CreateOrder(amountInPaise, booking.ID)
	if err != nil {
		return nil, apiErrors.Wrap(err, 500, apiErrors.PaymentFailed, "Failed to create payment order")
	}

	payment := &domain.Payment{
		ID:                uuid.NewString(),
		BookingID:         booking.ID,
		Provider:          "razorpay",
		ProviderReference: order.ID,
		Amount:            booking.TotalAmount,
		Status:            domain.PaymentStatusInitiated,
		RawResponse:       order.RawResponse,
		CreatedAt:         time.Now(),
	}

	if err := u.paymentRepo.CreatePayment(payment); err != nil {
		return nil, err
	}

	logger.Log.Info().
		Str("booking_id", booking.ID).
		Str("user_id", userID).
		Str("payment_id", payment.ID).
		Str("provider_reference", order.ID).
		Int64("amount", order.Amount).
		Msg("payment_order_created")

	return &domain.PaymentOrderResponse{
		OrderID:     order.ID,
		Amount:      order.Amount,
		Currency:    order.Currency,
		RazorpayKey: u.razorpayService.GetKeyID(),
	}, nil
}

func (u *PaymentUsecase) VerifyPayment(
	ctx context.Context,
	razorpayOrderID string,
	razorpayPaymentID string,
	razorpaySignature string,
) error {
	payment, err := u.paymentRepo.FindByProviderReference(razorpayOrderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apiErrors.ErrResourceNotFound
		}
		return err
	}

	if payment.Status == domain.PaymentStatusSuccess {
		return nil
	}

	isValid, err := u.razorpayService.VerifySignature(razorpayOrderID, razorpayPaymentID, razorpaySignature)
	if err != nil {
		return apiErrors.Wrap(err, 500, apiErrors.PaymentFailed, "Failed to verify payment signature")
	}

	if !isValid {
		if updateErr := u.paymentRepo.UpdateStatus(payment.ID, domain.PaymentStatusFailed); updateErr != nil {
			return updateErr
		}

		if u.notificationUsecase != nil {
			if booking, bookingErr := u.bookingRepo.FindByID(ctx, payment.BookingID); bookingErr == nil {
				if notifyErr := u.notificationUsecase.SendNotification(
					booking.UserID,
					domain.NotificationTypePaymentFailed,
					"Payment failed",
					"Payment failed. Please retry.",
					map[string]interface{}{
						"booking_id": payment.BookingID,
						"payment_id": payment.ID,
					},
				); notifyErr != nil {
					logger.Log.Warn().
						Err(notifyErr).
						Str("user_id", booking.UserID).
						Str("payment_id", payment.ID).
						Msg("notification_send_failed")
				}
			}
		}

		logger.Log.Warn().
			Str("payment_id", payment.ID).
			Str("booking_id", payment.BookingID).
			Str("provider_reference", razorpayOrderID).
			Str("razorpay_payment_id", razorpayPaymentID).
			Msg("payment_signature_invalid")
		return apiErrors.New(400, apiErrors.PaymentFailed, "Invalid payment signature")
	}

	booking, err := u.bookingRepo.FindByID(context.Background(), payment.BookingID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apiErrors.ErrResourceNotFound
		}
		return err
	}

	event, err := u.eventRepo.GetEventByID(booking.EventID)
	if err != nil {
		return err
	}

	if err := u.paymentRepo.MarkPaymentSuccess(payment.ID, payment.BookingID, event.OrganizerID, payment.Amount); err != nil {
		if errors.Is(err, apiErrors.ErrBookingExpiredPaymentSuccess) {
			if lateErr := u.paymentRepo.HandleLatePayment(payment.ID, payment.BookingID, booking.UserID, payment.Amount); lateErr != nil {
				return lateErr
			}
			if u.notificationUsecase != nil {
				u.notificationUsecase.SendNotification(
					booking.UserID,
					domain.NotificationTypePaymentSuccess,
					"Payment received but booking expired",
					"Your payment was successful but the booking expired. The amount will be refunded to your given payment details in 3-5 working days. Please update your payment details in your profile if you haven't.",
					map[string]interface{}{
						"booking_id": payment.BookingID,
						"payment_id": payment.ID,
						"is_late_payment": true,
					},
				)
			}
			logger.Log.Info().
				Str("payment_id", payment.ID).
				Str("booking_id", payment.BookingID).
				Msg("payment_success_but_booking_expired_handled")
			return apiErrors.ErrBookingExpiredPaymentSuccess
		}
		return err
	}

	if u.engagementRepo != nil {
		_ = u.engagementRepo.IncrementSuccessfulBookings(ctx, booking.EventID)
	}

	if u.notificationUsecase != nil {
		booking, bookingErr := u.bookingRepo.FindByID(ctx, payment.BookingID)
		if bookingErr != nil {
			logger.Log.Warn().
				Err(bookingErr).
				Str("payment_id", payment.ID).
				Msg("notification_context_lookup_failed")
		} else if event, eventErr := u.eventRepo.GetEventByID(booking.EventID); eventErr != nil {
			logger.Log.Warn().
				Err(eventErr).
				Str("payment_id", payment.ID).
				Msg("notification_context_lookup_failed")
		} else {
			if notifyErr := u.notificationUsecase.SendNotification(
				booking.UserID,
				domain.NotificationTypePaymentSuccess,
				"Payment successful",
				"Payment successful. Tickets confirmed.",
				map[string]interface{}{
					"booking_id":  booking.ID,
					"event_id":    event.ID,
					"event_title": event.Title,
					"amount":      payment.Amount,
				},
			); notifyErr != nil {
				logger.Log.Warn().
					Err(notifyErr).
					Str("user_id", booking.UserID).
					Str("payment_id", payment.ID).
					Msg("notification_send_failed")
			}

			if notifyErr := u.notificationUsecase.SendNotification(
				booking.UserID,
				domain.NotificationTypeTicketGenerated,
				"Your tickets are generated",
				"Your tickets are generated",
				map[string]interface{}{
					"booking_id":  booking.ID,
					"event_id":    event.ID,
					"event_title": event.Title,
				},
			); notifyErr != nil {
				logger.Log.Warn().
					Err(notifyErr).
					Str("user_id", booking.UserID).
					Str("payment_id", payment.ID).
					Msg("notification_send_failed")
			}

			ticketsInBooking, _ := u.bookingRepo.GetTicketCountByBookingID(ctx, booking.ID)
			organizerEarnings := payment.Amount - 30*float64(ticketsInBooking)
			if organizerEarnings < 0 {
				organizerEarnings = 0
			}

			if notifyErr := u.notificationUsecase.SendNotification(
				event.OrganizerID,
				domain.NotificationTypePaymentSuccess,
				"New booking received",
				"New booking received. You will earn ₹"+strconv.FormatFloat(organizerEarnings, 'f', 2, 64)+" after settlement.",
				map[string]interface{}{
					"booking_id":  booking.ID,
					"event_id":    event.ID,
					"event_title": event.Title,
					"amount":      organizerEarnings,
				},
			); notifyErr != nil {
				logger.Log.Warn().
					Err(notifyErr).
					Str("organizer_id", event.OrganizerID).
					Str("payment_id", payment.ID).
					Msg("notification_send_failed")
			}
		}
	}

	logger.Log.Info().
		Str("payment_id", payment.ID).
		Str("booking_id", payment.BookingID).
		Str("provider_reference", razorpayOrderID).
		Str("razorpay_payment_id", razorpayPaymentID).
		Msg("payment_verified")

	return nil
}

func (u *PaymentUsecase) RefundPaymentToWallet(ctx context.Context, bookingID string, userID string) error {
	booking, err := u.bookingRepo.FindByID(ctx, bookingID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apiErrors.ErrResourceNotFound
		}

		return err
	}

	if booking.UserID != userID {
		return apiErrors.ErrForbiddenAction
	}

	if booking.Status != "paid" && booking.Status != "cancelled" {
		return apiErrors.ErrInvalidStateTransition
	}

	payment, err := u.paymentRepo.FindByBookingID(bookingID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apiErrors.ErrResourceNotFound
		}

		return err
	}

	if payment.Status == domain.PaymentStatusRefunded {
		return apiErrors.New(409, apiErrors.DuplicateResource, "Refund already processed")
	}

	if payment.Status != domain.PaymentStatusSuccess {
		return apiErrors.New(400, apiErrors.InvalidStateTransition, "Payment is not eligible for refund")
	}

	refundAmount := normalizeRefundAmount(payment.Amount - booking.TotalAmount)
	if refundAmount <= 0 {
		return apiErrors.New(400, apiErrors.InvalidStateTransition, "No refundable ticket amount available")
	}

	platformFeeAmount := normalizeRefundAmount(payment.Amount - refundAmount)
	if platformFeeAmount < 0 {
		return apiErrors.New(400, apiErrors.InvalidStateTransition, "Invalid refund split")
	}

	if err := u.paymentRepo.RefundPaymentToWallet(
		booking.UserID,
		payment.ID,
		booking.ID,
		refundAmount,
		platformFeeAmount,
	); err != nil {
		return err
	}

	logger.Log.Info().
		Str("payment_id", payment.ID).
		Str("booking_id", booking.ID).
		Str("user_id", booking.UserID).
		Float64("refund_amount", refundAmount).
		Float64("platform_fee_amount", platformFeeAmount).
		Msg("payment_refunded_to_wallet")

	return nil
}

func normalizeRefundAmount(amount float64) float64 {
	return math.Round(amount*100) / 100
}
