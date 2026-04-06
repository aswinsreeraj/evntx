package repository

import "github.com/aswinsreeraj/evntx/internal/domain"

type SettingsRepository interface {
	GetPlatformSettings() (*domain.PlatformSettings, error)
	UpdatePlatformSettings(s *domain.PlatformSettings) error
	GetPaymentSettings() ([]domain.PaymentSettings, error)
	GetPaymentProviderConfig(provider string) (*domain.PaymentSettings, error)
	UpdatePaymentProvider(provider string, isEnabled bool, config map[string]interface{}) error
}
