

package mocks

import (
	domain "github.com/aswinsreeraj/evntx/internal/domain"
	mock "github.com/stretchr/testify/mock"
)


type UserRoleRepository struct {
	mock.Mock
}


func (_m *UserRoleRepository) AddRole(userID string, role domain.UserRole) error {
	ret := _m.Called(userID, role)

	if len(ret) == 0 {
		panic("no return value specified for AddRole")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, domain.UserRole) error); ok {
		r0 = rf(userID, role)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *UserRoleRepository) GetRolesByUserID(userID string) ([]domain.UserRole, error) {
	ret := _m.Called(userID)

	if len(ret) == 0 {
		panic("no return value specified for GetRolesByUserID")
	}

	var r0 []domain.UserRole
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]domain.UserRole, error)); ok {
		return rf(userID)
	}
	if rf, ok := ret.Get(0).(func(string) []domain.UserRole); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.UserRole)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *UserRoleRepository) RemoveRole(userID string, role domain.UserRole) error {
	ret := _m.Called(userID, role)

	if len(ret) == 0 {
		panic("no return value specified for RemoveRole")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, domain.UserRole) error); ok {
		r0 = rf(userID, role)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}



func NewUserRoleRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserRoleRepository {
	mock := &UserRoleRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
