package integration

import (
	"testing"

	"github.com/aswinsreeraj/evntx/internal/domain"
	repoImpl "github.com/aswinsreeraj/evntx/internal/infrastructure/repository"
	"github.com/aswinsreeraj/evntx/internal/usecase"
	"github.com/aswinsreeraj/evntx/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestWalletIntegration_Transactions(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.ClearDatabase(db)

	walletRepo := repoImpl.NewWalletGormRepository(db)
	roleRepo := repoImpl.NewUserRoleGormRepository(db)
	platformWalletRepo := repoImpl.NewPlatformWalletGormRepository(db)
	bookingRepo := repoImpl.NewBookingGormRepository(db)
	payoutRepo := repoImpl.NewPayoutGormRepository(db)
	notificationRepo := repoImpl.NewNotificationGormRepository(db)
	notificationUsecase := usecase.NewNotificationUsecase(notificationRepo)

	paymentService := &testutil.MockPaymentService{}
	walletUsecase := usecase.NewWalletUsecase(walletRepo, roleRepo, platformWalletRepo, paymentService, bookingRepo, payoutRepo, notificationUsecase)

	
	user := testutil.SeedUser(db, "wallet_user@evntx.com", "user")

	wallet, err := walletUsecase.GetWalletByUserID(user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, wallet)
	assert.Equal(t, 0.0, wallet.AvailableBalance)

	
	err = walletUsecase.ApplyTransaction(wallet.ID, domain.WalletTransactionTypeCredit, 500.0, domain.WalletReferenceTypeFundAddition, "ref_123")
	assert.NoError(t, err)

	wallet, _ = walletUsecase.GetWalletByUserID(user.ID)
	assert.Equal(t, 500.0, wallet.AvailableBalance)
	assert.Equal(t, 500.0, wallet.TotalCredited)

	
	err = walletUsecase.ApplyTransaction(wallet.ID, domain.WalletTransactionTypeDebit, 200.0, domain.WalletReferenceTypePurchase, "booking_123")
	assert.NoError(t, err)

	wallet, _ = walletUsecase.GetWalletByUserID(user.ID)
	assert.Equal(t, 300.0, wallet.AvailableBalance)
	assert.Equal(t, 200.0, wallet.TotalDebited)

	
	err = walletUsecase.ApplyTransaction(wallet.ID, domain.WalletTransactionTypeDebit, 400.0, domain.WalletReferenceTypePurchase, "booking_456")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nsufficient balance")

	wallet, _ = walletUsecase.GetWalletByUserID(user.ID)
	assert.Equal(t, 300.0, wallet.AvailableBalance) 
}
