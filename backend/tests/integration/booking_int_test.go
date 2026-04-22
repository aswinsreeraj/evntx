package integration

import (
	"context"
	"testing"

	"github.com/aswinsreeraj/evntx/internal/domain"
	repoImpl "github.com/aswinsreeraj/evntx/internal/infrastructure/repository"
	"github.com/aswinsreeraj/evntx/internal/usecase"
	"github.com/aswinsreeraj/evntx/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestBookingIntegration_ReserveTickets(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.ClearDatabase(db)

	bookingRepo := repoImpl.NewBookingGormRepository(db)
	eventRepo := repoImpl.NewEventGormRepository(db)
	roleRepo := repoImpl.NewUserRoleGormRepository(db)
	notificationRepo := repoImpl.NewNotificationGormRepository(db)
	settingsRepo := repoImpl.NewSettingsGormRepository(db)
	_ = settingsRepo.EnsureExists()

	notificationUsecase := usecase.NewNotificationUsecase(notificationRepo)
	bookingUsecase := usecase.NewBookingUsecase(bookingRepo, eventRepo, roleRepo, notificationUsecase, settingsRepo)

	
	organizer := testutil.SeedUser(db, "organizer@evntx.com", string(domain.RoleOrganizer))
	user := testutil.SeedUser(db, "booker@evntx.com", "user")
	event := testutil.SeedEvent(db, organizer.ID, "Tech Conference", "live")
	ticketType := testutil.SeedTicketType(db, event.ID, "General Admission", 100.0, 50)

	
	ticketsToReserve := []domain.TicketRequest{
		{
			TicketTypeID: ticketType.ID,
			Quantity:     2,
		},
	}

	booking, err := bookingUsecase.ReserveTickets(context.Background(), user.ID, event.ID, ticketsToReserve)
	assert.NoError(t, err)
	assert.NotNil(t, booking)
	assert.Equal(t, "reserved", booking.Status)
	assert.Equal(t, 260.0, booking.TotalAmount) 

	
	var dbTicketType repoImpl.TicketTypeModel
	err = db.First(&dbTicketType, "id = ?", ticketType.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, 48, dbTicketType.AvailableQuantity)

	
	largeRequest := []domain.TicketRequest{
		{
			TicketTypeID: ticketType.ID,
			Quantity:     50, 
		},
	}
	_, err = bookingUsecase.ReserveTickets(context.Background(), user.ID, event.ID, largeRequest)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "sold out")

	
	err = db.First(&dbTicketType, "id = ?", ticketType.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, 48, dbTicketType.AvailableQuantity)
}
