package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/internal/infrastructure/db"
	"github.com/aswinsreeraj/evntx/internal/infrastructure/repository"
	"github.com/google/uuid"
)

func main() {
	database, err := db.NewPostgresDB("postgres://postgres:postgres@localhost:5432/evntx?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewEngagementRepository(database)
	ctx := context.Background()
	testEventID := uuid.NewString()

	fmt.Printf("Starting verification for EventID: %s\n", testEventID)

	
	err = repo.IncrementSuccessfulBookings(ctx, testEventID)
	if err != nil {
		log.Fatalf("Failed to increment: %v", err)
	}
	fmt.Println("Incremented successful bookings.")

	
	loc, _ := time.LoadLocation("Asia/Calcutta")
	now := time.Now().In(loc)
	
	startDate := now
	endDate := now

	
	report, err := repo.GetEngagementReport(ctx, []string{testEventID}, startDate, endDate)
	if err != nil {
		log.Fatalf("Failed to get event report: %v", err)
	}

	foundInEvent := false
	for _, step := range report.UserJourney {
		if step.Label == "Successful Bookings" {
			fmt.Printf("Event Journey: %s = %d\n", step.Label, step.Count)
			if step.Count > 0 {
				foundInEvent = true
			}
		}
	}

	
	adminReport, err := repo.GetEngagementReport(ctx, []string{domain.PlatformEventID}, startDate, endDate)
	if err != nil {
		log.Fatalf("Failed to get admin report: %v", err)
	}

	foundInPlatform := false
	for _, step := range adminReport.UserJourney {
		if step.Label == "Successful Bookings" {
			fmt.Printf("Platform Journey: %s = %d\n", step.Label, step.Count)
			if step.Count > 0 {
				foundInPlatform = true
			}
		}
	}

	if foundInEvent && foundInPlatform {
		fmt.Println("VERIFICATION SUCCESSFUL: Both event and platform metrics incremented in IST.")
	} else {
		fmt.Printf("VERIFICATION FAILED: Event: %v, Platform: %v\n", foundInEvent, foundInPlatform)
	}
}
