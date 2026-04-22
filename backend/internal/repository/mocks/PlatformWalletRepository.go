

package mocks

import (
	domain "github.com/aswinsreeraj/evntx/internal/domain"
	mock "github.com/stretchr/testify/mock"
)


type PlatformWalletRepository struct {
	mock.Mock
}


func (_m *PlatformWalletRepository) ApplyPlatformTransaction(txnType string, amount float64, referenceType string, referenceID string) error {
	ret := _m.Called(txnType, amount, referenceType, referenceID)

	if len(ret) == 0 {
		panic("no return value specified for ApplyPlatformTransaction")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, float64, string, string) error); ok {
		r0 = rf(txnType, amount, referenceType, referenceID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *PlatformWalletRepository) EnsureExists() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for EnsureExists")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *PlatformWalletRepository) GetPlatformWallet() (*domain.PlatformWallet, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetPlatformWallet")
	}

	var r0 *domain.PlatformWallet
	var r1 error
	if rf, ok := ret.Get(0).(func() (*domain.PlatformWallet, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() *domain.PlatformWallet); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.PlatformWallet)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}



func NewPlatformWalletRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *PlatformWalletRepository {
	mock := &PlatformWalletRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
