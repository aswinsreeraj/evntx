package main

import (
	"fmt"
	"log"
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/internal/infrastructure/database"
	repoImpl "github.com/aswinsreeraj/evntx/internal/infrastructure/repository"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("failed to load env")
	}
	db, err := database.NewPostgresConnection()
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	userRepo := repoImpl.NewUserGormRepository(db)

	email := "test_db_debug@example.com"
	user := &domain.User{
		ID:            uuid.NewString(),
		Email:         email,
		Name:          "Test Debug",
		IsActive:      true,
		EmailVerified: true,
	}

	err = userRepo.Create(user)
	if err != nil {
		log.Fatal("ERROR ON USER CREATE: ", err)
	}
	fmt.Println("USER CREATE SUCCESS")

	sessionRepo := repoImpl.NewUserSessionGormRepository(db)
	session := &domain.UserSession{
		ID:               uuid.NewString(),
		UserID:           user.ID,
		RefreshTokenHash: "hashhash",
		UserAgent:        "fake agent",
		IPAddress:        "127.0.0.1",
		ExpiresAt:        time.Now().Add(7 * 24 * time.Hour),
		Revoked:          false,
	}

	err = sessionRepo.Create(session)
	if err != nil {
		log.Fatal("ERROR ON SESSION CREATE: ", err)
	}
	fmt.Println("SESSION CREATE SUCCESS")
}
