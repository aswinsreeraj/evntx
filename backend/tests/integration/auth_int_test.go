package integration

import (
	"testing"

	repoImpl "github.com/aswinsreeraj/evntx/internal/infrastructure/repository"
	"github.com/aswinsreeraj/evntx/internal/usecase"
	"github.com/aswinsreeraj/evntx/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

type CapturingEmailSender struct {
	LastEmail string
	LastOTP   string
}

func (s *CapturingEmailSender) SendOTP(to, otp string) error {
	s.LastEmail = to
	s.LastOTP = otp
	return nil
}

func (s *CapturingEmailSender) SendOrganizerApproval(email, name string) error {
	return nil
}

func TestAuthIntegration_OTPFlow(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.ClearDatabase(db)

	emailSender := &CapturingEmailSender{}

	otpRepo := repoImpl.NewEmailOTPGormRepository(db)
	userRepo := repoImpl.NewUserGormRepository(db)
	sessionRepo := repoImpl.NewUserSessionGormRepository(db)
	roleRepo := repoImpl.NewUserRoleGormRepository(db)
	walletRepo := repoImpl.NewWalletGormRepository(db)
	settingsRepo := repoImpl.NewSettingsGormRepository(db)
	_ = settingsRepo.EnsureExists()

	authUsecase := usecase.NewAuthUsecase(otpRepo, userRepo, sessionRepo, emailSender, roleRepo, walletRepo, settingsRepo)

	testEmail := "test_auth_flow@example.com"

	
	isNewUser, err := authUsecase.RequestEmailOTP(testEmail)
	assert.NoError(t, err)
	assert.True(t, isNewUser, "User should be identified as new")
	assert.Equal(t, testEmail, emailSender.LastEmail)
	assert.NotEmpty(t, emailSender.LastOTP)

	rawOTP := emailSender.LastOTP

	
	storedOTP, err := otpRepo.FindValidOTP(testEmail)
	assert.NoError(t, err)
	assert.NotNil(t, storedOTP)
	assert.Equal(t, testEmail, storedOTP.Email)
	assert.False(t, storedOTP.Consumed)

	
	user, roles, accessToken, refreshToken, err := authUsecase.VerifyEmailOTP(testEmail, rawOTP, "Test User", "Go-http-client/1.1", "127.0.0.1")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, testEmail, user.Email)
	assert.Equal(t, "Test User", user.Name)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)
	assert.Empty(t, roles, "New user should have no roles by default")

	
	dbUser, err := userRepo.FindByEmail(testEmail)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, dbUser.ID)
	assert.True(t, dbUser.EmailVerified)

	
	wallet, err := walletRepo.GetWalletByUserID(user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, wallet)
	assert.Equal(t, float64(0), wallet.AvailableBalance)

	
	consumedOTP, err := otpRepo.FindValidOTP(testEmail)
	assert.Error(t, err) 
	assert.Nil(t, consumedOTP)
}
