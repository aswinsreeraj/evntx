package main

import (
	"fmt"
	"log"

	"github.com/aswinsreeraj/evntx/internal/infrastructure/database"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	db, err := database.NewPostgresConnection()
	if err != nil {
		log.Fatal(err)
	}

	var events []map[string]interface{}
	db.Table("engagement_events").Order("created_at desc").Limit(10).Find(&events)

	fmt.Println("Recent Engagement Events:")
	for _, e := range events {
		fmt.Printf("Type: %v, EventID: %v, SessionID: %v\n", e["event_type"], e["event_id"], e["session_id"])
	}

	var daily []map[string]interface{}
	db.Table("event_engagement_daily").Limit(10).Find(&daily)
	fmt.Println("\nDaily Aggregates:")
	for _, d := range daily {
		fmt.Printf("EventID: %v, Date: %v, Visitors: %v, Views: %v\n", d["event_id"], d["date"], d["visitors"], d["event_views"])
	}
}
