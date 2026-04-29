package mocks

import mock "github.com/stretchr/testify/mock"

type EmailSender struct {
	mock.Mock
}

func (_m *EmailSender) SendOTP(email string, otp string) error {
	ret := _m.Called(email, otp)

	if len(ret) == 0 {
		panic("no return value specified for SendOTP")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(email, otp)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (_m *EmailSender) SendOrganizerApproval(email string, name string) error {
	ret := _m.Called(email, name)

	if len(ret) == 0 {
		panic("no return value specified for SendOrganizerApproval")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(email, name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func NewEmailSender(t interface {
	mock.TestingT
	Cleanup(func())
}) *EmailSender {
	mock := &EmailSender{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
