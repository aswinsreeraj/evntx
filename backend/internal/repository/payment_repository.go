package repository

import "github.com/aswinsreeraj/evntx/internal/domain"

type PaymentRepository interface {
	CreatePayment(payment *domain.Payment) error
	FindByProviderReference(orderID string) (*domain.Payment, error)
	FindByBookingID(bookingID string) (*domain.Payment, error)
	UpdateStatus(paymentID string, status string) error
	MarkPaymentSuccess(paymentID string, bookingID string, organizerID string, amount float64) error
	RefundPaymentToWallet(userID string, paymentID string, bookingID string, refundAmount float64, platformFeeAmount float64) error
}
