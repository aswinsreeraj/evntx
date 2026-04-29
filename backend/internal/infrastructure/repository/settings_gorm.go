package repository

import (
	"encoding/json"
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PlatformSettingsModel struct {
	ID					string		`gorm:"type:uuid;primaryKey"`
	EnableUserRegistration			bool		`gorm:"default:true;not null"`
	AllowGoogleLogin			bool		`gorm:"default:true;not null"`
	RequireAdminApprovalForOrganizers	bool		`gorm:"default:false;not null"`
	RequireAdminApprovalForEvents		bool		`gorm:"default:true;not null"`
	RefundWindowDays			int		`gorm:"default:3;not null"`
	AllowEventCancellation			bool		`gorm:"default:true;not null"`
	PlatformFeeValue			float64		`gorm:"type:numeric(10,2);default:30;not null"`
	PlatformFeeType				string		`gorm:"size:20;default:'fixed';not null"`
	UpdatedAt				time.Time	`gorm:"not null"`
}

type PaymentSettingsModel struct {
	ID		string		`gorm:"type:uuid;primaryKey"`
	Provider	string		`gorm:"uniqueIndex;not null"`
	IsEnabled	bool		`gorm:"default:false;not null"`
	Config		json.RawMessage	`gorm:"type:jsonb"`
	CreatedAt	time.Time	`gorm:"not null"`
	UpdatedAt	time.Time	`gorm:"not null"`
}

type settingsGormRepository struct {
	db *gorm.DB
}

func NewSettingsGormRepository(db *gorm.DB) *settingsGormRepository {
	return &settingsGormRepository{db: db}
}

func (r *settingsGormRepository) EnsureExists() error {
	platformModel := PlatformSettingsModel{
		ID:					domain.PlatformSettingsID,
		EnableUserRegistration:			true,
		AllowGoogleLogin:			true,
		RequireAdminApprovalForOrganizers:	false,
		RequireAdminApprovalForEvents:		true,
		RefundWindowDays:			3,
		AllowEventCancellation:			true,
		PlatformFeeValue:			30,
		PlatformFeeType:			string(domain.PlatformFeeTypeFixed),
		UpdatedAt:				time.Now(),
	}

	if err := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&platformModel).Error; err != nil {
		return err
	}

	razorpayModel := PaymentSettingsModel{
		ID:		uuid.NewString(),
		Provider:	"razorpay",
		IsEnabled:	true,
		Config:		json.RawMessage(`{}`),
		CreatedAt:	time.Now(),
		UpdatedAt:	time.Now(),
	}

	if err := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&razorpayModel).Error; err != nil {
		return err
	}

	walletModel := PaymentSettingsModel{
		ID:		uuid.NewString(),
		Provider:	"wallet",
		IsEnabled:	true,
		Config:		json.RawMessage(`{}`),
		CreatedAt:	time.Now(),
		UpdatedAt:	time.Now(),
	}

	return r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&walletModel).Error
}

func (r *settingsGormRepository) GetPlatformSettings() (*domain.PlatformSettings, error) {
	var model PlatformSettingsModel
	if err := r.db.Where("id = ?", domain.PlatformSettingsID).First(&model).Error; err != nil {
		return nil, err
	}
	return toDomainPlatformSettings(&model), nil
}

func (r *settingsGormRepository) UpdatePlatformSettings(s *domain.PlatformSettings) error {
	s.UpdatedAt = time.Now()
	return r.db.Model(&PlatformSettingsModel{}).
		Where("id = ?", domain.PlatformSettingsID).
		Updates(map[string]interface{}{
			"enable_user_registration":			s.EnableUserRegistration,
			"allow_google_login":				s.AllowGoogleLogin,
			"require_admin_approval_for_organizers":	s.RequireAdminApprovalForOrganizers,
			"require_admin_approval_for_events":		s.RequireAdminApprovalForEvents,
			"refund_window_days":				s.RefundWindowDays,
			"allow_event_cancellation":			s.AllowEventCancellation,
			"platform_fee_value":				s.PlatformFeeValue,
			"platform_fee_type":				string(s.PlatformFeeType),
			"updated_at":					s.UpdatedAt,
		}).Error
}

func (r *settingsGormRepository) GetPaymentSettings() ([]domain.PaymentSettings, error) {
	var models []PaymentSettingsModel
	if err := r.db.Find(&models).Error; err != nil {
		return nil, err
	}

	var result []domain.PaymentSettings
	for _, m := range models {
		result = append(result, *toDomainPaymentSettings(&m))
	}
	return result, nil
}

func (r *settingsGormRepository) GetPaymentProviderConfig(provider string) (*domain.PaymentSettings, error) {
	var model PaymentSettingsModel
	if err := r.db.Where("provider = ?", provider).First(&model).Error; err != nil {
		return nil, err
	}
	return toDomainPaymentSettings(&model), nil
}

func (r *settingsGormRepository) UpdatePaymentProvider(provider string, isEnabled bool, config map[string]interface{}) error {
	configBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}
	if configBytes == nil || string(configBytes) == "null" {
		configBytes = []byte(`{}`)
	}

	return r.db.Model(&PaymentSettingsModel{}).
		Where("provider = ?", provider).
		Updates(map[string]interface{}{
			"is_enabled":	isEnabled,
			"config":	json.RawMessage(configBytes),
			"updated_at":	time.Now(),
		}).Error
}

func toDomainPlatformSettings(m *PlatformSettingsModel) *domain.PlatformSettings {
	return &domain.PlatformSettings{
		ID:					m.ID,
		EnableUserRegistration:			m.EnableUserRegistration,
		AllowGoogleLogin:			m.AllowGoogleLogin,
		RequireAdminApprovalForOrganizers:	m.RequireAdminApprovalForOrganizers,
		RequireAdminApprovalForEvents:		m.RequireAdminApprovalForEvents,
		RefundWindowDays:			m.RefundWindowDays,
		AllowEventCancellation:			m.AllowEventCancellation,
		PlatformFeeValue:			m.PlatformFeeValue,
		PlatformFeeType:			domain.PlatformFeeType(m.PlatformFeeType),
		UpdatedAt:				m.UpdatedAt,
	}
}

func toDomainPaymentSettings(m *PaymentSettingsModel) *domain.PaymentSettings {
	var config map[string]interface{}
	if len(m.Config) > 0 {
		_ = json.Unmarshal(m.Config, &config)
	}
	if config == nil {
		config = make(map[string]interface{})
	}

	return &domain.PaymentSettings{
		ID:		m.ID,
		Provider:	m.Provider,
		IsEnabled:	m.IsEnabled,
		Config:		config,
		CreatedAt:	m.CreatedAt,
		UpdatedAt:	m.UpdatedAt,
	}
}
