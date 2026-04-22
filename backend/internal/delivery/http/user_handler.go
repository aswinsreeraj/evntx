package http

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/internal/infrastructure/storage"
	"github.com/aswinsreeraj/evntx/internal/usecase"
	apiErrors "github.com/aswinsreeraj/evntx/pkg/errors"
	apiResponse "github.com/aswinsreeraj/evntx/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserHandler struct {
	userUsecase    *usecase.UserUsecase
	walletUsecase  *usecase.WalletUsecase
	bookingUsecase *usecase.BookingUsecase
	auditUsecase   *usecase.AuditUsecase
}

func NewUserHandler(
	userUsecase *usecase.UserUsecase,
	walletUsecase *usecase.WalletUsecase,
	bookingUsecase *usecase.BookingUsecase,
	auditUsecase *usecase.AuditUsecase,
) *UserHandler {
	return &UserHandler{
		userUsecase:    userUsecase,
		walletUsecase:  walletUsecase,
		bookingUsecase: bookingUsecase,
		auditUsecase:   auditUsecase,
	}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")

	user, orgDetail, roles, err := h.userUsecase.GetProfile(userID)
	if err != nil {
		apiResponse.AppError(c, apiErrors.ErrResourceNotFound)
		return
	}

	resp := gin.H{
		"id":            user.ID,
		"name":          user.Name,
		"email":         user.Email,
		"mobile":        user.Mobile,
		"dob":           user.Dob,
		"gender":        user.Gender,
		"profile_image": user.ProfileImage,
		"locations":     user.Locations,
		"roles":         roles,
	}

	if orgDetail != nil {
		resp["organization_name"] = orgDetail.OrganizationName
		resp["address"] = orgDetail.Address
	}

	apiResponse.Success(c, "Profile retrieved successfully", resp)
}

func (h *UserHandler) GetWallet(c *gin.Context) {
	userID := c.GetString("user_id")

	wallet, err := h.walletUsecase.GetWalletByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			apiResponse.AppError(c, apiErrors.ErrResourceNotFound)
			return
		}

		apiResponse.AppError(c, apiErrors.ErrInternalServerError)
		return
	}

	apiResponse.Success(c, "Wallet retrieved successfully", gin.H{
		"available_balance": wallet.AvailableBalance,
		"pending_balance":   wallet.PendingBalance,
		"total_credited":    wallet.TotalCredited,
		"total_debited":     wallet.TotalDebited,
	})
}

func (h *UserHandler) AddPayoutCredentials(c *gin.Context) {
	userID := c.GetString("user_id")

	var req struct {
		AccountHolderName string `json:"account_holder_name" binding:"required"`
		AccountNumber     string `json:"account_number" binding:"required"`
		IFSCCode          string `json:"ifsc_code" binding:"required"`
		UPIID             string `json:"upi_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		apiResponse.AppError(c, apiErrors.ErrInvalidRequestBody)
		return
	}

	if err := h.walletUsecase.AddPayoutCredentials(c.Request.Context(), userID, req.AccountHolderName, req.AccountNumber, req.IFSCCode, req.UPIID); err != nil {
		apiResponse.AppError(c, err)
		return
	}

	apiResponse.Success(c, "Payout credentials saved securely", nil)
}

func (h *UserHandler) RequestPayout(c *gin.Context) {
	userID := c.GetString("user_id")

	var req struct {
		Amount float64 `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		apiResponse.AppError(c, apiErrors.ErrInvalidRequestBody)
		return
	}

	if err := h.walletUsecase.RequestPayout(c.Request.Context(), userID, req.Amount); err != nil {
		apiResponse.AppError(c, err)
		return
	}

	apiResponse.Success(c, "Payout request submitted", gin.H{
		"amount": req.Amount,
		"status": "pending",
	})
}

