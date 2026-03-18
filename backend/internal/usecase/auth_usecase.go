package usecase

import (
	"errors"
	"fmt"
	"os"
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
	emailSender repository.EmailSender
	roleRepo    repository.UserRoleRepository
}

func NewAuthUsecase(
	otpRepo repository.EmailOTPRepository,
	userRepo repository.UserRepository,
	sessionRepo repository.UserSessionRepository,
	emailSender repository.EmailSender,
	roleRepo repository.UserRoleRepository,
) *AuthUsecase {
	return &AuthUsecase{
		otpRepo:     otpRepo,
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		emailSender: emailSender,
		roleRepo:    roleRepo,
	}
}

func getAdminEmail() string {
	email := os.Getenv("ADMIN_EMAIL")
	if email == "" {
		return "admin@evntx.com"
	}
	return email
}

func (u *AuthUsecase) RequestEmailOTP(email string) (bool, error) {

	isNewUser := false
	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			isNewUser = true
		} else {
			return false, err
		}
	} else if !user.IsActive {
		return false, fmt.Errorf("Account has been blocked. Please send a mail to admin at %s", getAdminEmail())

	}

	rawOTP, err := otp.GenerateOTP()
	if err != nil {
		return false, err
	}

	hash, err := otp.HashOTP(rawOTP)
	if err != nil {
		return false, err
	}

	if err := u.otpRepo.InvalidatePrevious(email); err != nil {
		return false, err
	}

	emailOTP := &domain.EmailOTP{
		ID:        uuid.NewString(),
		Email:     email,
		OTPHash:   hash,
		ExpiresAt: time.Now().Add(5 * time.Minute),
		Consumed:  false,
	}

	if err := u.otpRepo.Create(emailOTP); err != nil {
		return false, err
	}

	if err := u.emailSender.SendOTP(email, rawOTP); err != nil {
		return false, err
	}

	return isNewUser, nil
}

func (u *AuthUsecase) VerifyEmailOTP(email, rawOTP, name, userAgent, ip string) (*domain.User, []domain.UserRole, string, string, error) {

	logger.Log.Info().Msgf("Verifying email: %s", email)

	storedOTP, err := u.otpRepo.FindValidOTP(email)
	if err != nil {
		logger.Log.Warn().Msgf("FindValidOTP failed: %v", err)
		return nil, nil, "", "", err
	}

	logger.Log.Info().Msgf("Stored OTP hash: %s", storedOTP.OTPHash)
	logger.Log.Info().Msgf("Raw OTP received: %s", rawOTP)

	if err := otp.CompareOTP(storedOTP.OTPHash, rawOTP); err != nil {
		logger.Log.Warn().Msgf("Compare failed: %v", err)
		return nil, nil, "", "", err
	}

	if err := u.otpRepo.MarkConsumed(storedOTP.ID); err != nil {
		return nil, nil, "", "", err
	}

	user, err := u.userRepo.FindByEmail(email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user = &domain.User{
				ID:            uuid.NewString(),
				Email:         email,
				Name:          name,
				IsActive:      true,
				EmailVerified: true,
			}

			if err := u.userRepo.Create(user); err != nil {
				return nil, nil, "", "", err
			}
		} else {
			return nil, nil, "", "", err
		}
	} else if !user.IsActive {
		return nil, nil, "", "", fmt.Errorf("Account has been blocked. Please send a mail to admin at %s", getAdminEmail())
	}

	accessToken, err := jwtutil.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, nil, "", "", err
	}

	refreshToken, err := jwtutil.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, nil, "", "", err
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
		return nil, nil, "", "", err
	}

	roles, err := u.roleRepo.GetRolesByUserID(user.ID)
	if err != nil {
		roles = []domain.UserRole{}
	}

	return user, roles, accessToken, refreshToken, nil
}

func (u *AuthUsecase) Register(email, rawOTP, name, dob, gender, roleStr, organizationName, userAgent, ip string) (*domain.User, []domain.UserRole, string, string, error) {

	logger.Log.Info().Msgf("Registering user: %s", email)

	storedOTP, err := u.otpRepo.FindValidOTP(email)
	if err != nil {
		logger.Log.Warn().Msgf("FindValidOTP failed: %v", err)
		return nil, nil, "", "", err
	}

	if err := otp.CompareOTP(storedOTP.OTPHash, rawOTP); err != nil {
		logger.Log.Warn().Msgf("Compare failed: %v", err)
		return nil, nil, "", "", err
	}

	if err := u.otpRepo.MarkConsumed(storedOTP.ID); err != nil {
		return nil, nil, "", "", err
	}

	user, err := u.userRepo.FindByEmail(email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user = &domain.User{
				ID:               uuid.NewString(),
				Email:            email,
				Name:             name,
				Dob:              dob,
				Gender:           gender,
				OrganizationName: organizationName,
				IsActive:         true,
				EmailVerified:    true,
			}

			if err := u.userRepo.Create(user); err != nil {
				return nil, nil, "", "", err
			}
		} else {
			return nil, nil, "", "", err
		}
	} else {
		if !user.IsActive {
			return nil, nil, "", "", fmt.Errorf("Account has been blocked. Please send a mail to admin at %s", getAdminEmail())
		}
		user.Name = name
		user.Dob = dob
		user.Gender = gender
		user.OrganizationName = organizationName
		if err := u.userRepo.Update(user); err != nil {
			return nil, nil, "", "", err
		}
	}

	accessToken, err := jwtutil.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, nil, "", "", err
	}

	refreshToken, err := jwtutil.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, nil, "", "", err
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
		return nil, nil, "", "", err
	}

	roles, err := u.roleRepo.GetRolesByUserID(user.ID)
	if err != nil {
		roles = []domain.UserRole{}
	}

	targetRole := domain.RoleGoer
	if roleStr == "organizer" {
		targetRole = domain.RoleOrganizer
	}

	hasRole := false
	for _, r := range roles {
		if r == targetRole {
			hasRole = true
			break
		}
	}

	if !hasRole {
		if err := u.roleRepo.AddRole(user.ID, targetRole); err == nil {
			roles = append(roles, targetRole)
		}
	}

	return user, roles, accessToken, refreshToken, nil
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
	} else if !user.IsActive {
		return "", "", fmt.Errorf("Account has been blocked. Please send a mail to admin at %s", getAdminEmail())
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
