package usecase

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/aswinsreeraj/evntx/pkg/logger"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/internal/repository"
	apiErrors "github.com/aswinsreeraj/evntx/pkg/errors"
	"github.com/aswinsreeraj/evntx/pkg/hash"
	jwtutil "github.com/aswinsreeraj/evntx/pkg/jwt"
	oauthutil "github.com/aswinsreeraj/evntx/pkg/oauth"
	"github.com/aswinsreeraj/evntx/pkg/otp"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthUsecase struct {
	otpRepo		repository.EmailOTPRepository
	userRepo	repository.UserRepository
	sessionRepo	repository.UserSessionRepository
	emailSender	repository.EmailSender
	roleRepo	repository.UserRoleRepository
	walletRepo	repository.WalletRepository
	settingsRepo	repository.SettingsRepository
}

func NewAuthUsecase(
	otpRepo repository.EmailOTPRepository,
	userRepo repository.UserRepository,
	sessionRepo repository.UserSessionRepository,
	emailSender repository.EmailSender,
	roleRepo repository.UserRoleRepository,
	walletRepo repository.WalletRepository,
	settingsRepo repository.SettingsRepository,
) *AuthUsecase {
	return &AuthUsecase{
		otpRepo:	otpRepo,
		userRepo:	userRepo,
		sessionRepo:	sessionRepo,
		emailSender:	emailSender,
		roleRepo:	roleRepo,
		walletRepo:	walletRepo,
		settingsRepo:	settingsRepo,
	}
}

func getAdminEmail() string {
	email := os.Getenv("ADMIN_EMAIL")
	if email == "" {
		return "admin@evntx.com"
	}
	return email
}

