package usecase_test

import (
	"testing"
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/internal/repository/mocks"
	"github.com/aswinsreeraj/evntx/internal/usecase"
	"github.com/aswinsreeraj/evntx/pkg/otp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestAuthUsecase_RequestEmailOTP(t *testing.T) {
	mockOTPRepo := mocks.NewEmailOTPRepository(t)
	mockUserRepo := mocks.NewUserRepository(t)
	mockEmailSender := mocks.NewEmailSender(t)
	mockSettingsRepo := mocks.NewSettingsRepository(t)

	authUsecase := usecase.NewAuthUsecase(
		mockOTPRepo,
		mockUserRepo,
		nil,
		mockEmailSender,
		nil,
		nil,
		mockSettingsRepo,
	)

	t.Run("Success_NewUser", func(t *testing.T) {
		email := "test@example.com"
		mockUserRepo.On("FindByEmail", email).Return(nil, gorm.ErrRecordNotFound).Once()
		mockSettingsRepo.On("GetPlatformSettings").Return(&domain.PlatformSettings{EnableUserRegistration: true}, nil).Once()
		mockOTPRepo.On("InvalidatePrevious", email).Return(nil).Once()
		mockOTPRepo.On("Create", mock.AnythingOfType("*domain.EmailOTP")).Return(nil).Once()
		mockEmailSender.On("SendOTP", email, mock.AnythingOfType("string")).Return(nil).Once()

		isNewUser, err := authUsecase.RequestEmailOTP(email)

		assert.NoError(t, err)
		assert.True(t, isNewUser)
	})

	t.Run("Failure_BlockedUser", func(t *testing.T) {
		email := "blocked@example.com"
		mockUserRepo.On("FindByEmail", email).Return(&domain.User{IsActive: false}, nil).Once()

		isNewUser, err := authUsecase.RequestEmailOTP(email)

		assert.Error(t, err)
		assert.False(t, isNewUser)
		assert.Contains(t, err.Error(), "Account has been blocked")
	})
}

func TestAuthUsecase_VerifyEmailOTP(t *testing.T) {
	mockOTPRepo := mocks.NewEmailOTPRepository(t)
	mockUserRepo := mocks.NewUserRepository(t)
	mockSessionRepo := mocks.NewUserSessionRepository(t)
	mockRoleRepo := mocks.NewUserRoleRepository(t)
	mockWalletRepo := mocks.NewWalletRepository(t)

	authUsecase := usecase.NewAuthUsecase(
		mockOTPRepo,
		mockUserRepo,
		mockSessionRepo,
		nil,
		mockRoleRepo,
		mockWalletRepo,
		nil,
	)

	t.Run("Success_VerifyOTP", func(t *testing.T) {
		email := "test@example.com"
		rawOTP := "123456"
		hashedOTP, _ := otp.HashOTP(rawOTP)

		storedOTP := &domain.EmailOTP{
			ID:		"otp-id",
			Email:		email,
			OTPHash:	hashedOTP,
			ExpiresAt:	time.Now().Add(5 * time.Minute),
			Consumed:	false,
		}

		user := &domain.User{ID: "user-id", Email: email, IsActive: true}

		mockOTPRepo.On("FindValidOTP", email).Return(storedOTP, nil).Once()
		mockOTPRepo.On("MarkConsumed", storedOTP.ID).Return(nil).Once()
		mockUserRepo.On("FindByEmail", email).Return(user, nil).Once()
		mockUserRepo.On("GetOrganizerDetails", user.ID).Return(nil, gorm.ErrRecordNotFound).Once()
		mockSessionRepo.On("Create", mock.AnythingOfType("*domain.UserSession")).Return(nil).Once()
		mockRoleRepo.On("GetRolesByUserID", user.ID).Return([]domain.UserRole{}, nil).Once()

		retUser, roles, accessToken, refreshToken, err := authUsecase.VerifyEmailOTP(email, rawOTP, "Test User", "Go-http-client/1.1", "127.0.0.1")

		assert.NoError(t, err)
		assert.NotNil(t, retUser)
		assert.Empty(t, roles)
		assert.NotEmpty(t, accessToken)
		assert.NotEmpty(t, refreshToken)
	})

	t.Run("Failure_InvalidOTP", func(t *testing.T) {
		email := "test@example.com"
		rawOTP := "123456"
		wrongHash, _ := otp.HashOTP("654321")

		storedOTP := &domain.EmailOTP{
			ID:		"otp-id",
			Email:		email,
			OTPHash:	wrongHash,
			ExpiresAt:	time.Now().Add(5 * time.Minute),
		}

		mockOTPRepo.On("FindValidOTP", email).Return(storedOTP, nil).Once()

		retUser, roles, accessToken, refreshToken, err := authUsecase.VerifyEmailOTP(email, rawOTP, "Test User", "Go-http-client/1.1", "127.0.0.1")

		assert.Error(t, err)
		assert.Nil(t, retUser)
		assert.Nil(t, roles)
		assert.Empty(t, accessToken)
		assert.Empty(t, refreshToken)
	})
}
