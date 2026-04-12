package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aswinsreeraj/evntx/internal/infrastructure/db"
	"github.com/aswinsreeraj/evntx/internal/infrastructure/repository"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load("../.env")

	database, err := db.NewPostgres()
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewEngagementRepository(database)

	// 1. Inspect the table directly
	fmt.Println("--- Database Content (event_engagement_daily) ---")
	var rows []struct {
		EventID string
		Date    time.Time
		PageViews int
		EventViews int
		Visitors int
	}
	database.Table("event_engagement_daily").Find(&rows)
	for _, r := range rows {
		fmt.Printf("Event: %s | Date: %s | PV: %d | EV: %d | Vis: %d\n", 
			r.EventID, r.Date.Format("2006-01-02"), r.PageViews, r.EventViews, r.Visitors)
	}

	if len(rows) == 0 {
		fmt.Println("No rows found in table!")
		return
	}

	// 2. Test the reporting query for the first event found
	fmt.Println("\n--- Testing Report Query ---")
	eventIDs := []string{rows[0].EventID}
	now := time.Now()
	startDate := now.AddDate(0, 0, -30)
	endDate := now

	report, err := repo.GetEngagementReport(context.Background(), eventIDs, startDate, endDate)
	if err != nil {
		fmt.Printf("Error running report: %v\n", err)
		return
	}

	fmt.Printf("Report Results for %s:\n", eventIDs[0])
	fmt.Printf("Page Views: %.0f\n", report.PageViews.Value)
	fmt.Printf("Conversion Rate: %.2f%%\n", report.ConversionRate.Value)
	if len(report.UserJourney) > 0 {
		fmt.Println("User Journey:")
		for _, step := range report.UserJourney {
			fmt.Printf("- %s: %d (%.1f%%)\n", step.Label, step.Count, step.Percentage)
		}
	}
}
