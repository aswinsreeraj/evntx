package repository

import "github.com/aswinsreeraj/evntx/internal/domain"

type RazorpayService interface {
	CreateOrder(amount int64, receipt string) (*domain.RazorpayOrder, error)
	GetKeyID() string
	VerifySignature(orderID string, paymentID string, signature string) (bool, error)
}
