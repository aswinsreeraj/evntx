

package mocks

import (
	domain "github.com/aswinsreeraj/evntx/internal/domain"
	mock "github.com/stretchr/testify/mock"
)


type SettingsRepository struct {
	mock.Mock
}


func (_m *SettingsRepository) GetPaymentProviderConfig(provider string) (*domain.PaymentSettings, error) {
	ret := _m.Called(provider)

	if len(ret) == 0 {
		panic("no return value specified for GetPaymentProviderConfig")
	}

	var r0 *domain.PaymentSettings
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*domain.PaymentSettings, error)); ok {
		return rf(provider)
	}
	if rf, ok := ret.Get(0).(func(string) *domain.PaymentSettings); ok {
		r0 = rf(provider)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.PaymentSettings)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(provider)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *SettingsRepository) GetPaymentSettings() ([]domain.PaymentSettings, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetPaymentSettings")
	}

	var r0 []domain.PaymentSettings
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]domain.PaymentSettings, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []domain.PaymentSettings); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.PaymentSettings)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *SettingsRepository) GetPlatformSettings() (*domain.PlatformSettings, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetPlatformSettings")
	}

	var r0 *domain.PlatformSettings
	var r1 error
	if rf, ok := ret.Get(0).(func() (*domain.PlatformSettings, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() *domain.PlatformSettings); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.PlatformSettings)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *SettingsRepository) UpdatePaymentProvider(provider string, isEnabled bool, config map[string]interface{}) error {
	ret := _m.Called(provider, isEnabled, config)

	if len(ret) == 0 {
		panic("no return value specified for UpdatePaymentProvider")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, bool, map[string]interface{}) error); ok {
		r0 = rf(provider, isEnabled, config)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *SettingsRepository) UpdatePlatformSettings(s *domain.PlatformSettings) error {
	ret := _m.Called(s)

	if len(ret) == 0 {
		panic("no return value specified for UpdatePlatformSettings")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.PlatformSettings) error); ok {
		r0 = rf(s)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}



func NewSettingsRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *SettingsRepository {
	mock := &SettingsRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
