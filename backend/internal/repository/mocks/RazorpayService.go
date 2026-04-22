

package mocks

import (
	domain "github.com/aswinsreeraj/evntx/internal/domain"
	mock "github.com/stretchr/testify/mock"
)


type RazorpayService struct {
	mock.Mock
}


func (_m *RazorpayService) CreateOrder(amount int64, receipt string) (*domain.RazorpayOrder, error) {
	ret := _m.Called(amount, receipt)

	if len(ret) == 0 {
		panic("no return value specified for CreateOrder")
	}

	var r0 *domain.RazorpayOrder
	var r1 error
	if rf, ok := ret.Get(0).(func(int64, string) (*domain.RazorpayOrder, error)); ok {
		return rf(amount, receipt)
	}
	if rf, ok := ret.Get(0).(func(int64, string) *domain.RazorpayOrder); ok {
		r0 = rf(amount, receipt)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.RazorpayOrder)
		}
	}

	if rf, ok := ret.Get(1).(func(int64, string) error); ok {
		r1 = rf(amount, receipt)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *RazorpayService) FetchOrder(orderID string) (*domain.RazorpayOrder, error) {
	ret := _m.Called(orderID)

	if len(ret) == 0 {
		panic("no return value specified for FetchOrder")
	}

	var r0 *domain.RazorpayOrder
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*domain.RazorpayOrder, error)); ok {
		return rf(orderID)
	}
	if rf, ok := ret.Get(0).(func(string) *domain.RazorpayOrder); ok {
		r0 = rf(orderID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.RazorpayOrder)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(orderID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *RazorpayService) GetKeyID() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetKeyID")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}


func (_m *RazorpayService) RefundPayment(paymentID string, amount int64) error {
	ret := _m.Called(paymentID, amount)

	if len(ret) == 0 {
		panic("no return value specified for RefundPayment")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, int64) error); ok {
		r0 = rf(paymentID, amount)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *RazorpayService) VerifySignature(orderID string, paymentID string, signature string) (bool, error) {
	ret := _m.Called(orderID, paymentID, signature)

	if len(ret) == 0 {
		panic("no return value specified for VerifySignature")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string, string) (bool, error)); ok {
		return rf(orderID, paymentID, signature)
	}
	if rf, ok := ret.Get(0).(func(string, string, string) bool); ok {
		r0 = rf(orderID, paymentID, signature)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(string, string, string) error); ok {
		r1 = rf(orderID, paymentID, signature)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}



func NewRazorpayService(t interface {
	mock.TestingT
	Cleanup(func())
}) *RazorpayService {
	mock := &RazorpayService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
