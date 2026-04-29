package mocks

import (
	domain "github.com/aswinsreeraj/evntx/internal/domain"
	mock "github.com/stretchr/testify/mock"
)

type UserSessionRepository struct {
	mock.Mock
}

func (_m *UserSessionRepository) Create(session *domain.UserSession) error {
	ret := _m.Called(session)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.UserSession) error); ok {
		r0 = rf(session)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (_m *UserSessionRepository) FindByUserID(userID string) (*domain.UserSession, error) {
	ret := _m.Called(userID)

	if len(ret) == 0 {
		panic("no return value specified for FindByUserID")
	}

	var r0 *domain.UserSession
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*domain.UserSession, error)); ok {
		return rf(userID)
	}
	if rf, ok := ret.Get(0).(func(string) *domain.UserSession); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.UserSession)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *UserSessionRepository) Revoke(sessionID string) error {
	ret := _m.Called(sessionID)

	if len(ret) == 0 {
		panic("no return value specified for Revoke")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(sessionID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func NewUserSessionRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserSessionRepository {
	mock := &UserSessionRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
