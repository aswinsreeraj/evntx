package testutil

import (
	"time"

	"github.com/aswinsreeraj/evntx/internal/cache"
	httpDelivery "github.com/aswinsreeraj/evntx/internal/delivery/http"
	"github.com/aswinsreeraj/evntx/internal/domain"
	repoImpl "github.com/aswinsreeraj/evntx/internal/infrastructure/repository"
	"github.com/aswinsreeraj/evntx/internal/middleware"
	"github.com/aswinsreeraj/evntx/internal/usecase"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MockEmailSender struct{}

func (m *MockEmailSender) SendOTP(to, otp string) error				{ return nil }
func (m *MockEmailSender) SendOrganizerApproval(email, name string) error	{ return nil }

type MockPaymentService struct{}

func (m *MockPaymentService) GetKeyID() string	{ return "test-key" }

func (m *MockPaymentService) CreateOrder(amount int64, receipt string) (*domain.RazorpayOrder, error) {
	return &domain.RazorpayOrder{
		ID:		"order_mock123",
		Amount:		amount,
		Currency:	"INR",
		Receipt:	receipt,
	}, nil
}

func (m *MockPaymentService) VerifySignature(orderID, paymentID, signature string) (bool, error) {
	return true, nil
}

func (m *MockPaymentService) FetchOrder(orderID string) (*domain.RazorpayOrder, error) {
	return &domain.RazorpayOrder{
		ID:		orderID,
		Amount:		50000,
		Currency:	"INR",
		Status:		"paid",
	}, nil
}

func (m *MockPaymentService) RefundPayment(paymentID string, amount int64) error {
	return nil
}

func SetupTestRouter(db *gorm.DB) *gin.Engine {
	emailSender := &MockEmailSender{}
	paymentService := &MockPaymentService{}

	roleRepo := repoImpl.NewUserRoleGormRepository(db)
	userRepo := repoImpl.NewUserGormRepository(db)
	walletRepo := repoImpl.NewWalletGormRepository(db)
	notificationRepo := repoImpl.NewNotificationGormRepository(db)
	platformWalletRepo := repoImpl.NewPlatformWalletGormRepository(db)
	_ = platformWalletRepo.EnsureExists()

	settingsRepo := repoImpl.NewSettingsGormRepository(db)
	_ = settingsRepo.EnsureExists()

	bookingRepo := repoImpl.NewBookingGormRepository(db)
	payoutRepo := repoImpl.NewPayoutGormRepository(db)
	auditRepo := repoImpl.NewAuditGormRepository(db)
	paymentRepo := repoImpl.NewPaymentGormRepository(db, settingsRepo)
	eventRepo := repoImpl.NewEventGormRepository(db)
	engagementRepo := repoImpl.NewEngagementGormRepository(db)
	otpRepo := repoImpl.NewEmailOTPGormRepository(db)
	sessionRepo := repoImpl.NewUserSessionGormRepository(db)

	notificationUsecase := usecase.NewNotificationUsecase(notificationRepo)
	userUsecase := usecase.NewUserUsecase(userRepo, roleRepo, walletRepo, emailSender)
	walletUsecase := usecase.NewWalletUsecase(walletRepo, roleRepo, platformWalletRepo, paymentService, bookingRepo, payoutRepo, notificationUsecase)
	auditUsecase := usecase.NewAuditUsecase(auditRepo, userRepo)
	bookingUsecase := usecase.NewBookingUsecase(bookingRepo, eventRepo, roleRepo, notificationUsecase, settingsRepo)
	paymentUsecase := usecase.NewPaymentUsecase(bookingRepo, eventRepo, paymentRepo, paymentService, notificationUsecase, engagementRepo)
	engagementUsecase := usecase.NewEngagementUsecase(engagementRepo)
	authUsecase := usecase.NewAuthUsecase(otpRepo, userRepo, sessionRepo, emailSender, roleRepo, walletRepo, settingsRepo)
	eventUsecase := usecase.NewEventUsecase(eventRepo, bookingRepo, notificationUsecase, settingsRepo)

	apiCache := cache.NewCache()

	notificationHandler := httpDelivery.NewNotificationHandler(notificationUsecase)
	userHandler := httpDelivery.NewUserHandler(userUsecase, walletUsecase, bookingUsecase, auditUsecase)
	authHandler := httpDelivery.NewAuthHandler(authUsecase)
	eventHandler := httpDelivery.NewEventHandler(eventUsecase, userUsecase, bookingUsecase, apiCache)
	adminHandler := httpDelivery.NewAdminHandler(eventUsecase, userUsecase, walletUsecase, platformWalletRepo, engagementUsecase, settingsRepo, roleRepo, auditUsecase)
	bookingHandler := httpDelivery.NewBookingHandler(bookingUsecase, paymentUsecase)
	paymentHandler := httpDelivery.NewPaymentHandler(paymentUsecase)
	engagementHandler := httpDelivery.NewEngagementHandler(engagementUsecase)
	organizerHandler := httpDelivery.NewOrganizerHandler(eventUsecase, userUsecase, walletUsecase, engagementUsecase, apiCache)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins:		[]string{"*"},
		AllowMethods:		[]string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:		[]string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials:	true,
		MaxAge:			12 * time.Hour,
	}))

	router.Use(middleware.LoggingMiddleware())
	router.Use(gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	router.POST("/auth/otp/request", middleware.RateLimitMiddleware(1, 1), authHandler.RequestOTP)
	router.POST("/auth/otp/verify", middleware.RateLimitMiddleware(5, 5), authHandler.VerifyOTP)
	router.POST("/auth/register", middleware.RateLimitMiddleware(5, 5), authHandler.Register)
	router.POST("/auth/refresh", authHandler.Refresh)
	router.POST("/auth/logout", authHandler.Logout)
	router.POST("/auth/oauth/google", authHandler.GoogleLogin)

	userGroup := router.Group("/users")
	userGroup.Use(middleware.JWTAuthMiddleware())
	userGroup.GET("/me", userHandler.GetProfile)
	userGroup.GET("/me/wallet", userHandler.GetWallet)
	userGroup.POST("/me/wallet/payout", userHandler.RequestPayout)
	userGroup.POST("/me/payout/credentials", userHandler.AddPayoutCredentials)
	userGroup.GET("/me/payouts", userHandler.GetPayouts)
	userGroup.POST("/me/wallet/add-fund", userHandler.CreateAddFundOrder)
	userGroup.POST("/me/wallet/add-fund/verify", userHandler.VerifyAddFundPayment)
	userGroup.GET("/me/wallet/transactions", userHandler.GetWalletTransactions)
	userGroup.GET("/me/bookings", userHandler.GetMyBookingsHandler)
	userGroup.GET("/me/tickets", userHandler.GetMyTicketsHandler)
	userGroup.PUT("/me", userHandler.UpdateProfile)
	userGroup.POST("/me/image", userHandler.UploadProfileImage)

	notificationGroup := router.Group("/notifications")
	notificationGroup.Use(middleware.JWTAuthMiddleware())
	notificationGroup.GET("", notificationHandler.GetNotifications)
	notificationGroup.PATCH("/:id/read", notificationHandler.MarkAsRead)
	notificationGroup.PATCH("/read-all", notificationHandler.MarkAllAsRead)
	notificationGroup.DELETE("", notificationHandler.ClearAll)

	bookingGroup := router.Group("/bookings")
	bookingGroup.Use(middleware.JWTAuthMiddleware())
	bookingGroup.POST("/reserve", bookingHandler.ReserveTickets)
	bookingGroup.POST("/:booking_id/cancel", bookingHandler.CancelBooking)
	bookingGroup.POST("/:booking_id/refund", bookingHandler.RefundBooking)
	bookingGroup.POST("/:booking_id/pay-with-wallet", bookingHandler.PayWithWallet)

	paymentGroup := router.Group("/payments")
	paymentGroup.Use(middleware.JWTAuthMiddleware())
	paymentGroup.POST("/razorpay/order", paymentHandler.CreateRazorpayOrder)
	paymentGroup.POST("/razorpay/verify", paymentHandler.VerifyRazorpayPayment)

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
	organizerGroup.GET("/events", organizerHandler.GetMyEvents)
	organizerGroup.GET("/events/slug/:slug", organizerHandler.GetEvent)
	organizerGroup.PUT("/events/:event_id", organizerHandler.UpdateEvent)
	organizerGroup.DELETE("/events/:event_id", organizerHandler.DeleteEvent)
	organizerGroup.POST("/events/:event_id/cancel-request", organizerHandler.RequestEventCancellation)
	organizerGroup.POST("/events/:event_id/submit", organizerHandler.SubmitEventHandler)
	organizerGroup.POST("/upload", organizerHandler.UploadImage)

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

	router.GET("/events", eventHandler.ListEvents)
	router.GET("/events/:slug", eventHandler.GetEvent)
	router.POST("/events/:event_id/check-in", middleware.JWTAuthMiddleware(), eventHandler.CheckInTicket)

	router.GET("/settings", adminHandler.GetPlatformSettings)
	router.GET("/payment-settings", adminHandler.GetPaymentSettings)

	engagementGroup := router.Group("/engagement")
	engagementGroup.POST("/session", engagementHandler.InitializeSession)
	engagementGroup.POST("/track", engagementHandler.TrackEvent)

	return router
}
