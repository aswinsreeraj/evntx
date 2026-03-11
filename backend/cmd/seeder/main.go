package main

import (
	"fmt"
	"log"
	"time"

	"github.com/aswinsreeraj/evntx/internal/infrastructure/database"
	"github.com/aswinsreeraj/evntx/internal/infrastructure/repository"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found")
	}

	db, err := database.NewPostgresConnection()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("Seeding users...")

	for i := 1; i <= 35; i++ {
		isActive := true
		if i%4 == 0 {
			isActive = false // make 25% of the users suspended
		}

		user := repository.UserModel{
			ID:            uuid.NewString(),
			Name:          fmt.Sprintf("Sample User %d", i),
			Email:         fmt.Sprintf("user%03d@example.com", i),
			Mobile:        fmt.Sprintf("555000%d", i),
			Dob:           "1990-01-01",
			Gender:        "Male",
			ProfileImage:  "",
			Locations:     []string{"Kochi"},
			IsActive:      isActive,
			EmailVerified: true,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		if err := db.Create(&user).Error; err != nil {
			log.Printf("Failed to create user %d: %v", i, err)
		} else {
			fmt.Printf("Created user %d\n", i)
		}
	}

	fmt.Println("Seeding complete! Check admin users list.")
}
