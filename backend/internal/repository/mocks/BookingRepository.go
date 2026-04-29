package mocks

import (
	context "context"

	domain "github.com/aswinsreeraj/evntx/internal/domain"
	mock "github.com/stretchr/testify/mock"
)

type BookingRepository struct {
	mock.Mock
}

func (_m *BookingRepository) CancelBooking(ctx context.Context, bookingID string, userID string, items []domain.TicketCancelRequest, isRefundable bool) error {
	ret := _m.Called(ctx, bookingID, userID, items, isRefundable)

	if len(ret) == 0 {
		panic("no return value specified for CancelBooking")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, []domain.TicketCancelRequest, bool) error); ok {
		r0 = rf(ctx, bookingID, userID, items, isRefundable)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (_m *BookingRepository) CheckInTicket(ctx context.Context, eventID string, ticketCode string) (*domain.TicketCheckIn, error) {
	ret := _m.Called(ctx, eventID, ticketCode)

	if len(ret) == 0 {
		panic("no return value specified for CheckInTicket")
	}

	var r0 *domain.TicketCheckIn
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (*domain.TicketCheckIn, error)); ok {
		return rf(ctx, eventID, ticketCode)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *domain.TicketCheckIn); ok {
		r0 = rf(ctx, eventID, ticketCode)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.TicketCheckIn)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, eventID, ticketCode)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *BookingRepository) ExpireBookings(ctx context.Context) ([]domain.Booking, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for ExpireBookings")
	}

	var r0 []domain.Booking
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]domain.Booking, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []domain.Booking); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Booking)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *BookingRepository) FindByID(ctx context.Context, bookingID string) (*domain.Booking, error) {
	ret := _m.Called(ctx, bookingID)

	if len(ret) == 0 {
		panic("no return value specified for FindByID")
	}

	var r0 *domain.Booking
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*domain.Booking, error)); ok {
		return rf(ctx, bookingID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *domain.Booking); ok {
		r0 = rf(ctx, bookingID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Booking)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, bookingID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *BookingRepository) GetBookingContextsByIDs(ctx context.Context, bookingIDs []string) (map[string]domain.BookingContextDetails, error) {
	ret := _m.Called(ctx, bookingIDs)

	if len(ret) == 0 {
		panic("no return value specified for GetBookingContextsByIDs")
	}

	var r0 map[string]domain.BookingContextDetails
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []string) (map[string]domain.BookingContextDetails, error)); ok {
		return rf(ctx, bookingIDs)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []string) map[string]domain.BookingContextDetails); ok {
		r0 = rf(ctx, bookingIDs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]domain.BookingContextDetails)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []string) error); ok {
		r1 = rf(ctx, bookingIDs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *BookingRepository) GetPaidBookingsByEventID(ctx context.Context, eventID string) ([]domain.Booking, error) {
	ret := _m.Called(ctx, eventID)

	if len(ret) == 0 {
		panic("no return value specified for GetPaidBookingsByEventID")
	}

	var r0 []domain.Booking
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]domain.Booking, error)); ok {
		return rf(ctx, eventID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []domain.Booking); ok {
		r0 = rf(ctx, eventID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Booking)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, eventID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *BookingRepository) GetTicketCountByBookingID(ctx context.Context, bookingID string) (int, error) {
	ret := _m.Called(ctx, bookingID)

	if len(ret) == 0 {
		panic("no return value specified for GetTicketCountByBookingID")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (int, error)); ok {
		return rf(ctx, bookingID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) int); ok {
		r0 = rf(ctx, bookingID)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, bookingID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *BookingRepository) GetUserBookings(ctx context.Context, userID string, page int, limit int, status string) ([]domain.BookingWithEvent, int64, error) {
	ret := _m.Called(ctx, userID, page, limit, status)

	if len(ret) == 0 {
		panic("no return value specified for GetUserBookings")
	}

	var r0 []domain.BookingWithEvent
	var r1 int64
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string, int, int, string) ([]domain.BookingWithEvent, int64, error)); ok {
		return rf(ctx, userID, page, limit, status)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, int, int, string) []domain.BookingWithEvent); ok {
		r0 = rf(ctx, userID, page, limit, status)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.BookingWithEvent)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, int, int, string) int64); ok {
		r1 = rf(ctx, userID, page, limit, status)
	} else {
		r1 = ret.Get(1).(int64)
	}

	if rf, ok := ret.Get(2).(func(context.Context, string, int, int, string) error); ok {
		r2 = rf(ctx, userID, page, limit, status)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

func (_m *BookingRepository) GetUserTickets(ctx context.Context, userID string, eventID string, bookingID string, status string) ([]domain.TicketWithEvent, error) {
	ret := _m.Called(ctx, userID, eventID, bookingID, status)

	if len(ret) == 0 {
		panic("no return value specified for GetUserTickets")
	}

	var r0 []domain.TicketWithEvent
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, string) ([]domain.TicketWithEvent, error)); ok {
		return rf(ctx, userID, eventID, bookingID, status)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, string) []domain.TicketWithEvent); ok {
		r0 = rf(ctx, userID, eventID, bookingID, status)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.TicketWithEvent)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, string, string) error); ok {
		r1 = rf(ctx, userID, eventID, bookingID, status)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *BookingRepository) PayWithWallet(ctx context.Context, bookingID string, userID string, amount float64) error {
	ret := _m.Called(ctx, bookingID, userID, amount)

	if len(ret) == 0 {
		panic("no return value specified for PayWithWallet")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, float64) error); ok {
		r0 = rf(ctx, bookingID, userID, amount)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (_m *BookingRepository) ReserveTickets(ctx context.Context, booking *domain.Booking, tickets []domain.BookingTicket) error {
	ret := _m.Called(ctx, booking, tickets)

	if len(ret) == 0 {
		panic("no return value specified for ReserveTickets")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Booking, []domain.BookingTicket) error); ok {
		r0 = rf(ctx, booking, tickets)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func NewBookingRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *BookingRepository {
	mock := &BookingRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
