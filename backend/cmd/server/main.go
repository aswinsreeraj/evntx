package main

import (
	"time"

	httpDelivery "github.com/aswinsreeraj/evntx/internal/delivery/http"
	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/internal/infrastructure/database"
	emailImpl "github.com/aswinsreeraj/evntx/internal/infrastructure/email"
	repoImpl "github.com/aswinsreeraj/evntx/internal/infrastructure/repository"
	"github.com/aswinsreeraj/evntx/internal/middleware"
	"github.com/aswinsreeraj/evntx/internal/usecase"
	"github.com/aswinsreeraj/evntx/pkg/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	logger.Init()
	if err := godotenv.Load(); err != nil {
		logger.Log.Warn().Msg("failed to load env")
	}
	db, err := database.NewPostgresConnection()
	if err != nil {
		logger.Log.Fatal().Msgf("failed to connect to database: %v", err)
	}

	// Auto migrate (temporary for development)
	db.AutoMigrate(&repoImpl.UserModel{})
	db.AutoMigrate(&repoImpl.EmailOTPModel{})
	db.AutoMigrate(&repoImpl.UserSessionModel{})
	db.AutoMigrate(&repoImpl.UserRoleModel{})

	userRepo := repoImpl.NewUserGormRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)

	router := gin.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.Use(middleware.LoggingMiddleware())
	router.Use(gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	roleRepo := repoImpl.NewUserRoleGormRepository(db)
	emailSender := emailImpl.NewSMTPSender()

	otpRepo := repoImpl.NewEmailOTPGormRepository(db)
	sessionRepo := repoImpl.NewUserSessionGormRepository(db)
	authUsecase := usecase.NewAuthUsecase(otpRepo, userRepo, sessionRepo, emailSender)
	authHandler := httpDelivery.NewAuthHandler(authUsecase)

	router.POST(
		"/auth/otp/request",
		middleware.RateLimitMiddleware(5, 5), // 5 req/sec burst 5
		authHandler.RequestOTP,
	)
	router.POST(
		"/auth/otp/verify",
		middleware.RateLimitMiddleware(5, 5),
		authHandler.VerifyOTP,
	)

	router.POST("/auth/refresh", authHandler.Refresh)
	router.POST("/auth/logout", authHandler.Logout)
	router.POST("/auth/oauth/google", authHandler.GoogleLogin)

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

	eventRepo := repoImpl.NewEventGormRepository(db)
	eventUsecase := usecase.NewEventUsecase(eventRepo)
	eventHandler := httpDelivery.NewEventHandler(eventUsecase)

	router.GET("/events", eventHandler.ListEvents)
	router.GET("/events/:slug", eventHandler.GetEvent)

	logger.Log.Info().Msg("Server running on :8080")
	router.Run(":8080")
}
