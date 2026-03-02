package main

import (
	"log"

	httpDelivery "github.com/aswinsreeraj/evntx/internal/delivery/http"
	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/internal/infrastructure/database"
	repoImpl "github.com/aswinsreeraj/evntx/internal/infrastructure/repository"
	"github.com/aswinsreeraj/evntx/internal/middleware"
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
	db.AutoMigrate(&repoImpl.EmailOTPModel{})
	db.AutoMigrate(&repoImpl.UserSessionModel{})
	db.AutoMigrate(&repoImpl.UserRoleModel{})

	userRepo := repoImpl.NewUserGormRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	roleRepo := repoImpl.NewUserRoleGormRepository(db)

	otpRepo := repoImpl.NewEmailOTPGormRepository(db)
	sessionRepo := repoImpl.NewUserSessionGormRepository(db)
	authUsecase := usecase.NewAuthUsecase(otpRepo, userRepo, sessionRepo)
	authHandler := httpDelivery.NewAuthHandler(authUsecase)

	router.POST("/auth/otp/request", authHandler.RequestOTP)
	router.POST("/auth/otp/verify", authHandler.VerifyOTP)

	router.POST("/auth/refresh", authHandler.Refresh)
	router.POST("/auth/logout", authHandler.Logout)

	protected := router.Group("/admin")
	protected.Use(middleware.JWTAuthMiddleware())
	protected.Use(middleware.RBACMiddleware(roleRepo, domain.RoleAdmin))

	protected.GET("/dashboard", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "admin access granted"})
	})

	userHandler := httpDelivery.NewUserHandler(userUsecase)

	userGroup := router.Group("/users")
	userGroup.Use(middleware.JWTAuthMiddleware())

	userGroup.GET("/me", userHandler.GetProfile)
	userGroup.PUT("/me", userHandler.UpdateProfile)

	adminGroup := router.Group("/admin")
	adminGroup.Use(middleware.JWTAuthMiddleware())
	adminGroup.Use(middleware.RBACMiddleware(roleRepo, domain.RoleAdmin))

	adminGroup.GET("/users", userHandler.AdminListUsers)
	adminGroup.PATCH("/users/:id/status", userHandler.AdminUpdateUserStatus)

	log.Println("Server running on :8080")
	router.Run(":8080")
}
