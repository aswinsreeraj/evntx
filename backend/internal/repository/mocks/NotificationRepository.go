package mocks

import (
	domain "github.com/aswinsreeraj/evntx/internal/domain"
	mock "github.com/stretchr/testify/mock"
)

type NotificationRepository struct {
	mock.Mock
}

func (_m *NotificationRepository) ClearAll(userID string) error {
	ret := _m.Called(userID)

	if len(ret) == 0 {
		panic("no return value specified for ClearAll")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (_m *NotificationRepository) CreateNotification(notification *domain.Notification) error {
	ret := _m.Called(notification)

	if len(ret) == 0 {
		panic("no return value specified for CreateNotification")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.Notification) error); ok {
		r0 = rf(notification)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (_m *NotificationRepository) GetNotificationsByUser(userID string, page int, limit int) ([]domain.Notification, int64, int64, error) {
	ret := _m.Called(userID, page, limit)

	if len(ret) == 0 {
		panic("no return value specified for GetNotificationsByUser")
	}

	var r0 []domain.Notification
	var r1 int64
	var r2 int64
	var r3 error
	if rf, ok := ret.Get(0).(func(string, int, int) ([]domain.Notification, int64, int64, error)); ok {
		return rf(userID, page, limit)
	}
	if rf, ok := ret.Get(0).(func(string, int, int) []domain.Notification); ok {
		r0 = rf(userID, page, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Notification)
		}
	}

	if rf, ok := ret.Get(1).(func(string, int, int) int64); ok {
		r1 = rf(userID, page, limit)
	} else {
		r1 = ret.Get(1).(int64)
	}

	if rf, ok := ret.Get(2).(func(string, int, int) int64); ok {
		r2 = rf(userID, page, limit)
	} else {
		r2 = ret.Get(2).(int64)
	}

	if rf, ok := ret.Get(3).(func(string, int, int) error); ok {
		r3 = rf(userID, page, limit)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}

func (_m *NotificationRepository) MarkAllAsRead(userID string) error {
	ret := _m.Called(userID)

	if len(ret) == 0 {
		panic("no return value specified for MarkAllAsRead")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (_m *NotificationRepository) MarkAsRead(notificationID string, userID string) error {
	ret := _m.Called(notificationID, userID)

	if len(ret) == 0 {
		panic("no return value specified for MarkAsRead")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(notificationID, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func NewNotificationRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *NotificationRepository {
	mock := &NotificationRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
