package main

import (
	"os"
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

	"github.com/aswinsreeraj/evntx/pkg/workers"
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

	db.AutoMigrate(&repoImpl.UserModel{})
	db.AutoMigrate(&repoImpl.OrganizerDetailModel{})
	db.AutoMigrate(&repoImpl.EmailOTPModel{})
	db.AutoMigrate(&repoImpl.UserSessionModel{})
	db.AutoMigrate(&repoImpl.UserRoleModel{})
	db.AutoMigrate(&repoImpl.EventModel{})
	db.AutoMigrate(&repoImpl.EventDetailsModel{})
	db.AutoMigrate(&repoImpl.EventPersonnelModel{})
	db.AutoMigrate(&repoImpl.TicketTypeModel{})
	db.AutoMigrate(&repoImpl.EventModerationLogModel{})
	db.AutoMigrate(&repoImpl.BookingModel{})
	db.AutoMigrate(&repoImpl.BookingTicketModel{})
	db.AutoMigrate(&repoImpl.TicketModel{})

	roleRepo := repoImpl.NewUserRoleGormRepository(db)
	userRepo := repoImpl.NewUserGormRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo, roleRepo)

	bookingRepo := repoImpl.NewBookingGormRepository(db)
	eventRepo := repoImpl.NewEventGormRepository(db)
	bookingUsecase := usecase.NewBookingUsecase(bookingRepo, eventRepo)

	userHandler := httpDelivery.NewUserHandler(userUsecase, bookingUsecase)

	emailSender := emailImpl.NewSMTPSender()

	otpRepo := repoImpl.NewEmailOTPGormRepository(db)
	sessionRepo := repoImpl.NewUserSessionGormRepository(db)
	authUsecase := usecase.NewAuthUsecase(otpRepo, userRepo, sessionRepo, emailSender, roleRepo)
	authHandler := httpDelivery.NewAuthHandler(authUsecase)

	eventUsecase := usecase.NewEventUsecase(eventRepo)
	eventHandler := httpDelivery.NewEventHandler(eventUsecase, userUsecase)
	adminHandler := httpDelivery.NewAdminHandler(eventUsecase, userUsecase)

	bookingHandler := httpDelivery.NewBookingHandler(bookingUsecase)

	expirationWorker := workers.NewBookingExpirationWorker(bookingUsecase)
	go expirationWorker.Start()

	organizerHandler := httpDelivery.NewOrganizerHandler(eventUsecase, userUsecase)

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

	router.Static("/assets", "./assets")

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	router.POST(
		"/auth/otp/request",
		middleware.RateLimitMiddleware(1, 1),
		authHandler.RequestOTP,
	)
	router.POST(
		"/auth/otp/verify",
		middleware.RateLimitMiddleware(5, 5),
		authHandler.VerifyOTP,
	)
	router.POST(
		"/auth/register",
		middleware.RateLimitMiddleware(5, 5),
		authHandler.Register,
	)

	router.POST("/auth/refresh", authHandler.Refresh)
	router.POST("/auth/logout", authHandler.Logout)
	router.POST("/auth/oauth/google", authHandler.GoogleLogin)

	err = os.MkdirAll("assets/images", os.ModePerm)
	if err != nil {
		logger.Log.Warn().Msg("failed to create assets/images directory")
	}
	router.Static("/assets", "./assets")

	// User endpoints
	userGroup := router.Group("/users")
	userGroup.Use(middleware.JWTAuthMiddleware())

	userGroup.GET("/me", userHandler.GetProfile)
	userGroup.GET("/me/bookings", userHandler.GetMyBookingsHandler)
	userGroup.GET("/me/tickets", userHandler.GetMyTicketsHandler)
	userGroup.PUT("/me", userHandler.UpdateProfile)
	userGroup.POST("/me/image", userHandler.UploadProfileImage)

	// Booking endpoint
	bookingGroup := router.Group("/bookings")
	bookingGroup.Use(middleware.JWTAuthMiddleware())
	bookingGroup.POST("/reserve", bookingHandler.ReserveTickets)
	bookingGroup.POST("/:booking_id/cancel", bookingHandler.CancelBooking)

	// Organizer endpoints
	organizerGroup := router.Group("/organizer")
	organizerGroup.Use(middleware.JWTAuthMiddleware())
	organizerGroup.Use(middleware.RBACMiddleware(roleRepo, domain.RoleOrganizer))

	organizerGroup.GET("/me", organizerHandler.GetProfile)
	organizerGroup.POST("/events", organizerHandler.CreateEvent)
	organizerGroup.GET("/events", organizerHandler.GetMyEvents)
	organizerGroup.GET("/events/slug/:slug", organizerHandler.GetEvent)
	organizerGroup.PUT("/events/:event_id", organizerHandler.UpdateEvent)
	organizerGroup.DELETE("/events/:event_id", organizerHandler.DeleteEvent)
	organizerGroup.POST("/events/:event_id/submit", organizerHandler.SubmitEventHandler)
	organizerGroup.POST("/upload", organizerHandler.UploadImage)

	// Admin endpoints
	adminGroup := router.Group("/admin")
	adminGroup.Use(middleware.JWTAuthMiddleware())
	adminGroup.Use(middleware.RBACMiddleware(roleRepo, domain.RoleAdmin))

	adminGroup.GET("/users", userHandler.AdminListUsers)
	adminGroup.GET("/organizers", userHandler.AdminListOrganizers)
	adminGroup.PATCH("/users/:id/status", userHandler.AdminUpdateUserStatus)
	adminGroup.GET("/events", adminHandler.AdminListEvents)
	adminGroup.GET("/events/slug/:slug", adminHandler.AdminGetEvent)
	adminGroup.PATCH("/events/:event_id/approve", adminHandler.ApproveEventHandler)
	adminGroup.PATCH("/events/:event_id/reject", adminHandler.RejectEventHandler)

	adminGroup.GET("/dashboard", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "admin access granted"})
	})

	// Event endpoints
	router.GET("/events", eventHandler.ListEvents)
	router.GET("/events/:slug", eventHandler.GetEvent)

	logger.Log.Info().Msg("Server running on :8080")
	router.Run(":8080")
}
