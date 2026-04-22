

package mocks

import (
	domain "github.com/aswinsreeraj/evntx/internal/domain"
	mock "github.com/stretchr/testify/mock"
)


type JobRepository struct {
	mock.Mock
}


func (_m *JobRepository) LogJob(log *domain.JobLog) error {
	ret := _m.Called(log)

	if len(ret) == 0 {
		panic("no return value specified for LogJob")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.JobLog) error); ok {
		r0 = rf(log)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}



func NewJobRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *JobRepository {
	mock := &JobRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
