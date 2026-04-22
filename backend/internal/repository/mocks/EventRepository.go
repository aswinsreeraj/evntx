

package mocks

import (
	context "context"

	domain "github.com/aswinsreeraj/evntx/internal/domain"
	mock "github.com/stretchr/testify/mock"

	time "time"
)


type EventRepository struct {
	mock.Mock
}


func (_m *EventRepository) AdminSearchEvents(search string, status string, page int, limit int) ([]domain.AdminEventDetails, int64, error) {
	ret := _m.Called(search, status, page, limit)

	if len(ret) == 0 {
		panic("no return value specified for AdminSearchEvents")
	}

	var r0 []domain.AdminEventDetails
	var r1 int64
	var r2 error
	if rf, ok := ret.Get(0).(func(string, string, int, int) ([]domain.AdminEventDetails, int64, error)); ok {
		return rf(search, status, page, limit)
	}
	if rf, ok := ret.Get(0).(func(string, string, int, int) []domain.AdminEventDetails); ok {
		r0 = rf(search, status, page, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.AdminEventDetails)
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


func (_m *EventRepository) ApproveEvent(ctx context.Context, eventID string) error {
	ret := _m.Called(ctx, eventID)

	if len(ret) == 0 {
		panic("no return value specified for ApproveEvent")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, eventID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *EventRepository) CancelLiveEvent(ctx context.Context, eventID string, organizerID string) error {
	ret := _m.Called(ctx, eventID, organizerID)

	if len(ret) == 0 {
		panic("no return value specified for CancelLiveEvent")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, eventID, organizerID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *EventRepository) CreateEvent(ctx context.Context, event *domain.Event, details *domain.EventDetails, tickets []domain.TicketType, personnels []domain.EventPersonnel) error {
	ret := _m.Called(ctx, event, details, tickets, personnels)

	if len(ret) == 0 {
		panic("no return value specified for CreateEvent")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Event, *domain.EventDetails, []domain.TicketType, []domain.EventPersonnel) error); ok {
		r0 = rf(ctx, event, details, tickets, personnels)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *EventRepository) DeleteEvent(ctx context.Context, eventID string) error {
	ret := _m.Called(ctx, eventID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteEvent")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, eventID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *EventRepository) FindPastLiveEvents(ctx context.Context, now time.Time) ([]domain.Event, error) {
	ret := _m.Called(ctx, now)

	if len(ret) == 0 {
		panic("no return value specified for FindPastLiveEvents")
	}

	var r0 []domain.Event
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, time.Time) ([]domain.Event, error)); ok {
		return rf(ctx, now)
	}
	if rf, ok := ret.Get(0).(func(context.Context, time.Time) []domain.Event); ok {
		r0 = rf(ctx, now)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Event)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, time.Time) error); ok {
		r1 = rf(ctx, now)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *EventRepository) GetAdminDashboardStats() (*domain.AdminDashboardStats, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetAdminDashboardStats")
	}

	var r0 *domain.AdminDashboardStats
	var r1 error
	if rf, ok := ret.Get(0).(func() (*domain.AdminDashboardStats, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() *domain.AdminDashboardStats); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.AdminDashboardStats)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *EventRepository) GetAdminRevenueReport(startDate time.Time, endDate time.Time) (*domain.AdminRevenueReport, error) {
	ret := _m.Called(startDate, endDate)

	if len(ret) == 0 {
		panic("no return value specified for GetAdminRevenueReport")
	}

	var r0 *domain.AdminRevenueReport
	var r1 error
	if rf, ok := ret.Get(0).(func(time.Time, time.Time) (*domain.AdminRevenueReport, error)); ok {
		return rf(startDate, endDate)
	}
	if rf, ok := ret.Get(0).(func(time.Time, time.Time) *domain.AdminRevenueReport); ok {
		r0 = rf(startDate, endDate)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.AdminRevenueReport)
		}
	}

	if rf, ok := ret.Get(1).(func(time.Time, time.Time) error); ok {
		r1 = rf(startDate, endDate)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *EventRepository) GetDashboardStats(organizerID string) (*domain.OrganizerDashboardStats, error) {
	ret := _m.Called(organizerID)

	if len(ret) == 0 {
		panic("no return value specified for GetDashboardStats")
	}

	var r0 *domain.OrganizerDashboardStats
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*domain.OrganizerDashboardStats, error)); ok {
		return rf(organizerID)
	}
	if rf, ok := ret.Get(0).(func(string) *domain.OrganizerDashboardStats); ok {
		r0 = rf(organizerID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.OrganizerDashboardStats)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(organizerID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *EventRepository) GetEventByID(eventID string) (*domain.Event, error) {
	ret := _m.Called(eventID)

	if len(ret) == 0 {
		panic("no return value specified for GetEventByID")
	}

	var r0 *domain.Event
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*domain.Event, error)); ok {
		return rf(eventID)
	}
	if rf, ok := ret.Get(0).(func(string) *domain.Event); ok {
		r0 = rf(eventID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Event)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(eventID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *EventRepository) GetEventBySlug(slug string) (*domain.Event, error) {
	ret := _m.Called(slug)

	if len(ret) == 0 {
		panic("no return value specified for GetEventBySlug")
	}

	var r0 *domain.Event
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*domain.Event, error)); ok {
		return rf(slug)
	}
	if rf, ok := ret.Get(0).(func(string) *domain.Event); ok {
		r0 = rf(slug)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Event)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(slug)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *EventRepository) GetEventDetails(eventID string) (*domain.EventDetails, error) {
	ret := _m.Called(eventID)

	if len(ret) == 0 {
		panic("no return value specified for GetEventDetails")
	}

	var r0 *domain.EventDetails
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*domain.EventDetails, error)); ok {
		return rf(eventID)
	}
	if rf, ok := ret.Get(0).(func(string) *domain.EventDetails); ok {
		r0 = rf(eventID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.EventDetails)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(eventID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *EventRepository) GetEventPersonnels(eventID string) ([]domain.EventPersonnel, error) {
	ret := _m.Called(eventID)

	if len(ret) == 0 {
		panic("no return value specified for GetEventPersonnels")
	}

	var r0 []domain.EventPersonnel
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]domain.EventPersonnel, error)); ok {
		return rf(eventID)
	}
	if rf, ok := ret.Get(0).(func(string) []domain.EventPersonnel); ok {
		r0 = rf(eventID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.EventPersonnel)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(eventID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *EventRepository) GetEventsByOrganizerID(organizerID string, status string) ([]domain.Event, error) {
	ret := _m.Called(organizerID, status)

	if len(ret) == 0 {
		panic("no return value specified for GetEventsByOrganizerID")
	}

	var r0 []domain.Event
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) ([]domain.Event, error)); ok {
		return rf(organizerID, status)
	}
	if rf, ok := ret.Get(0).(func(string, string) []domain.Event); ok {
		r0 = rf(organizerID, status)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Event)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(organizerID, status)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *EventRepository) GetSalesReport(organizerID string, eventID string, startDate string, endDate string) (*domain.SalesReportStats, error) {
	ret := _m.Called(organizerID, eventID, startDate, endDate)

	if len(ret) == 0 {
		panic("no return value specified for GetSalesReport")
	}

	var r0 *domain.SalesReportStats
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string, string, string) (*domain.SalesReportStats, error)); ok {
		return rf(organizerID, eventID, startDate, endDate)
	}
	if rf, ok := ret.Get(0).(func(string, string, string, string) *domain.SalesReportStats); ok {
		r0 = rf(organizerID, eventID, startDate, endDate)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.SalesReportStats)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string, string, string) error); ok {
		r1 = rf(organizerID, eventID, startDate, endDate)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *EventRepository) GetTicketTypesByEventID(eventID string) ([]domain.TicketType, error) {
	ret := _m.Called(eventID)

	if len(ret) == 0 {
		panic("no return value specified for GetTicketTypesByEventID")
	}

	var r0 []domain.TicketType
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]domain.TicketType, error)); ok {
		return rf(eventID)
	}
	if rf, ok := ret.Get(0).(func(string) []domain.TicketType); ok {
		r0 = rf(eventID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.TicketType)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(eventID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *EventRepository) ListLiveEvents(city string, category string, search string, sortBy string, minPrice string, maxPrice string, startDate string, endDate string, page int, limit int) ([]domain.Event, int64, float64, float64, error) {
	ret := _m.Called(city, category, search, sortBy, minPrice, maxPrice, startDate, endDate, page, limit)

	if len(ret) == 0 {
		panic("no return value specified for ListLiveEvents")
	}

	var r0 []domain.Event
	var r1 int64
	var r2 float64
	var r3 float64
	var r4 error
	if rf, ok := ret.Get(0).(func(string, string, string, string, string, string, string, string, int, int) ([]domain.Event, int64, float64, float64, error)); ok {
		return rf(city, category, search, sortBy, minPrice, maxPrice, startDate, endDate, page, limit)
	}
	if rf, ok := ret.Get(0).(func(string, string, string, string, string, string, string, string, int, int) []domain.Event); ok {
		r0 = rf(city, category, search, sortBy, minPrice, maxPrice, startDate, endDate, page, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Event)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string, string, string, string, string, string, string, int, int) int64); ok {
		r1 = rf(city, category, search, sortBy, minPrice, maxPrice, startDate, endDate, page, limit)
	} else {
		r1 = ret.Get(1).(int64)
	}

	if rf, ok := ret.Get(2).(func(string, string, string, string, string, string, string, string, int, int) float64); ok {
		r2 = rf(city, category, search, sortBy, minPrice, maxPrice, startDate, endDate, page, limit)
	} else {
		r2 = ret.Get(2).(float64)
	}

	if rf, ok := ret.Get(3).(func(string, string, string, string, string, string, string, string, int, int) float64); ok {
		r3 = rf(city, category, search, sortBy, minPrice, maxPrice, startDate, endDate, page, limit)
	} else {
		r3 = ret.Get(3).(float64)
	}

	if rf, ok := ret.Get(4).(func(string, string, string, string, string, string, string, string, int, int) error); ok {
		r4 = rf(city, category, search, sortBy, minPrice, maxPrice, startDate, endDate, page, limit)
	} else {
		r4 = ret.Error(4)
	}

	return r0, r1, r2, r3, r4
}


