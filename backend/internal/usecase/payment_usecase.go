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
	paymentRepo     repository.PaymentRepository
	razorpayService repository.RazorpayService
}

func NewPaymentUsecase(
	bookingRepo repository.BookingRepository,
	paymentRepo repository.PaymentRepository,
	razorpayService repository.RazorpayService,
) *PaymentUsecase {
	return &PaymentUsecase{
		bookingRepo:     bookingRepo,
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