func (h *UserHandler) GetPayouts(c *gin.Context) {
	userID := c.GetString("user_id")

	payouts, total, err := h.walletUsecase.GetPayoutRequestsByUser(c.Request.Context(), userID, 1, 50)
	if err != nil {
		apiResponse.AppError(c, apiErrors.ErrInternalServerError)
		return
	}

	apiResponse.Success(c, "Payouts retrieved successfully", gin.H{
		"payouts": payouts,
		"total":   total,
	})
}

func (h *UserHandler) CreateAddFundOrder(c *gin.Context) {
	userID := c.GetString("user_id")

	var req struct {
		Amount float64 `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		apiResponse.AppError(c, apiErrors.ErrInvalidRequestBody)
		return
	}

	resp, err := h.walletUsecase.CreateAddFundOrder(userID, req.Amount)
	if err != nil {
		apiResponse.AppError(c, err)
		return
	}

	apiResponse.Success(c, "Add Fund order created", resp)
}

func (h *UserHandler) VerifyAddFundPayment(c *gin.Context) {
	userID := c.GetString("user_id")

	var req struct {
		RazorpayOrderID   string `json:"razorpay_order_id" binding:"required"`
		RazorpayPaymentID string `json:"razorpay_payment_id" binding:"required"`
		RazorpaySignature string `json:"razorpay_signature" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		apiResponse.AppError(c, apiErrors.ErrInvalidRequestBody)
		return
	}

	if err := h.walletUsecase.VerifyAddFundPayment(
		userID,
		req.RazorpayOrderID,
		req.RazorpayPaymentID,
		req.RazorpaySignature,
	); err != nil {
		apiResponse.AppError(c, err)
		return
	}

	apiResponse.Success(c, "Funds added successfully", nil)
}

