

package mocks

import (
	domain "github.com/aswinsreeraj/evntx/internal/domain"
	mock "github.com/stretchr/testify/mock"
)


type UserRepository struct {
	mock.Mock
}


func (_m *UserRepository) Create(user *domain.User) error {
	ret := _m.Called(user)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.User) error); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *UserRepository) Delete(id string) error {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *UserRepository) FindByEmail(email string) (*domain.User, error) {
	ret := _m.Called(email)

	if len(ret) == 0 {
		panic("no return value specified for FindByEmail")
	}

	var r0 *domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*domain.User, error)); ok {
		return rf(email)
	}
	if rf, ok := ret.Get(0).(func(string) *domain.User); ok {
		r0 = rf(email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.User)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *UserRepository) FindByID(id string) (*domain.User, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for FindByID")
	}

	var r0 *domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*domain.User, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(string) *domain.User); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.User)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *UserRepository) FindUsersByRole(role domain.UserRole) ([]domain.User, error) {
	ret := _m.Called(role)

	if len(ret) == 0 {
		panic("no return value specified for FindUsersByRole")
	}

	var r0 []domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func(domain.UserRole) ([]domain.User, error)); ok {
		return rf(role)
	}
	if rf, ok := ret.Get(0).(func(domain.UserRole) []domain.User); ok {
		r0 = rf(role)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.User)
		}
	}

	if rf, ok := ret.Get(1).(func(domain.UserRole) error); ok {
		r1 = rf(role)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *UserRepository) GetOrganizerDetails(userID string) (*domain.OrganizerDetail, error) {
	ret := _m.Called(userID)

	if len(ret) == 0 {
		panic("no return value specified for GetOrganizerDetails")
	}

	var r0 *domain.OrganizerDetail
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*domain.OrganizerDetail, error)); ok {
		return rf(userID)
	}
	if rf, ok := ret.Get(0).(func(string) *domain.OrganizerDetail); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.OrganizerDetail)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *UserRepository) Search(search string, status string, page int, limit int) ([]domain.AdminUserDetails, int64, error) {
	ret := _m.Called(search, status, page, limit)

	if len(ret) == 0 {
		panic("no return value specified for Search")
	}

	var r0 []domain.AdminUserDetails
	var r1 int64
	var r2 error
	if rf, ok := ret.Get(0).(func(string, string, int, int) ([]domain.AdminUserDetails, int64, error)); ok {
		return rf(search, status, page, limit)
	}
	if rf, ok := ret.Get(0).(func(string, string, int, int) []domain.AdminUserDetails); ok {
		r0 = rf(search, status, page, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.AdminUserDetails)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string, int, int) int64); ok {
		r1 = rf(search, status, page, limit)
	} else {
		r1 = ret.Get(1).(int64)
	}

	if rf, ok := ret.Get(2).(func(string, string, int, int) error); ok {
		r2 = rf(search, status, page, limit)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}


func (_m *UserRepository) SearchOrganizers(search string, status string, page int, limit int) ([]domain.OrganizerDetails, int64, error) {
	ret := _m.Called(search, status, page, limit)

	if len(ret) == 0 {
		panic("no return value specified for SearchOrganizers")
	}

	var r0 []domain.OrganizerDetails
	var r1 int64
	var r2 error
	if rf, ok := ret.Get(0).(func(string, string, int, int) ([]domain.OrganizerDetails, int64, error)); ok {
		return rf(search, status, page, limit)
	}
	if rf, ok := ret.Get(0).(func(string, string, int, int) []domain.OrganizerDetails); ok {
		r0 = rf(search, status, page, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.OrganizerDetails)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string, int, int) int64); ok {
		r1 = rf(search, status, page, limit)
	} else {
		r1 = ret.Get(1).(int64)
	}

	if rf, ok := ret.Get(2).(func(string, string, int, int) error); ok {
		r2 = rf(search, status, page, limit)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}


func (_m *UserRepository) Update(user *domain.User) error {
	ret := _m.Called(user)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.User) error); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *UserRepository) UpdateOrganizerApprovalStatus(userID string, approvalStatus string) error {
	ret := _m.Called(userID, approvalStatus)

	if len(ret) == 0 {
		panic("no return value specified for UpdateOrganizerApprovalStatus")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(userID, approvalStatus)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *UserRepository) UpdateStatus(userID string, isActive bool) error {
	ret := _m.Called(userID, isActive)

	if len(ret) == 0 {
		panic("no return value specified for UpdateStatus")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, bool) error); ok {
		r0 = rf(userID, isActive)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *UserRepository) UpsertOrganizerDetails(detail *domain.OrganizerDetail) error {
	ret := _m.Called(detail)

	if len(ret) == 0 {
		panic("no return value specified for UpsertOrganizerDetails")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.OrganizerDetail) error); ok {
		r0 = rf(detail)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}



func NewUserRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserRepository {
	mock := &UserRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
