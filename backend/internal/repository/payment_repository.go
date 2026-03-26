package repository

import "github.com/aswinsreeraj/evntx/internal/domain"

type PaymentRepository interface {
	CreatePayment(payment *domain.Payment) error
	FindByProviderReference(orderID string) (*domain.Payment, error)
	UpdateStatus(paymentID string, status string) error
}