func (_m *EventRepository) RejectEvent(ctx context.Context, eventID string, adminID string, reason string) error {
	ret := _m.Called(ctx, eventID, adminID, reason)

	if len(ret) == 0 {
		panic("no return value specified for RejectEvent")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) error); ok {
		r0 = rf(ctx, eventID, adminID, reason)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *EventRepository) RejectEventCancellation(ctx context.Context, eventID string, adminID string, reason string) error {
	ret := _m.Called(ctx, eventID, adminID, reason)

	if len(ret) == 0 {
		panic("no return value specified for RejectEventCancellation")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) error); ok {
		r0 = rf(ctx, eventID, adminID, reason)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *EventRepository) RequestEventCancellation(ctx context.Context, eventID string, organizerID string, reason string) error {
	ret := _m.Called(ctx, eventID, organizerID, reason)

	if len(ret) == 0 {
		panic("no return value specified for RequestEventCancellation")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) error); ok {
		r0 = rf(ctx, eventID, organizerID, reason)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *EventRepository) SettleEventEarnings(ctx context.Context, eventID string, organizerID string, totalAmount float64) error {
	ret := _m.Called(ctx, eventID, organizerID, totalAmount)

	if len(ret) == 0 {
		panic("no return value specified for SettleEventEarnings")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, float64) error); ok {
		r0 = rf(ctx, eventID, organizerID, totalAmount)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *EventRepository) SuspendLiveEvent(ctx context.Context, eventID string, adminID string, reason string) error {
	ret := _m.Called(ctx, eventID, adminID, reason)

	if len(ret) == 0 {
		panic("no return value specified for SuspendLiveEvent")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) error); ok {
		r0 = rf(ctx, eventID, adminID, reason)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *EventRepository) UpdateEvent(ctx context.Context, eventID string, eventUpdates map[string]interface{}, detailUpdates map[string]interface{}, ticketUpdates []domain.TicketType, personnelUpdates []domain.EventPersonnel) error {
	ret := _m.Called(ctx, eventID, eventUpdates, detailUpdates, ticketUpdates, personnelUpdates)

	if len(ret) == 0 {
		panic("no return value specified for UpdateEvent")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, map[string]interface{}, map[string]interface{}, []domain.TicketType, []domain.EventPersonnel) error); ok {
		r0 = rf(ctx, eventID, eventUpdates, detailUpdates, ticketUpdates, personnelUpdates)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


func (_m *EventRepository) UpdateEventStatus(ctx context.Context, eventID string, status string) error {
	ret := _m.Called(ctx, eventID, status)

	if len(ret) == 0 {
		panic("no return value specified for UpdateEventStatus")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, eventID, status)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}



func NewEventRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *EventRepository {
	mock := &EventRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
