package main

import (
	"log"

	"github.com/aswinsreeraj/evntx/internal/infrastructure/database"
	repoImpl "github.com/aswinsreeraj/evntx/internal/infrastructure/repository"
	"github.com/aswinsreeraj/evntx/internal/usecase"
	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.NewPostgresConnection()
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	// Auto migrate (temporary for development)
	db.AutoMigrate(&repoImpl.UserModel{})

	userRepo := repoImpl.NewUserGormRepository(db)
	_ = usecase.NewUserUsecase(userRepo)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	log.Println("Server running on :8080")
	router.Run(":8080")
}