func (u *AuthUsecase) ensureOrganizerNotPending(userID string) error {
	detail, err := u.userRepo.GetOrganizerDetails(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	if detail != nil && detail.ApprovalStatus == "pending" {
		return apiErrors.New(403, apiErrors.ForbiddenAction, "Your organizer account is pending admin approval. You can log in after approval.")
	}
	return nil
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
		return false, apiErrors.New(403, apiErrors.ForbiddenAction, fmt.Sprintf("Account has been blocked. Please send a mail to admin at %s", getAdminEmail()))

	}

	if isNewUser && u.settingsRepo != nil {
		if settings, settingsErr := u.settingsRepo.GetPlatformSettings(); settingsErr == nil && !settings.EnableUserRegistration {
			return false, apiErrors.New(403, apiErrors.ForbiddenAction, "New registrations are currently disabled by admin")
		}
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
		ID:		uuid.NewString(),
		Email:		email,
		OTPHash:	hash,
		ExpiresAt:	time.Now().Add(5 * time.Minute),
		Consumed:	false,
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

	storedOTP, err := u.otpRepo.FindValidOTP(email)
	if err != nil {
		logger.Log.Warn().Err(err).Str("email", email).Msg("OTP lookup failed")
		return nil, nil, "", "", err
	}

	if err := otp.CompareOTP(storedOTP.OTPHash, rawOTP); err != nil {
		logger.Log.Warn().Str("email", email).Msg("OTP verification failed: invalid code")
		return nil, nil, "", "", err
	}

	if err := u.otpRepo.MarkConsumed(storedOTP.ID); err != nil {
		return nil, nil, "", "", err
	}

	user, err := u.userRepo.FindByEmail(email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user = &domain.User{
				ID:		uuid.NewString(),
				Email:		email,
				Name:		name,
				IsActive:	true,
				EmailVerified:	true,
			}

			if err := u.userRepo.Create(user); err != nil {
				return nil, nil, "", "", err
			}

			wallet := &domain.Wallet{
				ID:			uuid.NewString(),
				UserID:			user.ID,
				AvailableBalance:	0,
				PendingBalance:		0,
				TotalCredited:		0,
				TotalDebited:		0,
				UpdatedAt:		time.Now(),
			}
			if err := u.walletRepo.CreateWallet(wallet); err != nil {
				logger.Log.Error().Err(err).Msg("failed to create wallet for user during authentication/registration")
			}
		} else {
			return nil, nil, "", "", err
		}
	} else if !user.IsActive {
		return nil, nil, "", "", apiErrors.New(403, apiErrors.ForbiddenAction, fmt.Sprintf("Account has been blocked. Please send a mail to admin at %s", getAdminEmail()))
	}

	if pendingErr := u.ensureOrganizerNotPending(user.ID); pendingErr != nil {
		return nil, nil, "", "", pendingErr
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
		ID:			uuid.NewString(),
		UserID:			user.ID,
		RefreshTokenHash:	refreshHash,
		UserAgent:		userAgent,
		IPAddress:		ip,
		ExpiresAt:		time.Now().Add(7 * 24 * time.Hour),
		Revoked:		false,
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
	if u.settingsRepo != nil {
		if settings, settingsErr := u.settingsRepo.GetPlatformSettings(); settingsErr == nil {
			if !settings.EnableUserRegistration {
				return nil, nil, "", "", apiErrors.New(403, apiErrors.ForbiddenAction, "New registrations are currently disabled by admin")
			}
		}
	}

	storedOTP, err := u.otpRepo.FindValidOTP(email)
	if err != nil {
		logger.Log.Warn().Err(err).Str("email", email).Msg("OTP lookup failed during registration")
		return nil, nil, "", "", err
	}

	if err := otp.CompareOTP(storedOTP.OTPHash, rawOTP); err != nil {
		logger.Log.Warn().Str("email", email).Msg("OTP verification failed during registration")
		return nil, nil, "", "", err
	}

	if err := u.otpRepo.MarkConsumed(storedOTP.ID); err != nil {
		return nil, nil, "", "", err
	}

	user, err := u.userRepo.FindByEmail(email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user = &domain.User{
				ID:		uuid.NewString(),
				Email:		email,
				Name:		name,
				Dob:		dob,
				Gender:		gender,
				IsActive:	true,
				EmailVerified:	true,
			}

			if err := u.userRepo.Create(user); err != nil {
				return nil, nil, "", "", err
			}
		} else {
			return nil, nil, "", "", err
		}
	} else {
		if !user.IsActive {
			return nil, nil, "", "", apiErrors.New(403, apiErrors.ForbiddenAction, fmt.Sprintf("Account has been blocked. Please send a mail to admin at %s", getAdminEmail()))
		}
		user.Name = name
		user.Dob = dob
		user.Gender = gender
		if err := u.userRepo.Update(user); err != nil {
			return nil, nil, "", "", err
		}
	}

	roles, err := u.roleRepo.GetRolesByUserID(user.ID)
	if err != nil {
		roles = []domain.UserRole{}
	}

	requiresApproval := false
	if roleStr == "organizer" && u.settingsRepo != nil {
		if settings, settingsErr := u.settingsRepo.GetPlatformSettings(); settingsErr == nil {
			requiresApproval = settings.RequireAdminApprovalForOrganizers
		}
	}

	if roleStr == "organizer" {
		hasOrganizer := false
		for _, r := range roles {
			if r == domain.RoleOrganizer {
				hasOrganizer = true
				break
			}
		}
		if !hasOrganizer && !requiresApproval {
			if err := u.roleRepo.AddRole(user.ID, domain.RoleOrganizer); err == nil {
				roles = append(roles, domain.RoleOrganizer)
			}
		}
	}

	if roleStr == "organizer" && organizationName != "" {
		approvalStatus := "approved"
		if requiresApproval {
			approvalStatus = "pending"
		}
		if err := u.userRepo.UpsertOrganizerDetails(&domain.OrganizerDetail{
			UserID:			user.ID,
			OrganizationName:	organizationName,
			Address:		"",
			ApprovalStatus:		approvalStatus,
		}); err != nil {
			return nil, nil, "", "", err
		}
	}

	if roleStr == "organizer" && requiresApproval {
		return nil, nil, "", "", apiErrors.New(403, apiErrors.ForbiddenAction, "Organizer registration submitted and pending admin approval. You can log in once approved.")
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
		ID:			uuid.NewString(),
		UserID:			user.ID,
		RefreshTokenHash:	refreshHash,
		UserAgent:		userAgent,
		IPAddress:		ip,
		ExpiresAt:		time.Now().Add(7 * 24 * time.Hour),
		Revoked:		false,
	}

	if err := u.sessionRepo.Create(session); err != nil {
		return nil, nil, "", "", err
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
	if u.settingsRepo != nil {
		if settings, settingsErr := u.settingsRepo.GetPlatformSettings(); settingsErr == nil {
			if !settings.AllowGoogleLogin {
				return "", "", apiErrors.New(403, apiErrors.ForbiddenAction, "Google login is currently disabled by admin")
			}
		}
	}

	googleUser, err := oauthutil.VerifyGoogleIDToken(idToken)
	if err != nil {
		return "", "", err
	}

	user, err := u.userRepo.FindByEmail(googleUser.Email)
	if err != nil {
		if u.settingsRepo != nil {
			if settings, settingsErr := u.settingsRepo.GetPlatformSettings(); settingsErr == nil && !settings.EnableUserRegistration {
				return "", "", apiErrors.New(403, apiErrors.ForbiddenAction, "New registrations are currently disabled by admin")
			}
		}

		user = &domain.User{
			ID:		uuid.NewString(),
			Email:		googleUser.Email,
			Name:		googleUser.Name,
			IsActive:	true,
			EmailVerified:	true,
		}
		if err := u.userRepo.Create(user); err != nil {
			return "", "", err
		}

		wallet := &domain.Wallet{
			ID:			uuid.NewString(),
			UserID:			user.ID,
			AvailableBalance:	0,
			PendingBalance:		0,
			TotalCredited:		0,
			TotalDebited:		0,
			UpdatedAt:		time.Now(),
		}
		if err := u.walletRepo.CreateWallet(wallet); err != nil {
			logger.Log.Error().Err(err).Msg("failed to create wallet for user during Google login")
		}
	} else if !user.IsActive {
		return "", "", apiErrors.New(403, apiErrors.ForbiddenAction, fmt.Sprintf("Account has been blocked. Please send a mail to admin at %s", getAdminEmail()))
	}

	if pendingErr := u.ensureOrganizerNotPending(user.ID); pendingErr != nil {
		return "", "", pendingErr
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
		ID:			uuid.NewString(),
		UserID:			user.ID,
		RefreshTokenHash:	refreshHash,
		UserAgent:		userAgent,
		IPAddress:		ip,
		ExpiresAt:		time.Now().Add(7 * 24 * time.Hour),
		Revoked:		false,
	}

	if err := u.sessionRepo.Create(session); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
