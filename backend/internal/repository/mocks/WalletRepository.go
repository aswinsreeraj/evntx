

package mocks

import (
	domain "github.com/aswinsreeraj/evntx/internal/domain"
	mock "github.com/stretchr/testify/mock"

	repository "github.com/aswinsreeraj/evntx/internal/repository"
)


type WalletRepository struct {
	mock.Mock
}


func (_m *WalletRepository) CreateTransaction(txn *domain.WalletTransaction) error {
	ret := _m.Called(txn)

	if len(ret) == 0 {
		panic("no return value specified for CreateTransaction")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.WalletTransaction) error); ok {
		r0 = rf(txn)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *WalletRepository) CreateWallet(wallet *domain.Wallet) error {
	ret := _m.Called(wallet)

	if len(ret) == 0 {
		panic("no return value specified for CreateWallet")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.Wallet) error); ok {
		r0 = rf(wallet)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *WalletRepository) GetTransactionsByWalletID(walletID string, filters domain.WalletTransactionFilter, page int, limit int) ([]domain.WalletTransaction, int64, error) {
	ret := _m.Called(walletID, filters, page, limit)

	if len(ret) == 0 {
		panic("no return value specified for GetTransactionsByWalletID")
	}

	var r0 []domain.WalletTransaction
	var r1 int64
	var r2 error
	if rf, ok := ret.Get(0).(func(string, domain.WalletTransactionFilter, int, int) ([]domain.WalletTransaction, int64, error)); ok {
		return rf(walletID, filters, page, limit)
	}
	if rf, ok := ret.Get(0).(func(string, domain.WalletTransactionFilter, int, int) []domain.WalletTransaction); ok {
		r0 = rf(walletID, filters, page, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.WalletTransaction)
		}
	}

	if rf, ok := ret.Get(1).(func(string, domain.WalletTransactionFilter, int, int) int64); ok {
		r1 = rf(walletID, filters, page, limit)
	} else {
		r1 = ret.Get(1).(int64)
	}

	if rf, ok := ret.Get(2).(func(string, domain.WalletTransactionFilter, int, int) error); ok {
		r2 = rf(walletID, filters, page, limit)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}


func (_m *WalletRepository) GetWalletByID(walletID string) (*domain.Wallet, error) {
	ret := _m.Called(walletID)

	if len(ret) == 0 {
		panic("no return value specified for GetWalletByID")
	}

	var r0 *domain.Wallet
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*domain.Wallet, error)); ok {
		return rf(walletID)
	}
	if rf, ok := ret.Get(0).(func(string) *domain.Wallet); ok {
		r0 = rf(walletID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Wallet)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(walletID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *WalletRepository) GetWalletByUserID(userID string) (*domain.Wallet, error) {
	ret := _m.Called(userID)

	if len(ret) == 0 {
		panic("no return value specified for GetWalletByUserID")
	}

	var r0 *domain.Wallet
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*domain.Wallet, error)); ok {
		return rf(userID)
	}
	if rf, ok := ret.Get(0).(func(string) *domain.Wallet); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Wallet)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *WalletRepository) UpdateTransactionStatusByReference(refType string, refID string, status string) error {
	ret := _m.Called(refType, refID, status)

	if len(ret) == 0 {
		panic("no return value specified for UpdateTransactionStatusByReference")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string) error); ok {
		r0 = rf(refType, refID, status)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *WalletRepository) UpdateWallet(wallet *domain.Wallet) error {
	ret := _m.Called(wallet)

	if len(ret) == 0 {
		panic("no return value specified for UpdateWallet")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.Wallet) error); ok {
		r0 = rf(wallet)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *WalletRepository) WithTransaction(fn func(repository.WalletRepository) error) error {
	ret := _m.Called(fn)

	if len(ret) == 0 {
		panic("no return value specified for WithTransaction")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(func(repository.WalletRepository) error) error); ok {
		r0 = rf(fn)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}



func NewWalletRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *WalletRepository {
	mock := &WalletRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
