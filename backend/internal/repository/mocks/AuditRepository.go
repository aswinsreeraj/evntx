

package mocks

import (
	domain "github.com/aswinsreeraj/evntx/internal/domain"
	mock "github.com/stretchr/testify/mock"
)


type AuditRepository struct {
	mock.Mock
}


func (_m *AuditRepository) Create(log *domain.AuditLog) error {
	ret := _m.Called(log)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.AuditLog) error); ok {
		r0 = rf(log)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *AuditRepository) GetLogs(page int, limit int) ([]domain.AuditLog, int64, error) {
	ret := _m.Called(page, limit)

	if len(ret) == 0 {
		panic("no return value specified for GetLogs")
	}

	var r0 []domain.AuditLog
	var r1 int64
	var r2 error
	if rf, ok := ret.Get(0).(func(int, int) ([]domain.AuditLog, int64, error)); ok {
		return rf(page, limit)
	}
	if rf, ok := ret.Get(0).(func(int, int) []domain.AuditLog); ok {
		r0 = rf(page, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.AuditLog)
		}
	}

	if rf, ok := ret.Get(1).(func(int, int) int64); ok {
		r1 = rf(page, limit)
	} else {
		r1 = ret.Get(1).(int64)
	}

	if rf, ok := ret.Get(2).(func(int, int) error); ok {
		r2 = rf(page, limit)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}



func NewAuditRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *AuditRepository {
	mock := &AuditRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
