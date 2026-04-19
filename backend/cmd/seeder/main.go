package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/internal/infrastructure/database"
	infraRepo "github.com/aswinsreeraj/evntx/internal/infrastructure/repository"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type seedEvent struct {
	Title       string
	City        string
	VenueName   string
	Category    string
	Tags        string
	CoverImage  string
	Description string
	Address     string
	StartTime   time.Time
	EndTime     time.Time
	Tickets     []domain.TicketType
	Personnels  []domain.EventPersonnel
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found")
	}

	db, err := database.NewPostgresConnection()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := db.AutoMigrate(
		&infraRepo.UserModel{},
		&infraRepo.WalletModel{},
		&infraRepo.WalletTransactionModel{},
		&infraRepo.OrganizerDetailModel{},
		&infraRepo.UserRoleModel{},
		&infraRepo.EventModel{},
		&infraRepo.EventDetailsModel{},
		&infraRepo.EventPersonnelModel{},
		&infraRepo.TicketTypeModel{},
	); err != nil {
		log.Fatal("Failed to run seeder migrations:", err)
	}

	now := time.Now()
	organizer := ensureOrganizer(db, now)
	seedGoers(db, now)
	seedEventsForOrganizer(db, organizer.ID, now)

	fmt.Println("Seeder completed successfully.")
}

func ensureOrganizer(db *gorm.DB, now time.Time) infraRepo.UserModel {
	var organizer infraRepo.UserModel
	err := db.Where("email = ?", "test@organizer.com").First(&organizer).Error

	if err == nil {
		log.Println("Sample organizer already exists")
	} else if err == gorm.ErrRecordNotFound {
		organizer = infraRepo.UserModel{
			ID:            uuid.NewString(),
			Name:          "Sample Organizer",
			Email:         "test@organizer.com",
			Mobile:        "9876543210",
			Dob:           "1992-06-12",
			Gender:        "Female",
			ProfileImage:  "",
			Locations:     []string{"Kochi", "Bangalore", "Pune"},
			IsActive:      true,
			EmailVerified: true,
			CreatedAt:     now,
			UpdatedAt:     now,
		}

		if err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&organizer).Error; err != nil {
				return err
			}

			wallet := infraRepo.WalletModel{
				ID:     uuid.NewString(),
				UserID: organizer.ID,
			}

			return tx.Create(&wallet).Error
		}); err != nil {
			log.Fatalf("Failed to create sample organizer: %v", err)
		}
		log.Println("Created sample organizer test@organizer.com")
	} else {
		log.Fatalf("Failed to query organizer: %v", err)
	}

	role := infraRepo.UserRoleModel{
		UserID: organizer.ID,
		Role:   string(domain.RoleOrganizer),
	}
	if err := db.Save(&role).Error; err != nil {
		log.Fatalf("Failed to assign organizer role: %v", err)
	}

	detail := infraRepo.OrganizerDetailModel{
		UserID:           organizer.ID,
		OrganizationName: "EVNTX Live",
		Address:          "55 Marine Drive, Kochi",
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	if err := db.Where("user_id = ?", organizer.ID).Assign(detail).FirstOrCreate(&detail).Error; err != nil {
		log.Fatalf("Failed to create organizer details: %v", err)
	}

	return organizer
}

func seedGoers(db *gorm.DB, now time.Time) {
	for i := 1; i <= 20; i++ {
		email := fmt.Sprintf("user%03d@example.com", i)
		var user infraRepo.UserModel
		err := db.Where("email = ?", email).First(&user).Error

		if err == gorm.ErrRecordNotFound {
			user = infraRepo.UserModel{
				ID:            uuid.NewString(),
				Name:          fmt.Sprintf("Sample User %d", i),
				Email:         email,
				Mobile:        fmt.Sprintf("555000%03d", i),
				Dob:           "1990-01-01",
				Gender:        "Male",
				Locations:     []string{"Kochi"},
				IsActive:      i%4 != 0,
				EmailVerified: true,
				CreatedAt:     now,
				UpdatedAt:     now,
			}

			if err := db.Transaction(func(tx *gorm.DB) error {
				if err := tx.Create(&user).Error; err != nil {
					return err
				}

				wallet := infraRepo.WalletModel{
					ID:     uuid.NewString(),
					UserID: user.ID,
				}

				return tx.Create(&wallet).Error
			}); err != nil {
				log.Printf("Failed to create user %s: %v", email, err)
				continue
			}
		} else if err != nil {
			log.Printf("Failed to query user %s: %v", email, err)
		}
	}
}

