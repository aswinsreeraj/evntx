package usecase

import (
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/internal/repository"
	jwtutil "github.com/aswinsreeraj/evntx/pkg/jwt"
	"github.com/aswinsreeraj/evntx/pkg/otp"
	"github.com/google/uuid"
)

type AuthUsecase struct {
	otpRepo     repository.EmailOTPRepository
	userRepo    repository.UserRepository
	sessionRepo repository.UserSessionRepository
}

func NewAuthUsecase(
	otpRepo repository.EmailOTPRepository,
	userRepo repository.UserRepository,
	sessionRepo repository.UserSessionRepository,
) *AuthUsecase {
	return &AuthUsecase{
		otpRepo:     otpRepo,
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
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

func (u *AuthUsecase) VerifyEmailOTP(email, rawOTP, userAgent, ip string) (string, string, error) {

	storedOTP, err := u.otpRepo.FindValidOTP(email)
	if err != nil {
		return "", "", err
	}

	if err := otp.CompareOTP(storedOTP.OTPHash, rawOTP); err != nil {
		return "", "", err
	}

	if err := u.otpRepo.MarkConsumed(storedOTP.ID); err != nil {
		return "", "", err
	}

	// find or create user
	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		user = &domain.User{
			ID:            uuid.NewString(),
			Email:         email,
			IsActive:      true,
			EmailVerified: true,
		}
		if err := u.userRepo.Create(user); err != nil {
			return "", "", err
		}
	}

	accessToken, err := jwtutil.GenerateAccessToken(user.ID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwtutil.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", err
	}

	refreshHash, err := otp.HashOTP(refreshToken)
	if err != nil {
		return "", "", err
	}

	session := &domain.UserSession{
		ID:               uuid.NewString(),
		UserID:           user.ID,
		RefreshTokenHash: refreshHash,
		UserAgent:        userAgent,
		IPAddress:        ip,
		ExpiresAt:        time.Now().Add(7 * 24 * time.Hour),
		Revoked:          false,
	}

	if err := u.sessionRepo.Create(session); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
