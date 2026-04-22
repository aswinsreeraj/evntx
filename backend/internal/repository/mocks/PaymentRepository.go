

package mocks

import (
	domain "github.com/aswinsreeraj/evntx/internal/domain"
	mock "github.com/stretchr/testify/mock"
)


type PaymentRepository struct {
	mock.Mock
}


func (_m *PaymentRepository) CreatePayment(payment *domain.Payment) error {
	ret := _m.Called(payment)

	if len(ret) == 0 {
		panic("no return value specified for CreatePayment")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.Payment) error); ok {
		r0 = rf(payment)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *PaymentRepository) FindByBookingID(bookingID string) (*domain.Payment, error) {
	ret := _m.Called(bookingID)

	if len(ret) == 0 {
		panic("no return value specified for FindByBookingID")
	}

	var r0 *domain.Payment
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*domain.Payment, error)); ok {
		return rf(bookingID)
	}
	if rf, ok := ret.Get(0).(func(string) *domain.Payment); ok {
		r0 = rf(bookingID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Payment)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(bookingID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *PaymentRepository) FindByProviderReference(orderID string) (*domain.Payment, error) {
	ret := _m.Called(orderID)

	if len(ret) == 0 {
		panic("no return value specified for FindByProviderReference")
	}

	var r0 *domain.Payment
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*domain.Payment, error)); ok {
		return rf(orderID)
	}
	if rf, ok := ret.Get(0).(func(string) *domain.Payment); ok {
		r0 = rf(orderID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Payment)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(orderID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *PaymentRepository) MarkPaymentSuccess(paymentID string, bookingID string, organizerID string, amount float64) error {
	ret := _m.Called(paymentID, bookingID, organizerID, amount)

	if len(ret) == 0 {
		panic("no return value specified for MarkPaymentSuccess")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string, float64) error); ok {
		r0 = rf(paymentID, bookingID, organizerID, amount)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *PaymentRepository) RefundPaymentToWallet(userID string, paymentID string, bookingID string, refundAmount float64, platformFeeAmount float64) error {
	ret := _m.Called(userID, paymentID, bookingID, refundAmount, platformFeeAmount)

	if len(ret) == 0 {
		panic("no return value specified for RefundPaymentToWallet")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string, float64, float64) error); ok {
		r0 = rf(userID, paymentID, bookingID, refundAmount, platformFeeAmount)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *PaymentRepository) UpdateStatus(paymentID string, status string) error {
	ret := _m.Called(paymentID, status)

	if len(ret) == 0 {
		panic("no return value specified for UpdateStatus")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(paymentID, status)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}



func NewPaymentRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *PaymentRepository {
	mock := &PaymentRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