func seedEventsForOrganizer(db *gorm.DB, organizerID string, now time.Time) {
	events := []seedEvent{
		{
			Title:       "Friday Night at Vapour Ladies Night",
			City:        "Kochi",
			VenueName:   "JLN Stadium",
			Category:    "Music",
			Tags:        "Music,Party,Dance",
			CoverImage:  "/assets/images/badass-bollywood.png",
			Description: "A high-energy Bollywood party night with guest performers, curated experiences, and a late evening dance floor.",
			Address:     "JLN Stadium, Kochi",
			StartTime:   now.AddDate(0, 0, 7).Truncate(time.Hour).Add(19 * time.Hour),
			EndTime:     now.AddDate(0, 0, 8).Truncate(time.Hour),
			Tickets: []domain.TicketType{
				{Name: "Standard", Price: 5000, TotalQuantity: 120},
				{Name: "Premium", Price: 7500, TotalQuantity: 60},
				{Name: "VIP Access", Price: 10000, TotalQuantity: 30},
			},
			Personnels: []domain.EventPersonnel{
				{Name: "Jane Doe", Role: "Host", Image: "/assets/images/host.jpg"},
				{Name: "DJ Jazee", Role: "Professional DJ", Image: "/assets/images/dj.jpg"},
			},
		},
		{
			Title:       "Sand Castle Workshop",
			City:        "Pune",
			VenueName:   "XYZ Conference Hall",
			Category:    "Art",
			Tags:        "Art,Workshop",
			CoverImage:  "/assets/images/sand-castle.png",
			Description: "A hands-on creative workshop that guides participants through sculpting, texture work, and collaborative sand installations.",
			Address:     "XYZ Conference Hall, Pune",
			StartTime:   now.AddDate(0, 0, 10).Truncate(time.Hour).Add(12 * time.Hour),
			EndTime:     now.AddDate(0, 0, 10).Truncate(time.Hour).Add(15 * time.Hour),
			Tickets: []domain.TicketType{
				{Name: "Standard", Price: 5000, TotalQuantity: 80},
				{Name: "Premium", Price: 7500, TotalQuantity: 25},
			},
			Personnels: []domain.EventPersonnel{
				{Name: "Asha Menon", Role: "Workshop Lead", Image: "/assets/images/host.jpg"},
			},
		},
		{
			Title:       "Premium Roy by Shreya",
			City:        "Chennai",
			VenueName:   "ABC Cafe",
			Category:    "Comedy",
			Tags:        "Comedy,Live Show",
			CoverImage:  "/assets/images/premium-roy.png",
			Description: "A stand-up comedy special with a tightly written set, crowd-work moments, and a premium lounge-style venue setup.",
			Address:     "ABC Cafe, Chennai",
			StartTime:   now.AddDate(0, 0, 14).Truncate(time.Hour).Add(18 * time.Hour),
			EndTime:     now.AddDate(0, 0, 14).Truncate(time.Hour).Add(20 * time.Hour),
			Tickets: []domain.TicketType{
				{Name: "Standard", Price: 5000, TotalQuantity: 100},
				{Name: "Premium", Price: 7500, TotalQuantity: 40},
				{Name: "VIP Access", Price: 10000, TotalQuantity: 20},
			},
			Personnels: []domain.EventPersonnel{
				{Name: "Shreya", Role: "Lead Comedian", Image: "/assets/images/perfomer.jpg"},
			},
		},
		{
			Title:       "Scorpions Coming Home Live 2026",
			City:        "Pune",
			VenueName:   "MNO Stadium",
			Category:    "Music",
			Tags:        "Music,Live Show",
			CoverImage:  "/assets/images/scorpions.png",
			Description: "A stadium concert experience with legacy rock hits, premium viewing zones, and crowd-scale production.",
			Address:     "MNO Stadium, Pune",
			StartTime:   now.AddDate(0, 0, 18).Truncate(time.Hour).Add(20 * time.Hour),
			EndTime:     now.AddDate(0, 0, 18).Truncate(time.Hour).Add(23 * time.Hour),
			Tickets: []domain.TicketType{
				{Name: "Standard", Price: 5000, TotalQuantity: 200},
				{Name: "Premium", Price: 7500, TotalQuantity: 80},
				{Name: "VIP Access", Price: 10000, TotalQuantity: 40},
			},
			Personnels: []domain.EventPersonnel{
				{Name: "Scorpions", Role: "Headline Act", Image: "/assets/images/perfomer.jpg"},
			},
		},
		{
			Title:       "Advancing Passive Fire Protection",
			City:        "Pune",
			VenueName:   "XYZ Conference Hall",
			Category:    "Business",
			Tags:        "Art,Workshop",
			CoverImage:  "/assets/images/fire-protect.png",
			Description: "An industry conference on passive fire protection systems, case studies, and practical implementation strategies.",
			Address:     "XYZ Conference Hall, Pune",
			StartTime:   now.AddDate(0, 0, 22).Truncate(time.Hour).Add(12 * time.Hour),
			EndTime:     now.AddDate(0, 0, 22).Truncate(time.Hour).Add(17 * time.Hour),
			Tickets: []domain.TicketType{
				{Name: "Standard", Price: 2200, TotalQuantity: 150},
			},
			Personnels: []domain.EventPersonnel{
				{Name: "Joe Smith", Role: "Lead Speaker", Image: "/assets/images/perfomer.jpg"},
			},
		},
		{
			Title:       "If I'm Not Wrong By Tarang Hardikar",
			City:        "Chennai",
			VenueName:   "ABC Cafe",
			Category:    "Comedy",
			Tags:        "Comedy,Live Show",
			CoverImage:  "/assets/images/if-im-not-wrong.png",
			Description: "A sharp stand-up set built around observational comedy, long-form stories, and intimate club energy.",
			Address:     "ABC Cafe, Chennai",
			StartTime:   now.AddDate(0, 0, 26).Truncate(time.Hour).Add(18 * time.Hour),
			EndTime:     now.AddDate(0, 0, 26).Truncate(time.Hour).Add(20 * time.Hour),
			Tickets: []domain.TicketType{
				{Name: "Standard", Price: 1500, TotalQuantity: 120},
			},
			Personnels: []domain.EventPersonnel{
				{Name: "Tarang Hardikar", Role: "Performer", Image: "/assets/images/perfomer.jpg"},
			},
		},
	}

	for _, event := range events {
		if err := upsertSeedEvent(db, organizerID, event, now); err != nil {
			log.Printf("Failed to seed event %s: %v", event.Title, err)
		}
	}
}

