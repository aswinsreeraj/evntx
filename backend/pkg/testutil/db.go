package testutil

import (
	"fmt"
	"os"
	"testing"

	repoImpl "github.com/aswinsreeraj/evntx/internal/infrastructure/repository"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupTestDB(t *testing.T) *gorm.DB {

	_ = godotenv.Load("../../.env")

	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=evntx_test port=5432 sslmode=disable TimeZone=UTC"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	err = db.AutoMigrate(
		&repoImpl.UserModel{},
		&repoImpl.WalletModel{},
		&repoImpl.WalletTransactionModel{},
		&repoImpl.OrganizerDetailModel{},
		&repoImpl.EmailOTPModel{},
		&repoImpl.UserSessionModel{},
		&repoImpl.UserRoleModel{},
		&repoImpl.EventModel{},
		&repoImpl.EventDetailsModel{},
		&repoImpl.EventPersonnelModel{},
		&repoImpl.TicketTypeModel{},
		&repoImpl.EventModerationLogModel{},
		&repoImpl.BookingModel{},
		&repoImpl.BookingTicketModel{},
		&repoImpl.TicketModel{},
		&repoImpl.PaymentModel{},
		&repoImpl.NotificationModel{},
		&repoImpl.PlatformWalletModel{},
		&repoImpl.PlatformWalletTransactionModel{},
		&repoImpl.PayoutRequestModel{},
		&repoImpl.PayoutCredentialModel{},
		&repoImpl.VisitorSessionModel{},
		&repoImpl.EngagementEventModel{},
		&repoImpl.EventEngagementDailyModel{},
		&repoImpl.PlatformSettingsModel{},
		&repoImpl.PaymentSettingsModel{},
		&repoImpl.AuditLogModel{},
		&repoImpl.JobLogModel{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func ClearDatabase(db *gorm.DB) {

	db.Exec("SET session_replication_role = 'replica';")

	tables := []string{
		"job_log_models", "audit_log_models", "payment_settings_models", "platform_settings_models",
		"event_engagement_daily_models", "engagement_event_models", "visitor_session_models",
		"payout_credential_models", "payout_request_models", "platform_wallet_transaction_models", "platform_wallet_models",
		"notification_models", "payment_models", "ticket_models", "booking_ticket_models", "booking_models",
		"event_moderation_log_models", "ticket_type_models", "event_personnel_models", "event_details_models", "event_models",
		"user_role_models", "user_session_models", "email_otp_models", "organizer_detail_models", "wallet_transaction_models", "wallet_models", "user_models",
	}

	for _, table := range tables {
		db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE;", table))
	}

	db.Exec("SET session_replication_role = 'origin';")
}
