package repository

import "github.com/aswinsreeraj/evntx/internal/domain"

type RazorpayService interface {
	GetKeyID() string
	CreateOrder(amount int64, receipt string) (*domain.RazorpayOrder, error)
	VerifySignature(orderID string, paymentID string, signature string) (bool, error)
	FetchOrder(orderID string) (*domain.RazorpayOrder, error)
	RefundPayment(paymentID string, amount int64) error
}
