package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/internal/infrastructure/database"
	repoImpl "github.com/aswinsreeraj/evntx/internal/infrastructure/repository"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Note: .env file not found, using environment variables")
	}

	db, err := database.NewPostgresConnection()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	repo := repoImpl.NewEngagementGormRepository(db)
	ctx := context.Background()

	eventID := "125d2ef8-cbcc-4f2b-be97-917abd80024f" // Example Event ID
	sessionID := "test-session-" + fmt.Sprintf("%d", time.Now().Unix())

	fmt.Println("--- TEST 1: Multiple Hits Same Session ---")
	for i := 1; i <= 3; i++ {
		evt := &domain.EngagementEvent{
			ID:        fmt.Sprintf("evt-s1-%d", i),
			SessionID: sessionID,
			EventID:   &eventID,
			EventType: domain.EngagementEventEventView,
			CreatedAt: time.Now(),
		}
		repo.LogEvent(ctx, evt)
		fmt.Printf("Logged Hit %d\n", i)
	}

	fmt.Println("\n--- TEST 2: New Session Hit ---")
	newSessionID := "test-session-new-" + fmt.Sprintf("%d", time.Now().Unix())
	evt := &domain.EngagementEvent{
		ID:        "evt-s2-1",
		SessionID: newSessionID,
		EventID:   &eventID,
		EventType: domain.EngagementEventEventView,
		CreatedAt: time.Now(),
	}
	repo.LogEvent(ctx, evt)
	fmt.Println("Logged Hit from New Session")

	fmt.Println("\n--- VERIFYING RESULTS ---")
	
	printStats := func(title string, id string) {
		stats, _ := repo.GetDailyAggregates(ctx, id, time.Now().AddDate(0, 0, -1), time.Now().AddDate(0, 0, 1))
		fmt.Printf("\n[%s] (%s):\n", title, id)
		for _, s := range stats {
			if s.Date.Format("2006-01-02") == time.Now().Format("2006-01-02") {
				fmt.Printf("  Date: %s, Visitors: %d, PageViews: %d, EventViews: %d\n", 
					s.Date.Format("2006-01-02"), s.Visitors, s.PageViews, s.EventViews)
			}
		}
	}

	printStats("Target Event", eventID)
	printStats("Global Platform", domain.PlatformEventID)
}
