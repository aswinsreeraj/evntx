package mocks

import (
	context "context"

	domain "github.com/aswinsreeraj/evntx/internal/domain"
	mock "github.com/stretchr/testify/mock"

	time "time"
)

type EngagementRepository struct {
	mock.Mock
}

func (_m *EngagementRepository) CreateSession(ctx context.Context, session *domain.VisitorSession) error {
	ret := _m.Called(ctx, session)

	if len(ret) == 0 {
		panic("no return value specified for CreateSession")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.VisitorSession) error); ok {
		r0 = rf(ctx, session)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (_m *EngagementRepository) GetDailyAggregates(ctx context.Context, eventID string, startDate time.Time, endDate time.Time) ([]domain.EventEngagementDaily, error) {
	ret := _m.Called(ctx, eventID, startDate, endDate)

	if len(ret) == 0 {
		panic("no return value specified for GetDailyAggregates")
	}

	var r0 []domain.EventEngagementDaily
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, time.Time, time.Time) ([]domain.EventEngagementDaily, error)); ok {
		return rf(ctx, eventID, startDate, endDate)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, time.Time, time.Time) []domain.EventEngagementDaily); ok {
		r0 = rf(ctx, eventID, startDate, endDate)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.EventEngagementDaily)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, time.Time, time.Time) error); ok {
		r1 = rf(ctx, eventID, startDate, endDate)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *EngagementRepository) GetEngagementReport(ctx context.Context, eventIDs []string, startDate time.Time, endDate time.Time) (*domain.EngagementReportStats, error) {
	ret := _m.Called(ctx, eventIDs, startDate, endDate)

	if len(ret) == 0 {
		panic("no return value specified for GetEngagementReport")
	}

	var r0 *domain.EngagementReportStats
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []string, time.Time, time.Time) (*domain.EngagementReportStats, error)); ok {
		return rf(ctx, eventIDs, startDate, endDate)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []string, time.Time, time.Time) *domain.EngagementReportStats); ok {
		r0 = rf(ctx, eventIDs, startDate, endDate)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.EngagementReportStats)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []string, time.Time, time.Time) error); ok {
		r1 = rf(ctx, eventIDs, startDate, endDate)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *EngagementRepository) GetSessionByID(ctx context.Context, sessionID string) (*domain.VisitorSession, error) {
	ret := _m.Called(ctx, sessionID)

	if len(ret) == 0 {
		panic("no return value specified for GetSessionByID")
	}

	var r0 *domain.VisitorSession
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*domain.VisitorSession, error)); ok {
		return rf(ctx, sessionID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *domain.VisitorSession); ok {
		r0 = rf(ctx, sessionID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.VisitorSession)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, sessionID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *EngagementRepository) IncrementSuccessfulBookings(ctx context.Context, eventID string, userID string) error {
	ret := _m.Called(ctx, eventID, userID)

	if len(ret) == 0 {
		panic("no return value specified for IncrementSuccessfulBookings")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, eventID, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (_m *EngagementRepository) LogEvent(ctx context.Context, event *domain.EngagementEvent) error {
	ret := _m.Called(ctx, event)

	if len(ret) == 0 {
		panic("no return value specified for LogEvent")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.EngagementEvent) error); ok {
		r0 = rf(ctx, event)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (_m *EngagementRepository) UpdateSessionLastSeen(ctx context.Context, sessionID string, userID *string) error {
	ret := _m.Called(ctx, sessionID, userID)

	if len(ret) == 0 {
		panic("no return value specified for UpdateSessionLastSeen")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *string) error); ok {
		r0 = rf(ctx, sessionID, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func NewEngagementRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *EngagementRepository {
	mock := &EngagementRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
