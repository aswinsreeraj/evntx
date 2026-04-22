package testutil

import (
	"github.com/aswinsreeraj/evntx/internal/domain"
	repoImpl "github.com/aswinsreeraj/evntx/internal/infrastructure/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

func SeedUser(db *gorm.DB, email string, role string) domain.User {

	user := repoImpl.UserModel{
		ID:            uuid.NewString(),
		Email:         email,
		Name:          "Test User",
		IsActive:      true,
		EmailVerified: true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	if err := db.Create(&user).Error; err != nil {
		panic(err)
	}

	
	roleMapping := repoImpl.UserRoleModel{
		UserID: user.ID,
		Role:   role,
	}
	if err := db.Create(&roleMapping).Error; err != nil {
		panic(err)
	}

	
	wallet := repoImpl.WalletModel{
		ID:               uuid.NewString(),
		UserID:           user.ID,
		AvailableBalance: 0,
		PendingBalance:   0,
	}
	if err := db.Create(&wallet).Error; err != nil {
		panic(err)
	}

	return domain.User{
		ID:       user.ID,
		Email:    user.Email,
		Name:     user.Name,
		IsActive: user.IsActive,
	}
}

func SeedEvent(db *gorm.DB, organizerID string, title string, status string) domain.Event {
	eventID := uuid.NewString()
	event := repoImpl.EventModel{
		ID:          eventID,
		OrganizerID: organizerID,
		Title:       title,
		Slug:        title + "-slug",
		Status:      status,
		StartTime:   time.Now().Add(24 * time.Hour).Unix(),
		EndTime:     time.Now().Add(48 * time.Hour).Unix(),
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}
	if err := db.Create(&event).Error; err != nil {
		panic(err)
	}

	eventDetails := repoImpl.EventDetailsModel{
		EventID:     eventID,
		Description: "Test event description",
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}
	if err := db.Create(&eventDetails).Error; err != nil {
		panic(err)
	}

	return domain.Event{
		ID:          event.ID,
		OrganizerID: event.OrganizerID,
		Title:       event.Title,
		Slug:        event.Slug,
		Status:      event.Status,
	}
}

func SeedTicketType(db *gorm.DB, eventID string, name string, price float64, totalTickets int) domain.TicketType {
	ticketType := repoImpl.TicketTypeModel{
		ID:                uuid.NewString(),
		EventID:           eventID,
		Name:              name,
		Price:             price,
		TotalQuantity:     totalTickets,
		AvailableQuantity: totalTickets,
		CreatedAt:         time.Now().Unix(),
		UpdatedAt:         time.Now().Unix(),
	}
	if err := db.Create(&ticketType).Error; err != nil {
		panic(err)
	}

	return domain.TicketType{
		ID:                ticketType.ID,
		EventID:           ticketType.EventID,
		Name:              ticketType.Name,
		Price:             ticketType.Price,
		TotalQuantity:     ticketType.TotalQuantity,
		AvailableQuantity: ticketType.AvailableQuantity,
	}
}
