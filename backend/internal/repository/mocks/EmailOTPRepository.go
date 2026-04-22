

package mocks

import (
	domain "github.com/aswinsreeraj/evntx/internal/domain"
	mock "github.com/stretchr/testify/mock"
)


type EmailOTPRepository struct {
	mock.Mock
}


func (_m *EmailOTPRepository) Create(otp *domain.EmailOTP) error {
	ret := _m.Called(otp)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.EmailOTP) error); ok {
		r0 = rf(otp)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *EmailOTPRepository) FindValidOTP(email string) (*domain.EmailOTP, error) {
	ret := _m.Called(email)

	if len(ret) == 0 {
		panic("no return value specified for FindValidOTP")
	}

	var r0 *domain.EmailOTP
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*domain.EmailOTP, error)); ok {
		return rf(email)
	}
	if rf, ok := ret.Get(0).(func(string) *domain.EmailOTP); ok {
		r0 = rf(email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.EmailOTP)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *EmailOTPRepository) InvalidatePrevious(email string) error {
	ret := _m.Called(email)

	if len(ret) == 0 {
		panic("no return value specified for InvalidatePrevious")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(email)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *EmailOTPRepository) MarkConsumed(id string) error {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for MarkConsumed")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}



func NewEmailOTPRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *EmailOTPRepository {
	mock := &EmailOTPRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
