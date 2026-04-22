package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/internal/repository/mocks"
	"github.com/aswinsreeraj/evntx/internal/usecase"
	apiErrors "github.com/aswinsreeraj/evntx/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBookingUsecase_ReserveTickets(t *testing.T) {
	mockBookingRepo := mocks.NewBookingRepository(t)
	mockEventRepo := mocks.NewEventRepository(t)
	mockSettingsRepo := mocks.NewSettingsRepository(t)

	bookingUsecase := usecase.NewBookingUsecase(
		mockBookingRepo,
		mockEventRepo,
		nil,
		nil,
		mockSettingsRepo,
	)

	t.Run("Success_ReserveTickets", func(t *testing.T) {
		ctx := context.Background()
		userID := "user-123"
		eventID := "event-123"
		requests := []domain.TicketRequest{
			{TicketTypeID: "tt-1", Quantity: 2},
		}

		event := &domain.Event{ID: eventID, Status: "live"}
		ticketTypes := []domain.TicketType{
			{ID: "tt-1", Price: 100, AvailableQuantity: 10},
		}

		mockEventRepo.On("GetEventByID", eventID).Return(event, nil).Once()
		mockEventRepo.On("GetTicketTypesByEventID", eventID).Return(ticketTypes, nil).Once()
		
		
		mockSettingsRepo.On("GetPlatformSettings").Return(&domain.PlatformSettings{
			PlatformFeeType:  domain.PlatformFeeTypeFixed,
			PlatformFeeValue: 10,
		}, nil).Once()

		mockBookingRepo.On("ReserveTickets", ctx, mock.AnythingOfType("*domain.Booking"), mock.AnythingOfType("[]domain.BookingTicket")).Return(nil).Once()

		booking, err := bookingUsecase.ReserveTickets(ctx, userID, eventID, requests)

		assert.NoError(t, err)
		assert.NotNil(t, booking)
		assert.Equal(t, "reserved", booking.Status)
		
		
		
		assert.Equal(t, float64(220), booking.TotalAmount)
	})

	t.Run("Failure_SoldOut", func(t *testing.T) {
		ctx := context.Background()
		userID := "user-123"
		eventID := "event-123"
		requests := []domain.TicketRequest{
			{TicketTypeID: "tt-1", Quantity: 5},
		}

		event := &domain.Event{ID: eventID, Status: "live"}
		ticketTypes := []domain.TicketType{
			{ID: "tt-1", Price: 100, AvailableQuantity: 2}, 
		}

		mockEventRepo.On("GetEventByID", eventID).Return(event, nil).Once()
		mockEventRepo.On("GetTicketTypesByEventID", eventID).Return(ticketTypes, nil).Once()

		booking, err := bookingUsecase.ReserveTickets(ctx, userID, eventID, requests)

		assert.ErrorIs(t, err, apiErrors.ErrTicketSoldOut)
		assert.Nil(t, booking)
	})

	t.Run("Failure_EventNotLive", func(t *testing.T) {
		ctx := context.Background()
		userID := "user-123"
		eventID := "event-123"
		requests := []domain.TicketRequest{
			{TicketTypeID: "tt-1", Quantity: 1},
		}

		event := &domain.Event{ID: eventID, Status: "draft"}

		mockEventRepo.On("GetEventByID", eventID).Return(event, nil).Once()

		booking, err := bookingUsecase.ReserveTickets(ctx, userID, eventID, requests)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "EVT_012: Event not live")
		assert.Nil(t, booking)
	})
}

func TestBookingUsecase_CheckInTicket(t *testing.T) {
	mockBookingRepo := mocks.NewBookingRepository(t)
	mockEventRepo := mocks.NewEventRepository(t)
	mockRoleRepo := mocks.NewUserRoleRepository(t)

	bookingUsecase := usecase.NewBookingUsecase(
		mockBookingRepo,
		mockEventRepo,
		mockRoleRepo,
		nil,
		nil,
	)

	t.Run("Success_CheckIn", func(t *testing.T) {
		ctx := context.Background()
		eventID := "event-123"
		actorID := "org-123"
		ticketCode := "TICKET-123"

		mockRoleRepo.On("GetRolesByUserID", actorID).Return([]domain.UserRole{domain.RoleOrganizer}, nil).Once()
		mockEventRepo.On("GetEventByID", eventID).Return(&domain.Event{ID: eventID, OrganizerID: actorID}, nil).Once()
		
		expectedCheckIn := &domain.TicketCheckIn{
			TicketID: "ticket-1",
			TicketCode: ticketCode,
			Status: "checked-in",
			CheckedInAt: time.Now(),
		}
		mockBookingRepo.On("CheckInTicket", ctx, eventID, ticketCode).Return(expectedCheckIn, nil).Once()

		res, err := bookingUsecase.CheckInTicket(ctx, eventID, actorID, ticketCode)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "checked-in", res.Status)
	})

	t.Run("Failure_Unauthorized", func(t *testing.T) {
		ctx := context.Background()
		eventID := "event-123"
		actorID := "user-456" 
		ticketCode := "TICKET-123"

		mockRoleRepo.On("GetRolesByUserID", actorID).Return([]domain.UserRole{domain.UserRole("user")}, nil).Once()

		res, err := bookingUsecase.CheckInTicket(ctx, eventID, actorID, ticketCode)

		assert.ErrorIs(t, err, apiErrors.ErrForbiddenAction)
		assert.Nil(t, res)
	})
}
