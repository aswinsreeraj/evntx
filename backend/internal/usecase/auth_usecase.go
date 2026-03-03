package usecase

import (
	"errors"
	"time"

	"github.com/aswinsreeraj/evntx/pkg/logger"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/internal/repository"
	"github.com/aswinsreeraj/evntx/pkg/hash"
	jwtutil "github.com/aswinsreeraj/evntx/pkg/jwt"
	oauthutil "github.com/aswinsreeraj/evntx/pkg/oauth"
	"github.com/aswinsreeraj/evntx/pkg/otp"
	"github.com/google/uuid"
	"gorm.io/gorm"
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

	logger.Log.Info().Msgf("Verifying email: %s", email)

	storedOTP, err := u.otpRepo.FindValidOTP(email)
	if err != nil {
		logger.Log.Warn().Msgf("FindValidOTP failed: %v", err)
		return "", "", err
	}

	logger.Log.Info().Msgf("Stored OTP hash: %s", storedOTP.OTPHash)
	logger.Log.Info().Msgf("Raw OTP received: %s", rawOTP)

	if err := otp.CompareOTP(storedOTP.OTPHash, rawOTP); err != nil {
		logger.Log.Warn().Msgf("Compare failed: %v", err)
		return "", "", err
	}

	if err := u.otpRepo.MarkConsumed(storedOTP.ID); err != nil {
		return "", "", err
	}

	// find or create user
	user, err := u.userRepo.FindByEmail(email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user = &domain.User{
				ID:            uuid.NewString(),
				Email:         email,
				IsActive:      true,
				EmailVerified: true,
			}

			if err := u.userRepo.Create(user); err != nil {
				return "", "", err
			}
		} else {
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

	refreshHash := hash.HashToken(refreshToken)

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

func (u *AuthUsecase) RefreshToken(refreshToken string) (string, error) {

	userID, err := jwtutil.ParseRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}

	session, err := u.sessionRepo.FindByUserID(userID)
	if err != nil {
		return "", err
	}

	if session.Revoked || session.ExpiresAt.Before(time.Now()) {
		return "", err
	}

	if session.RefreshTokenHash != hash.HashToken(refreshToken) {
		return "", errors.New("invalid refresh token")
	}

	return jwtutil.GenerateAccessToken(userID)
}

func (u *AuthUsecase) Logout(refreshToken string) error {

	userID, err := jwtutil.ParseRefreshToken(refreshToken)
	if err != nil {
		return err
	}

	session, err := u.sessionRepo.FindByUserID(userID)
	if err != nil {
		return err
	}

	return u.sessionRepo.Revoke(session.ID)
}

func (u *AuthUsecase) GoogleLogin(idToken, userAgent, ip string) (string, string, error) {

	googleUser, err := oauthutil.VerifyGoogleIDToken(idToken)
	if err != nil {
		return "", "", err
	}

	user, err := u.userRepo.FindByEmail(googleUser.Email)
	if err != nil {
		// create new user
		user = &domain.User{
			ID:            uuid.NewString(),
			Email:         googleUser.Email,
			Name:          googleUser.Name,
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

	refreshHash := hash.HashToken(refreshToken)

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
