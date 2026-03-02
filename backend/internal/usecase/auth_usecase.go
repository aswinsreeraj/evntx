package usecase

import (
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/internal/repository"
	"github.com/aswinsreeraj/evntx/pkg/otp"
	"github.com/google/uuid"
)

type AuthUsecase struct {
	otpRepo repository.EmailOTPRepository
}

func NewAuthUsecase(repo repository.EmailOTPRepository) *AuthUsecase {
	return &AuthUsecase{otpRepo: repo}
}

func (u *AuthUsecase) RequestEmailOTP(email string) (string, error) {
	// generate OTP
	rawOTP, err := otp.GenerateOTP()
	if err != nil {
		return "", err
	}

	hash, err := otp.HashOTP(rawOTP)
	if err != nil {
		return "", err
	}

	// invalidate previous unused OTPs
	if err := u.otpRepo.InvalidatePrevious(email); err != nil {
		return "", err
	}

	emailOTP := &domain.EmailOTP{
		ID:        uuid.NewString(),
		Email:     email,
		OTPHash:   hash,
		ExpiresAt: time.Now().Add(5 * time.Minute),
		Consumed:  false,
	}

	if err := u.otpRepo.Create(emailOTP); err != nil {
		return "", err
	}

	// return raw OTP (for now, for testing)
	return rawOTP, nil
}
