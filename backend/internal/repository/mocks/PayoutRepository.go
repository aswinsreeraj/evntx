

package mocks

import (
	context "context"

	domain "github.com/aswinsreeraj/evntx/internal/domain"
	mock "github.com/stretchr/testify/mock"

	repository "github.com/aswinsreeraj/evntx/internal/repository"
)


type PayoutRepository struct {
	mock.Mock
}


func (_m *PayoutRepository) AdminGetPayoutRequests(ctx context.Context, status string, page int, limit int) ([]domain.AdminPayoutDetail, int64, error) {
	ret := _m.Called(ctx, status, page, limit)

	if len(ret) == 0 {
		panic("no return value specified for AdminGetPayoutRequests")
	}

	var r0 []domain.AdminPayoutDetail
	var r1 int64
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string, int, int) ([]domain.AdminPayoutDetail, int64, error)); ok {
		return rf(ctx, status, page, limit)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, int, int) []domain.AdminPayoutDetail); ok {
		r0 = rf(ctx, status, page, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.AdminPayoutDetail)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, int, int) int64); ok {
		r1 = rf(ctx, status, page, limit)
	} else {
		r1 = ret.Get(1).(int64)
	}

	if rf, ok := ret.Get(2).(func(context.Context, string, int, int) error); ok {
		r2 = rf(ctx, status, page, limit)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}


func (_m *PayoutRepository) CreatePayoutRequest(ctx context.Context, req *domain.PayoutRequest) error {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for CreatePayoutRequest")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.PayoutRequest) error); ok {
		r0 = rf(ctx, req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *PayoutRepository) GetCredentialByUserID(ctx context.Context, userID string) (*domain.PayoutCredential, error) {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetCredentialByUserID")
	}

	var r0 *domain.PayoutCredential
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*domain.PayoutCredential, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *domain.PayoutCredential); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.PayoutCredential)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *PayoutRepository) GetPayoutRequestByID(ctx context.Context, payoutID string) (*domain.PayoutRequest, error) {
	ret := _m.Called(ctx, payoutID)

	if len(ret) == 0 {
		panic("no return value specified for GetPayoutRequestByID")
	}

	var r0 *domain.PayoutRequest
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*domain.PayoutRequest, error)); ok {
		return rf(ctx, payoutID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *domain.PayoutRequest); ok {
		r0 = rf(ctx, payoutID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.PayoutRequest)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, payoutID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *PayoutRepository) GetPayoutRequestsByUserID(ctx context.Context, userID string, page int, limit int) ([]domain.PayoutRequest, int64, error) {
	ret := _m.Called(ctx, userID, page, limit)

	if len(ret) == 0 {
		panic("no return value specified for GetPayoutRequestsByUserID")
	}

	var r0 []domain.PayoutRequest
	var r1 int64
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string, int, int) ([]domain.PayoutRequest, int64, error)); ok {
		return rf(ctx, userID, page, limit)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, int, int) []domain.PayoutRequest); ok {
		r0 = rf(ctx, userID, page, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.PayoutRequest)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, int, int) int64); ok {
		r1 = rf(ctx, userID, page, limit)
	} else {
		r1 = ret.Get(1).(int64)
	}

	if rf, ok := ret.Get(2).(func(context.Context, string, int, int) error); ok {
		r2 = rf(ctx, userID, page, limit)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}


func (_m *PayoutRepository) UpdatePayoutRequestStatus(ctx context.Context, payoutID string, newStatus domain.PayoutStatus, adminID *string, failureReason *string) error {
	ret := _m.Called(ctx, payoutID, newStatus, adminID, failureReason)

	if len(ret) == 0 {
		panic("no return value specified for UpdatePayoutRequestStatus")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, domain.PayoutStatus, *string, *string) error); ok {
		r0 = rf(ctx, payoutID, newStatus, adminID, failureReason)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *PayoutRepository) UpsertCredential(ctx context.Context, cred *domain.PayoutCredential) error {
	ret := _m.Called(ctx, cred)

	if len(ret) == 0 {
		panic("no return value specified for UpsertCredential")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.PayoutCredential) error); ok {
		r0 = rf(ctx, cred)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *PayoutRepository) WithTransaction(fn func(repository.PayoutRepository) error) error {
	ret := _m.Called(fn)

	if len(ret) == 0 {
		panic("no return value specified for WithTransaction")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(func(repository.PayoutRepository) error) error); ok {
		r0 = rf(fn)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}



func NewPayoutRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *PayoutRepository {
	mock := &PayoutRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