func (h *UserHandler) GetWalletTransactions(c *gin.Context) {
	userID := c.GetString("user_id")

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	filters := domain.WalletTransactionFilter{
		Type:   c.Query("type"),
		Status: c.Query("status"),
	}

	if filters.Type != "" &&
		filters.Type != domain.WalletTransactionTypeCredit &&
		filters.Type != domain.WalletTransactionTypeDebit {
		apiResponse.AppError(c, apiErrors.New(400, apiErrors.InvalidRequestBody, "Invalid wallet transaction type"))
		return
	}

	if filters.Status != "" &&
		filters.Status != domain.WalletTransactionStatusPending &&
		filters.Status != domain.WalletTransactionStatusCompleted &&
		filters.Status != domain.WalletTransactionStatusFailed {
		apiResponse.AppError(c, apiErrors.New(400, apiErrors.InvalidRequestBody, "Invalid wallet transaction status"))
		return
	}

	transactions, total, err := h.walletUsecase.GetTransactionsByUserID(userID, filters, page, limit)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			apiResponse.AppError(c, apiErrors.ErrResourceNotFound)
			return
		}

		apiResponse.AppError(c, apiErrors.ErrInternalServerError)
		return
	}

	responseTransactions := make([]gin.H, 0, len(transactions))
	for _, txn := range transactions {
		responseTransactions = append(responseTransactions, gin.H{
			"id":             txn.ID,
			"wallet_id":      txn.WalletID,
			"type":           txn.Type,
			"amount":         txn.Amount,
			"reference_type": txn.ReferenceType,
			"reference_id":   txn.ReferenceID,
			"status":         txn.Status,
			"created_at":     txn.CreatedAt,
			"context":        txn.Context,
		})
	}

	apiResponse.Success(c, "Wallet transactions retrieved successfully", gin.H{
		"transactions": responseTransactions,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

type updateProfileRequest struct {
	Name             string   `json:"name" binding:"required"`
	Mobile           string   `json:"mobile"`
	Dob              string   `json:"dob"`
	Gender           string   `json:"gender"`
	OrganizationName string   `json:"organization_name"`
	Address          string   `json:"address"`
	Locations        []string `json:"locations"`
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id")

	var req updateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiResponse.AppError(c, apiErrors.ErrInvalidRequestBody)
		return
	}

	nameRegex := regexp.MustCompile(`^[a-zA-Z\s]+$`)
	if !nameRegex.MatchString(req.Name) {
		apiResponse.AppError(c, apiErrors.New(400, apiErrors.InvalidRequestBody, "Name can only contain alphabets and spaces"))
		return
	}

	if err := h.userUsecase.UpdateProfile(userID, req.Name, req.Mobile, req.Dob, req.Gender, req.OrganizationName, req.Address, req.Locations); err != nil {
		apiResponse.AppError(c, apiErrors.ErrInternalServerError)
		return
	}

	apiResponse.Success(c, "Profile updated successfully", nil)
}

func (h *UserHandler) AdminListUsers(c *gin.Context) {
	search := c.Query("search")
	status := c.Query("status")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "5")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	users, total, err := h.userUsecase.AdminSearchUsers(search, status, page, limit)
	if err != nil {
		apiResponse.AppError(c, apiErrors.ErrInternalServerError)
		return
	}

	resp := make([]gin.H, 0, len(users))
	for _, u := range users {
		resp = append(resp, gin.H{
			"id":             u.ID,
			"name":           u.Name,
			"email":          u.Email,
			"is_active":      u.IsActive,
			"total_bookings": u.TotalBookings,
			"wallet_balance": u.WalletBalance,
			"created_at":     u.CreatedAt,
		})
	}

	apiResponse.Success(c, "Users retrieved successfully", gin.H{
		"users": resp,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

type updateStatusRequest struct {
	IsActive bool `json:"is_active"`
}

func (h *UserHandler) AdminUpdateUserStatus(c *gin.Context) {
	userID := c.Param("id")

	var req updateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiResponse.AppError(c, apiErrors.ErrInvalidRequestBody)
		return
	}

	if err := h.userUsecase.AdminUpdateUserStatus(userID, req.IsActive); err != nil {
		apiResponse.AppError(c, apiErrors.ErrInternalServerError)
		return
	}

	if h.auditUsecase != nil {
		adminID := c.GetString("user_id")
		clientIP := c.ClientIP()
		statusText := "suspended"
		if req.IsActive {
			statusText = "changed to active"
		}
		
		
		
		actionText := "User #" + userID[:6] + " " + statusText
		
		h.auditUsecase.LogAction(adminID, actionText, domain.ActionTagUser, map[string]interface{}{
			"user_id": userID,
			"status": req.IsActive,
		}, clientIP)
	}

	apiResponse.Success(c, "User status updated successfully", nil)
}

func (h *UserHandler) AdminListOrganizers(c *gin.Context) {
	search := c.Query("search")
	status := c.Query("status")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "5")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	orgs, total, err := h.userUsecase.AdminSearchOrganizers(search, status, page, limit)
	if err != nil {
		apiResponse.AppError(c, apiErrors.ErrInternalServerError)
		return
	}

	resp := make([]gin.H, 0, len(orgs))
	for _, u := range orgs {
		resp = append(resp, gin.H{
			"id":                      u.ID,
			"name":                    u.Name,
			"email":                   u.Email,
			"is_active":               u.IsActive,
			"approval_status":         u.ApprovalStatus,
			"total_bookings":          u.TotalBookings,
			"total_events":            u.TotalEvents,
			"wallet_balance":          u.WalletBalance,
			"total_revenue_generated": u.TotalRevenue,
		})
	}

	apiResponse.Success(c, "Organizers retrieved successfully", gin.H{
		"organizers": resp,
		"pagination": gin.H{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

func (h *UserHandler) AdminApproveOrganizer(c *gin.Context) {
	organizerID := c.Param("id")
	if err := h.userUsecase.AdminApproveOrganizer(organizerID); err != nil {
		apiResponse.AppError(c, err)
		return
	}
	apiResponse.Success(c, "Organizer approved successfully", nil)
}

func (h *UserHandler) AdminRejectOrganizer(c *gin.Context) {
	organizerID := c.Param("id")
	if err := h.userUsecase.AdminRejectOrganizer(organizerID); err != nil {
		apiResponse.AppError(c, err)
		return
	}
	apiResponse.Success(c, "Organizer rejected successfully", nil)
}

func (h *UserHandler) GetMyBookingsHandler(c *gin.Context) {
	userID := c.GetString("user_id")

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	status := c.Query("status")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	bookings, total, err := h.bookingUsecase.GetUserBookings(c.Request.Context(), userID, page, limit, status)
	if err != nil {
		apiResponse.AppError(c, apiErrors.ErrInternalServerError)
		return
	}

	responseBookings := make([]map[string]interface{}, 0, len(bookings))
	for _, b := range bookings {
		responseBookings = append(responseBookings, map[string]interface{}{
			"booking_id":       b.BookingID,
			"event_id":         b.EventID,
			"event_title":      b.EventTitle,
			"event_city":       b.EventCity,
			"event_start_time": b.EventStartTime,
			"status":           b.Status,
			"total_amount":     b.TotalAmount,
			"ticket_count":     b.TicketCount,
			"created_at":       b.CreatedAt,
			"coverImageUrl":    b.CoverImageURL,
			"venue":            b.VenueName,
			"tags":             strings.Split(b.Tags, ","),
			"event_status":     b.EventStatus,
		})
	}

	apiResponse.Success(c, "Bookings fetched successfully", gin.H{
		"bookings": responseBookings,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (h *UserHandler) GetMyTicketsHandler(c *gin.Context) {
	userID := c.GetString("user_id")

	eventID := c.Query("event_id")
	bookingID := c.Query("booking_id")
	status := c.Query("status")

	tickets, err := h.bookingUsecase.GetUserTickets(c.Request.Context(), userID, eventID, bookingID, status)
	if err != nil {
		apiResponse.AppError(c, apiErrors.ErrInternalServerError)
		return
	}

	responseTickets := make([]map[string]interface{}, 0, len(tickets))
	for _, t := range tickets {
		responseTickets = append(responseTickets, map[string]interface{}{
			"ticket_id":     t.TicketID,
			"ticket_code":   t.TicketCode,
			"event_id":      t.EventID,
			"event_title":   t.EventTitle,
			"ticket_type":   t.TicketType,
			"status":        t.Status,
			"checked_in_at": t.CheckedInAt,
		})
	}

	apiResponse.Success(c, "Tickets fetched successfully", gin.H{
		"tickets": responseTickets,
	})
}

func (h *UserHandler) UploadProfileImage(c *gin.Context) {
	userID := c.GetString("user_id")

	file, err := c.FormFile("profile_image")
	if err != nil {
		apiResponse.AppError(c, apiErrors.New(400, apiErrors.InvalidRequestBody, "Image file is required"))
		return
	}

	src, err := file.Open()
	if err != nil {
		apiResponse.AppError(c, apiErrors.ErrInternalServerError)
		return
	}
	defer src.Close()

	ext := filepath.Ext(file.Filename)
	key := fmt.Sprintf("profiles/%s/%s%s", userID, uuid.NewString(), ext)
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg"
	}

	imageURL, err := storage.UploadFile(c.Request.Context(), key, contentType, src)
	if err != nil {
		apiResponse.AppError(c, apiErrors.New(500, apiErrors.InternalServerError, "Failed to upload image"))
		return
	}

	if err := h.userUsecase.UploadProfileImage(userID, imageURL); err != nil {
		apiResponse.AppError(c, apiErrors.New(500, apiErrors.InternalServerError, "Failed to update profile image"))
		return
	}

	apiResponse.Success(c, "Profile image uploaded successfully", gin.H{
		"profile_image": imageURL,
	})
}
