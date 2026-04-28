package main

import (
	"os"
	"time"

	"github.com/aswinsreeraj/evntx/internal/cache"
	httpDelivery "github.com/aswinsreeraj/evntx/internal/delivery/http"
	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/internal/infrastructure/database"
	emailImpl "github.com/aswinsreeraj/evntx/internal/infrastructure/email"
	paymentImpl "github.com/aswinsreeraj/evntx/internal/infrastructure/payment"
	repoImpl "github.com/aswinsreeraj/evntx/internal/infrastructure/repository"
	"github.com/aswinsreeraj/evntx/internal/infrastructure/storage"
	"github.com/aswinsreeraj/evntx/internal/middleware"
	"github.com/aswinsreeraj/evntx/internal/usecase"
	"github.com/aswinsreeraj/evntx/pkg/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/aswinsreeraj/evntx/pkg/workers"
)

func main() {
	logger.Init()
	if err := godotenv.Load(); err != nil {
		logger.Log.Warn().Msg("failed to load env")
	}

	// Initialize S3 storage
	if err := storage.Init(); err != nil {
		logger.Log.Fatal().Msgf("failed to initialize S3 storage: %v", err)
	}

	db, err := database.NewPostgresConnection()
	if err != nil {
		logger.Log.Fatal().Msgf("failed to connect to database: %v", err)
	}

	db.AutoMigrate(&repoImpl.UserModel{})
	db.AutoMigrate(&repoImpl.WalletModel{})
	db.AutoMigrate(&repoImpl.WalletTransactionModel{})
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
	db.AutoMigrate(&repoImpl.PaymentModel{})
	db.AutoMigrate(&repoImpl.NotificationModel{})
	db.AutoMigrate(&repoImpl.PlatformWalletModel{})
	db.AutoMigrate(&repoImpl.PlatformWalletTransactionModel{})
	db.AutoMigrate(&repoImpl.PayoutRequestModel{})
	db.AutoMigrate(&repoImpl.PayoutCredentialModel{})
	db.AutoMigrate(&repoImpl.VisitorSessionModel{})
	db.AutoMigrate(&repoImpl.EngagementEventModel{})
	db.AutoMigrate(&repoImpl.EventEngagementDailyModel{})
	db.AutoMigrate(&repoImpl.PlatformSettingsModel{})
	db.AutoMigrate(&repoImpl.PaymentSettingsModel{})
	db.AutoMigrate(&repoImpl.AuditLogModel{})
	db.AutoMigrate(&repoImpl.JobLogModel{})

	roleRepo := repoImpl.NewUserRoleGormRepository(db)
	userRepo := repoImpl.NewUserGormRepository(db)
	walletRepo := repoImpl.NewWalletGormRepository(db)
	notificationRepo := repoImpl.NewNotificationGormRepository(db)
	platformWalletRepo := repoImpl.NewPlatformWalletGormRepository(db)
	if err := platformWalletRepo.EnsureExists(); err != nil {
		logger.Log.Fatal().Msgf("failed to initialize platform wallet: %v", err)
	}

	settingsRepo := repoImpl.NewSettingsGormRepository(db)
	if err := settingsRepo.EnsureExists(); err != nil {
		logger.Log.Fatal().Msgf("failed to initialize platform settings: %v", err)
	}

	bookingRepo := repoImpl.NewBookingGormRepository(db)
	payoutRepo := repoImpl.NewPayoutGormRepository(db)
	emailSender := emailImpl.NewSMTPSender()

	notificationUsecase := usecase.NewNotificationUsecase(notificationRepo)
	userUsecase := usecase.NewUserUsecase(userRepo, roleRepo, walletRepo, emailSender)
	razorpayService := paymentImpl.NewRazorpayService()
	walletUsecase := usecase.NewWalletUsecase(walletRepo, roleRepo, platformWalletRepo, razorpayService, bookingRepo, payoutRepo, notificationUsecase)

	auditRepo := repoImpl.NewAuditGormRepository(db)
	auditUsecase := usecase.NewAuditUsecase(auditRepo, userRepo)

	paymentRepo := repoImpl.NewPaymentGormRepository(db, settingsRepo)
	eventRepo := repoImpl.NewEventGormRepository(db)
	engagementRepo := repoImpl.NewEngagementGormRepository(db)

	bookingUsecase := usecase.NewBookingUsecase(bookingRepo, eventRepo, roleRepo, notificationUsecase, settingsRepo)
	paymentUsecase := usecase.NewPaymentUsecase(bookingRepo, eventRepo, paymentRepo, razorpayService, notificationUsecase, engagementRepo)
	engagementUsecase := usecase.NewEngagementUsecase(engagementRepo)

	notificationHandler := httpDelivery.NewNotificationHandler(notificationUsecase)
	userHandler := httpDelivery.NewUserHandler(userUsecase, walletUsecase, bookingUsecase, auditUsecase)

	otpRepo := repoImpl.NewEmailOTPGormRepository(db)
	sessionRepo := repoImpl.NewUserSessionGormRepository(db)
	authUsecase := usecase.NewAuthUsecase(otpRepo, userRepo, sessionRepo, emailSender, roleRepo, walletRepo, settingsRepo)
	authHandler := httpDelivery.NewAuthHandler(authUsecase)

	eventUsecase := usecase.NewEventUsecase(eventRepo, bookingRepo, notificationUsecase, settingsRepo)
	apiCache := cache.NewCache()
	eventHandler := httpDelivery.NewEventHandler(eventUsecase, userUsecase, bookingUsecase, apiCache)
	adminHandler := httpDelivery.NewAdminHandler(eventUsecase, userUsecase, walletUsecase, platformWalletRepo, engagementUsecase, settingsRepo, roleRepo, auditUsecase)

	bookingHandler := httpDelivery.NewBookingHandler(bookingUsecase, paymentUsecase)
	paymentHandler := httpDelivery.NewPaymentHandler(paymentUsecase)
	engagementHandler := httpDelivery.NewEngagementHandler(engagementUsecase)

	jobRepo := repoImpl.NewJobGormRepository(db)
	scheduler := workers.NewCronScheduler(jobRepo)

	scheduler.RegisterJob("BookingExpirationJob", "*/1 * * * *", workers.ProcessExpiredBookingsJob(bookingUsecase), 3)
	scheduler.RegisterJob("AutoProcessCompletedEventsJob", "0 * * * *", workers.AutoProcessCompletedEventsJob(eventUsecase), 3)
	scheduler.RegisterJob("ProcessPayoutSettlementsJob", "*/1 * * * *", workers.ProcessPayoutSettlementsJob(walletUsecase), 3)

	scheduler.Start()
	defer scheduler.Stop()

	organizerHandler := httpDelivery.NewOrganizerHandler(eventUsecase, userUsecase, walletUsecase, engagementUsecase, apiCache)

	router := gin.New()
	pprof.Register(router)

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://evntx-self.vercel.app", "http://localhost:5173", "https://evntx.aswinsreeraj.online"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.Use(middleware.LoggingMiddleware())
	router.Use(gin.Recovery())

	router.Static("/assets", "./assets")
	router.Static("/uploads", "./uploads")

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

	// User endpoints
	userGroup := router.Group("/users")
	userGroup.Use(middleware.JWTAuthMiddleware())

	userGroup.GET("/me", userHandler.GetProfile)
	//=== Wallet
	userGroup.GET("/me/wallet", userHandler.GetWallet)
	userGroup.POST("/me/wallet/payout", userHandler.RequestPayout)
	userGroup.POST("/me/payout/credentials", userHandler.AddPayoutCredentials)
	userGroup.GET("/me/payouts", userHandler.GetPayouts)
	userGroup.POST("/me/wallet/add-fund", userHandler.CreateAddFundOrder)
	userGroup.POST("/me/wallet/add-fund/verify", userHandler.VerifyAddFundPayment)
	userGroup.GET("/me/wallet/transactions", userHandler.GetWalletTransactions)
	//=== Booking
	userGroup.GET("/me/bookings", userHandler.GetMyBookingsHandler)
	userGroup.GET("/me/tickets", userHandler.GetMyTicketsHandler)
	//=== Profile
	userGroup.PUT("/me", userHandler.UpdateProfile)
	userGroup.POST("/me/image", userHandler.UploadProfileImage)

	// Notification endpoints
	notificationGroup := router.Group("/notifications")
	notificationGroup.Use(middleware.JWTAuthMiddleware())
	notificationGroup.GET("", notificationHandler.GetNotifications)
	notificationGroup.PATCH("/:id/read", notificationHandler.MarkAsRead)
	notificationGroup.PATCH("/read-all", notificationHandler.MarkAllAsRead)
	notificationGroup.DELETE("", notificationHandler.ClearAll)

	// Booking endpoints
	bookingGroup := router.Group("/bookings")
	bookingGroup.Use(middleware.JWTAuthMiddleware())
	bookingGroup.POST("/reserve", bookingHandler.ReserveTickets)
	bookingGroup.POST("/:booking_id/cancel", bookingHandler.CancelBooking)
	bookingGroup.POST("/:booking_id/refund", bookingHandler.RefundBooking)
	bookingGroup.POST("/:booking_id/pay-with-wallet", bookingHandler.PayWithWallet)

	// Payment endpoints
	paymentGroup := router.Group("/payments")
	paymentGroup.Use(middleware.JWTAuthMiddleware())
	paymentGroup.POST("/razorpay/order", paymentHandler.CreateRazorpayOrder)
	paymentGroup.POST("/razorpay/verify", paymentHandler.VerifyRazorpayPayment)

	// Organizer endpoints
	organizerGroup := router.Group("/organizer")
	organizerGroup.Use(middleware.JWTAuthMiddleware())
	organizerGroup.Use(middleware.RBACMiddleware(roleRepo, domain.RoleOrganizer))

	organizerGroup.GET("/dashboard", organizerHandler.GetDashboard)
	organizerGroup.GET("/reports/sales", organizerHandler.GetSalesReport)
	organizerGroup.GET("/reports/engagement", organizerHandler.GetEngagementReport)
	organizerGroup.GET("/me", organizerHandler.GetProfile)
	organizerGroup.GET("/wallet", organizerHandler.GetWallet)
	organizerGroup.POST("/wallet/payout", organizerHandler.RequestPayout)
	organizerGroup.POST("/payout/credentials", organizerHandler.AddPayoutCredentials)
	organizerGroup.GET("/payouts", organizerHandler.GetPayouts)
	organizerGroup.POST("/events", organizerHandler.CreateEvent)
	//=== Event
	organizerGroup.GET("/events", organizerHandler.GetMyEvents)
	organizerGroup.GET("/events/slug/:slug", organizerHandler.GetEvent)
	organizerGroup.PUT("/events/:event_id", organizerHandler.UpdateEvent)
	organizerGroup.DELETE("/events/:event_id", organizerHandler.DeleteEvent)
	organizerGroup.POST("/events/:event_id/cancel-request", organizerHandler.RequestEventCancellation)
	organizerGroup.POST("/events/:event_id/submit", organizerHandler.SubmitEventHandler)
	organizerGroup.POST("/upload", organizerHandler.UploadImage)

	// Admin endpoints
	adminGroup := router.Group("/admin")
	adminGroup.Use(middleware.JWTAuthMiddleware())
	adminGroup.Use(middleware.RBACMiddleware(roleRepo, domain.RoleAdmin))

	adminGroup.GET("/dashboard", adminHandler.GetAdminDashboard)
	adminGroup.GET("/reports/revenue", adminHandler.GetAdminRevenueReport)
	adminGroup.GET("/reports/engagement", adminHandler.GetAdminEngagementReport)
	adminGroup.GET("/users", userHandler.AdminListUsers)
	adminGroup.GET("/organizers", userHandler.AdminListOrganizers)
	adminGroup.PATCH("/organizers/:id/approve", userHandler.AdminApproveOrganizer)
	adminGroup.PATCH("/organizers/:id/reject", userHandler.AdminRejectOrganizer)
	adminGroup.PATCH("/users/:id/status", userHandler.AdminUpdateUserStatus)
	//=== Event
	adminGroup.GET("/events", adminHandler.AdminListEvents)
	adminGroup.GET("/events/slug/:slug", adminHandler.AdminGetEvent)
	adminGroup.PATCH("/events/:event_id/approve", adminHandler.ApproveEventHandler)
	adminGroup.PATCH("/events/:event_id/reject", adminHandler.RejectEventHandler)
	adminGroup.PATCH("/events/:event_id/suspend", adminHandler.SuspendEventHandler)
	adminGroup.PATCH("/events/:event_id/cancellation/approve", adminHandler.ApproveEventCancellationHandler)
	adminGroup.PATCH("/events/:event_id/cancellation/reject", adminHandler.RejectEventCancellationHandler)
	adminGroup.POST("/events/:event_id/complete", adminHandler.CompleteEventHandler)
	adminGroup.POST("/events/:event_id/settle", adminHandler.SettleEventHandler)
	adminGroup.GET("/platform-wallet", adminHandler.GetPlatformWallet)
	adminGroup.GET("/platform-wallet/transactions", adminHandler.GetPlatformTransactions)

	adminGroup.GET("/payouts", adminHandler.AdminGetPayouts)
	adminGroup.PATCH("/payouts/:id/approve", adminHandler.AdminApprovePayout)
	adminGroup.PATCH("/payouts/:id/reject", adminHandler.AdminRejectPayout)
	adminGroup.POST("/payouts/bulk-approve", adminHandler.AdminBulkApprovePayouts)

	adminGroup.GET("/settings", adminHandler.GetPlatformSettings)
	adminGroup.PUT("/settings", adminHandler.UpdatePlatformSettings)
	adminGroup.GET("/payment-settings", adminHandler.GetPaymentSettings)
	adminGroup.PUT("/payment-settings/:provider", adminHandler.UpdatePaymentProvider)
	adminGroup.GET("/admins", adminHandler.ListAdmins)
	adminGroup.POST("/admins", adminHandler.AddAdmin)
	adminGroup.DELETE("/admins/:id", adminHandler.DeleteAdmin)
	adminGroup.GET("/audit-logs", adminHandler.GetAuditLogs)

	// Event endpoints
	router.GET("/events", eventHandler.ListEvents)
	router.GET("/events/:slug", eventHandler.GetEvent)
	router.POST("/events/:event_id/check-in", middleware.JWTAuthMiddleware(), eventHandler.CheckInTicket)

	// Settings endpoints
	router.GET("/settings", adminHandler.GetPlatformSettings)
	router.GET("/payment-settings", adminHandler.GetPaymentSettings)

	// Engagement endpoints
	engagementGroup := router.Group("/engagement")
	engagementGroup.POST("/session", engagementHandler.InitializeSession)
	engagementGroup.POST("/track", engagementHandler.TrackEvent)

	logger.Log.Info().Msg("Server running on :8080")
	router.Run(":8080")
}
