package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/internal/infrastructure/database"
	"github.com/aswinsreeraj/evntx/internal/infrastructure/repository"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: go run cmd/admin/main.go <email>\nExample: go run cmd/admin/main.go admin@evntx.com")
	}
	email := os.Args[1]

	
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found or failed to load. Using system environment variables.")
	}

	db, err := database.NewPostgresConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	
	db.AutoMigrate(&repository.UserModel{}, &repository.UserRoleModel{})

	var user repository.UserModel
	res := db.Where("email = ?", email).First(&user)

	if res.Error != nil {
		
		fmt.Printf("User %s not found. Creating...\n", email)
		user = repository.UserModel{
			ID:            uuid.NewString(),
			Email:         email,
			Name:          "Admin User",
			IsActive:      true,
			EmailVerified: true,
		}
		if err := db.Create(&user).Error; err != nil {
			log.Fatalf("Failed to create user: %v", err)
		}
	} else {
		fmt.Printf("User %s found (ID: %s).\n", email, user.ID)
	}

	
	var role repository.UserRoleModel
	res = db.Where("user_id = ? AND role = ?", user.ID, domain.RoleAdmin).First(&role)

	if res.Error != nil {
		fmt.Printf("Granting admin role to %s...\n", email)
		role = repository.UserRoleModel{
			UserID: user.ID,
			Role:   string(domain.RoleAdmin),
		}
		if err := db.Create(&role).Error; err != nil {
			log.Fatalf("Failed to grant admin role: %v", err)
		}
	} else {
		fmt.Printf("User %s already has the admin role.\n", email)
	}

	fmt.Printf("\nSuccess! %s is now an admin.\n", email)
}
