package usecase

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/internal/repository"
	apiErrors "github.com/aswinsreeraj/evntx/pkg/errors"
	"github.com/aswinsreeraj/evntx/pkg/logger"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentUsecase struct {
	bookingRepo     repository.BookingRepository
	eventRepo       repository.EventRepository
	paymentRepo     repository.PaymentRepository
	razorpayService repository.RazorpayService
}

func NewPaymentUsecase(
	bookingRepo repository.BookingRepository,
	eventRepo repository.EventRepository,
	paymentRepo repository.PaymentRepository,
	razorpayService repository.RazorpayService,
) *PaymentUsecase {
	return &PaymentUsecase{
		bookingRepo:     bookingRepo,
		eventRepo:       eventRepo,
		paymentRepo:     paymentRepo,
		razorpayService: razorpayService,
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

func (u *PaymentUsecase) VerifyPayment(razorpayOrderID string, razorpayPaymentID string, razorpaySignature string) error {
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
		return err
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

	if err := u.paymentRepo.RefundPaymentToWallet(booking.UserID, payment.ID, booking.ID, payment.Amount); err != nil {
		return err
	}

	logger.Log.Info().
		Str("payment_id", payment.ID).
		Str("booking_id", booking.ID).
		Str("user_id", booking.UserID).
		Float64("amount", payment.Amount).
		Msg("payment_refunded_to_wallet")

	return nil
}
