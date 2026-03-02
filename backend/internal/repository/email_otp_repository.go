package repository

import "github.com/aswinsreeraj/evntx/internal/domain"

type EmailOTPRepository interface {
	Create(otp *domain.EmailOTP) error
	InvalidatePrevious(email string) error
	FindValidOTP(email string) (*domain.EmailOTP, error)
	MarkConsumed(id string) error
}