func upsertSeedEvent(db *gorm.DB, organizerID string, event seedEvent, now time.Time) error {
	slug := generateSlug(event.Title)

	return db.Transaction(func(tx *gorm.DB) error {
		var existing infraRepo.EventModel
		if err := tx.Where("slug = ?", slug).First(&existing).Error; err == nil {
			if err := tx.Where("event_id = ?", existing.ID).Delete(&infraRepo.EventDetailsModel{}).Error; err != nil {
				return err
			}
			if err := tx.Where("event_id = ?", existing.ID).Delete(&infraRepo.EventPersonnelModel{}).Error; err != nil {
				return err
			}
			if err := tx.Where("event_id = ?", existing.ID).Delete(&infraRepo.TicketTypeModel{}).Error; err != nil {
				return err
			}
			if err := tx.Delete(&existing).Error; err != nil {
				return err
			}
		} else if err != gorm.ErrRecordNotFound {
			return err
		}

		eventID := uuid.NewString()
		eventModel := infraRepo.EventModel{
			ID:            eventID,
			OrganizerID:   organizerID,
			Title:         event.Title,
			Slug:          slug,
			City:          event.City,
			VenueName:     event.VenueName,
			Category:      event.Category,
			StartTime:     event.StartTime.Unix(),
			EndTime:       event.EndTime.Unix(),
			Tags:          event.Tags,
			Status:        "live",
			CoverImageURL: event.CoverImage,
		}
		if err := tx.Create(&eventModel).Error; err != nil {
			return err
		}

		details := infraRepo.EventDetailsModel{
			EventID:            eventID,
			Description:        event.Description,
			VenueAddress:       event.Address,
			MapURL:             "",
			TotalCapacity:      totalCapacity(event.Tickets),
			AvailableCapacity:  totalCapacity(event.Tickets),
			Rating:             4.7,
			TotalReviews:       128,
			TermsAndConditions: "Tickets are non-transferable. Please carry a valid ID proof at entry.",
		}
		if err := tx.Create(&details).Error; err != nil {
			return err
		}

		for _, ticket := range event.Tickets {
			model := infraRepo.TicketTypeModel{
				ID:                uuid.NewString(),
				EventID:           eventID,
				Name:              ticket.Name,
				Price:             ticket.Price,
				TotalQuantity:     ticket.TotalQuantity,
				AvailableQuantity: ticket.TotalQuantity,
				Version:           1,
				CreatedAt:         now.Unix(),
				UpdatedAt:         now.Unix(),
			}
			if err := tx.Create(&model).Error; err != nil {
				return err
			}
		}

		for _, person := range event.Personnels {
			model := infraRepo.EventPersonnelModel{
				ID:          uuid.NewString(),
				EventID:     eventID,
				Name:        person.Name,
				Role:        person.Role,
				Image:       person.Image,
				ProfileLink: "",
			}
			if err := tx.Create(&model).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func totalCapacity(tickets []domain.TicketType) int {
	total := 0
	for _, ticket := range tickets {
		total += ticket.TotalQuantity
	}
	return total
}

func generateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "'", "")
	return slug
}
