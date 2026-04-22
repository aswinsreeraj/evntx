package http

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/aswinsreeraj/evntx/internal/cache"
	"github.com/aswinsreeraj/evntx/internal/infrastructure/storage"
	"github.com/aswinsreeraj/evntx/pkg/logger"
	"github.com/google/uuid"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/internal/usecase"
	apiErrors "github.com/aswinsreeraj/evntx/pkg/errors"
	"github.com/aswinsreeraj/evntx/pkg/response"
	"github.com/gin-gonic/gin"
)

type OrganizerHandler struct {
	eventUsecase      *usecase.EventUsecase
	userUsecase       *usecase.UserUsecase
	walletUsecase     *usecase.WalletUsecase
	engagementUsecase *usecase.EngagementUsecase
	cache             *cache.Cache
}

func NewOrganizerHandler(eu *usecase.EventUsecase, uu *usecase.UserUsecase, wu *usecase.WalletUsecase, engUsecase *usecase.EngagementUsecase, c *cache.Cache) *OrganizerHandler {
	return &OrganizerHandler{eventUsecase: eu, userUsecase: uu, walletUsecase: wu, engagementUsecase: engUsecase, cache: c}
}

func (h *OrganizerHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")

	user, organizerDetail, _, err := h.userUsecase.GetProfile(userID)
	if err != nil {
		response.AppError(c, apiErrors.ErrResourceNotFound)
		return
	}

	orgName := ""
	address := ""
	if organizerDetail != nil {
		orgName = organizerDetail.OrganizationName
		address = organizerDetail.Address
	}

	response.Success(c, "Organizer profile retrieved successfully", gin.H{
		"id":                user.ID,
		"name":              user.Name,
		"email":             user.Email,
		"mobile":            user.Mobile,
		"dob":               user.Dob,
		"gender":            user.Gender,
		"profile_image":     user.ProfileImage,
		"locations":         user.Locations,
		"organization_name": orgName,
		"address":           address,
	})
}

func (h *OrganizerHandler) GetDashboard(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	stats, err := h.eventUsecase.GetOrganizerDashboardStats(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve dashboard stats"})
		return
	}

	c.JSON(200, gin.H{"data": stats})
}

func (h *OrganizerHandler) GetSalesReport(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	eventID := c.Query("event_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	stats, err := h.eventUsecase.GetSalesReport(userID, eventID, startDate, endDate)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve sales report stats"})
		return
	}

	c.JSON(200, gin.H{"data": stats})
}

func (h *OrganizerHandler) GetEngagementReport(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	eventIDParam := c.Query("event_id")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	
	events, err := h.eventUsecase.GetOrganizerEvents(c.Request.Context(), userID, "")
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve organizer events"})
		return
	}

	var eventIDs []string
	for _, e := range events {
		if eventIDParam == "" || eventIDParam == "all" || eventIDParam == e.ID || eventIDParam == e.Slug {
			eventIDs = append(eventIDs, e.ID)
		}
	}

	var startDate, endDate time.Time
	if startDateStr != "" {
		startDate, _ = time.Parse(time.RFC3339, startDateStr)
	}
	if endDateStr != "" {
		endDate, _ = time.Parse(time.RFC3339, endDateStr)
	}
	loc, _ := time.LoadLocation("Asia/Calcutta")
	if startDate.IsZero() {
		endDate = time.Now().In(loc)
		startDate = endDate.AddDate(0, 0, -30)
	}

	report, err := h.engagementUsecase.GetEngagementReport(c.Request.Context(), eventIDs, startDate, endDate)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve engagement report"})
		return
	}

	c.JSON(200, gin.H{"data": report})
}

func (h *OrganizerHandler) GetWallet(c *gin.Context) {
	userID := c.GetString("user_id")

	wallet, err := h.walletUsecase.GetWalletByUserID(userID)
	if err != nil {
		response.AppError(c, apiErrors.ErrResourceNotFound)
		return
	}

	response.Success(c, "Organizer wallet retrieved successfully", gin.H{
		"available_balance": wallet.AvailableBalance,
		"pending_balance":   wallet.PendingBalance,
		"reserve_balance":   wallet.ReserveBalance,
		"total_credited":    wallet.TotalCredited,
		"total_debited":     wallet.TotalDebited,
	})
}

func (h *OrganizerHandler) AddPayoutCredentials(c *gin.Context) {
	userID := c.GetString("user_id")

	var req struct {
		AccountHolderName string `json:"account_holder_name" binding:"required"`
		AccountNumber     string `json:"account_number" binding:"required"`
		IFSCCode          string `json:"ifsc_code" binding:"required"`
		UPIID             string `json:"upi_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.AppError(c, apiErrors.ErrInvalidRequestBody)
		return
	}

	if err := h.walletUsecase.AddPayoutCredentials(c.Request.Context(), userID, req.AccountHolderName, req.AccountNumber, req.IFSCCode, req.UPIID); err != nil {
		response.AppError(c, err)
		return
	}

	response.Success(c, "Payout credentials saved securely", nil)
}

func (h *OrganizerHandler) RequestPayout(c *gin.Context) {
	userID := c.GetString("user_id")

	var req struct {
		Amount float64 `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.AppError(c, apiErrors.ErrInvalidRequestBody)
		return
	}

	if err := h.walletUsecase.RequestPayout(c.Request.Context(), userID, req.Amount); err != nil {
		response.AppError(c, err)
		return
	}

	response.Success(c, "Payout request submitted", gin.H{
		"amount": req.Amount,
		"status": "pending",
	})
}

func (h *OrganizerHandler) GetPayouts(c *gin.Context) {
	userID := c.GetString("user_id")

	payouts, total, err := h.walletUsecase.GetPayoutRequestsByUser(c.Request.Context(), userID, 1, 50)
	if err != nil {
		response.AppError(c, apiErrors.ErrInternalServerError)
		return
	}

	response.Success(c, "Payouts retrieved successfully", gin.H{
		"payouts": payouts,
		"total":   total,
	})
}

type detailsInput struct {
	Description        string `json:"description" binding:"required"`
	VenueAddress       string `json:"venue_address" binding:"required"`
	MapURL             string `json:"map_url"`
	TotalCapacity      int    `json:"total_capacity" binding:"required,gt=0"`
	TermsAndConditions string `json:"terms_and_conditions"`
}

type ticketInput struct {
	Name          string  `json:"name" binding:"required"`
	Price         float64 `json:"price" binding:"gte=0"`
	TotalQuantity int     `json:"total_quantity" binding:"required,gt=0"`
}

type personnelInput struct {
	Name        string `json:"name" binding:"required"`
	Role        string `json:"role" binding:"required"`
	Image       string `json:"image"`
	ProfileLink string `json:"profile_link"`
}

type createEventRequest struct {
	Title         string           `json:"title" binding:"required"`
	City          string           `json:"city" binding:"required"`
	VenueName     string           `json:"venue_name" binding:"required"`
	Category      string           `json:"category"`
	StartTime     time.Time        `json:"start_time" binding:"required"`
	EndTime       time.Time        `json:"end_time"`
	Tags          []string         `json:"tags"`
	CoverImageURL string           `json:"cover_image_url"`
	Details       detailsInput     `json:"details" binding:"required"`
	TicketTypes   []ticketInput    `json:"ticket_types" binding:"required,min=1,dive"`
	KeyPersonnel  []personnelInput `json:"key_personnel" binding:"omitempty,dive"`
}

func (h *OrganizerHandler) CreateEvent(c *gin.Context) {
	organizerID := c.GetString("user_id")

	var req createEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Error().Err(err).Msg("Invalid JSON or validation error")
		response.AppError(c, apiErrors.ErrInvalidRequestBody)
		return
	}

	tagsStr := ""
	for i, tag := range req.Tags {
		if i > 0 {
			tagsStr += ","
		}
		tagsStr += tag
	}

	event := &domain.Event{
		Title:         req.Title,
		City:          req.City,
		VenueName:     req.VenueName,
		Category:      req.Category,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		Tags:          tagsStr,
		CoverImageURL: req.CoverImageURL,
	}

	details := &domain.EventDetails{
		Description:        req.Details.Description,
		VenueAddress:       req.Details.VenueAddress,
		MapURL:             req.Details.MapURL,
		TotalCapacity:      req.Details.TotalCapacity,
		TermsAndConditions: req.Details.TermsAndConditions,
	}

	var tickets []domain.TicketType
	for _, t := range req.TicketTypes {
		tickets = append(tickets, domain.TicketType{
			Name:          t.Name,
			Price:         t.Price,
			TotalQuantity: t.TotalQuantity,
		})
	}

	var personnels []domain.EventPersonnel
	for _, p := range req.KeyPersonnel {
		personnels = append(personnels, domain.EventPersonnel{
			Name:        p.Name,
			Role:        p.Role,
			Image:       p.Image,
			ProfileLink: p.ProfileLink,
		})
	}

	eventID, err := h.eventUsecase.CreateEvent(c.Request.Context(), organizerID, event, details, tickets, personnels)
	if err != nil {
		response.AppError(c, apiErrors.New(400, apiErrors.InvalidStateTransition, err.Error()))
		fmt.Println("Error here", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Event created successfully",
		"data": gin.H{
			"event_id": eventID,
			"status":   "draft",
		},
	})
}

type detailsUpdateInput struct {
	Description        *string `json:"description"`
	VenueAddress       *string `json:"venue_address"`
	MapURL             *string `json:"map_url"`
	TotalCapacity      *int    `json:"total_capacity" binding:"omitempty,gt=0"`
	TermsAndConditions *string `json:"terms_and_conditions"`
}

type ticketUpdateInput struct {
	ID            *string  `json:"id"`
	Name          *string  `json:"name"`
	Price         *float64 `json:"price" binding:"omitempty,gte=0"`
	TotalQuantity *int     `json:"total_quantity" binding:"omitempty,gt=0"`
}

type personnelUpdateInput struct {
	ID          *string `json:"id"`
	Name        *string `json:"name"`
	Role        *string `json:"role"`
	Image       *string `json:"image"`
	ProfileLink *string `json:"profile_link"`
}

type updateEventRequest struct {
	Title         *string                `json:"title"`
	City          *string                `json:"city"`
	VenueName     *string                `json:"venue_name"`
	Category      *string                `json:"category"`
	StartTime     *time.Time             `json:"start_time"`
	EndTime       *time.Time             `json:"end_time"`
	Tags          []string               `json:"tags"`
	CoverImageURL *string                `json:"cover_image_url"`
	Details       *detailsUpdateInput    `json:"details"`
	TicketTypes   []ticketUpdateInput    `json:"ticket_types" binding:"omitempty,dive"`
	KeyPersonnel  []personnelUpdateInput `json:"key_personnel" binding:"omitempty,dive"`
}

func (h *OrganizerHandler) UpdateEvent(c *gin.Context) {
	organizerID := c.GetString("user_id")
	eventID := c.Param("event_id")

	var req updateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Error().Err(err).Msg("Invalid JSON or validation error")
		response.AppError(c, apiErrors.ErrInvalidRequestBody)
		return
	}

	eventUpdates := make(map[string]interface{})
	detailsUpdates := make(map[string]interface{})

	if req.Title != nil {
		eventUpdates["title"] = *req.Title
		eventUpdates["slug"] = strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(*req.Title, " ", "-"), "'", ""))
	}
	if req.City != nil {
		eventUpdates["city"] = *req.City
	}
	if req.VenueName != nil {
		eventUpdates["venue_name"] = *req.VenueName
	}
	if req.Category != nil {
		eventUpdates["category"] = *req.Category
	}
	if req.StartTime != nil {
		eventUpdates["start_time"] = req.StartTime.Unix()
	}
	if req.EndTime != nil {
		eventUpdates["end_time"] = req.EndTime.Unix()
	}
	if req.Tags != nil {
		tagsStr := ""
		for i, tag := range req.Tags {
			if i > 0 {
				tagsStr += ","
			}
			tagsStr += tag
		}
		eventUpdates["tags"] = tagsStr
	}
	if req.CoverImageURL != nil {
		eventUpdates["cover_image_url"] = *req.CoverImageURL
	}

	if req.Details != nil {
		if req.Details.Description != nil {
			detailsUpdates["description"] = *req.Details.Description
		}
		if req.Details.VenueAddress != nil {
			detailsUpdates["venue_address"] = *req.Details.VenueAddress
		}
		if req.Details.MapURL != nil {
			detailsUpdates["map_url"] = *req.Details.MapURL
		}
		if req.Details.TotalCapacity != nil {
			detailsUpdates["total_capacity"] = *req.Details.TotalCapacity
		}
		if req.Details.TermsAndConditions != nil {
			detailsUpdates["terms_and_conditions"] = *req.Details.TermsAndConditions
		}
	}

	var tickets []domain.TicketType
	if req.TicketTypes != nil {
		tickets = make([]domain.TicketType, 0, len(req.TicketTypes))
		for _, t := range req.TicketTypes {
			ticket := domain.TicketType{}
			if t.ID != nil {
				ticket.ID = *t.ID
			}
			if t.Name != nil {
				ticket.Name = *t.Name
			}
			if t.Price != nil {
				ticket.Price = *t.Price
			}
			if t.TotalQuantity != nil {
				ticket.TotalQuantity = *t.TotalQuantity
			}
			tickets = append(tickets, ticket)
		}
	}

	var personnels []domain.EventPersonnel
	if req.KeyPersonnel != nil {
		personnels = make([]domain.EventPersonnel, 0, len(req.KeyPersonnel))
		for _, p := range req.KeyPersonnel {
			personnel := domain.EventPersonnel{}
			if p.ID != nil {
				personnel.ID = *p.ID
			}
			if p.Name != nil {
				personnel.Name = *p.Name
			}
			if p.Role != nil {
				personnel.Role = *p.Role
			}
			if p.Image != nil {
				personnel.Image = *p.Image
			}
			if p.ProfileLink != nil {
				personnel.ProfileLink = *p.ProfileLink
			}
			personnels = append(personnels, personnel)
		}
	}

	err := h.eventUsecase.UpdateEvent(c.Request.Context(), organizerID, eventID, eventUpdates, detailsUpdates, tickets, personnels)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "EVT_004") {
			response.AppError(c, apiErrors.New(403, apiErrors.ForbiddenAction, "Forbidden action"))
			return
		} else if strings.Contains(errMsg, "EVT_006") {
			response.AppError(c, apiErrors.New(409, apiErrors.InvalidStateTransition, "Event cannot be updated in current state"))
			return
		}
		response.AppError(c, apiErrors.New(400, apiErrors.InvalidRequestBody, errMsg))
		return
	}

	if h.cache != nil {
		if cachedEvent, slugErr := h.eventUsecase.GetEventByID(eventID); slugErr == nil && cachedEvent != nil {
			h.cache.Delete("event:" + cachedEvent.Slug)
		}
	}

	response.Success(c, "Event updated successfully", gin.H{
		"event_id": eventID,
		"status":   "draft",
	})
}

func (h *OrganizerHandler) SubmitEventHandler(c *gin.Context) {
	organizerID := c.GetString("user_id")
	eventID := c.Param("event_id")

	err := h.eventUsecase.SubmitEventForApproval(c.Request.Context(), organizerID, eventID)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "EVT_004") {
			response.AppError(c, apiErrors.New(403, apiErrors.ForbiddenAction, "Forbidden action"))
			return
		} else if strings.Contains(errMsg, "EVT_006") {
			response.AppError(c, apiErrors.New(409, apiErrors.InvalidStateTransition, "Event cannot be submitted in current state"))
			return
		}
		response.AppError(c, apiErrors.New(400, apiErrors.InvalidRequestBody, errMsg))
		return
	}

	response.Success(c, "Event submitted for approval", gin.H{
		"event_id": eventID,
		"status":   "pending",
	})
}

func (h *OrganizerHandler) UploadImage(c *gin.Context) {
	organizerID := c.GetString("user_id")

	file, err := c.FormFile("image")
	if err != nil {
		response.AppError(c, apiErrors.New(400, apiErrors.InvalidRequestBody, "Image file is required"))
		return
	}

	src, err := file.Open()
	if err != nil {
		response.AppError(c, apiErrors.ErrInternalServerError)
		return
	}
	defer src.Close()

	ext := filepath.Ext(file.Filename)
	key := fmt.Sprintf("events/%s/%s%s", organizerID, uuid.NewString(), ext)
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg"
	}

	imageURL, err := storage.UploadFile(c.Request.Context(), key, contentType, src)
	if err != nil {
		response.AppError(c, apiErrors.New(500, apiErrors.InternalServerError, "Failed to upload image"))
		return
	}

	response.Success(c, "Image uploaded successfully", gin.H{
		"url": imageURL,
	})
}

func (h *OrganizerHandler) GetMyEvents(c *gin.Context) {
	organizerID := c.GetString("user_id")
	status := c.Query("status")

	events, err := h.eventUsecase.GetOrganizerEvents(c.Request.Context(), organizerID, status)
	if err != nil {
		response.AppError(c, apiErrors.ErrInternalServerError)
		return
	}

	response.Success(c, "Events fetched successfully", gin.H{
		"events": events,
	})
}

func (h *OrganizerHandler) GetEvent(c *gin.Context) {
	organizerID := c.GetString("user_id")
	slug := c.Param("slug")

	event, details, personnels, tickets, err := h.eventUsecase.GetOrganizerEvent(slug, organizerID)
	if err != nil {
		response.AppError(c, apiErrors.ErrResourceNotFound)
		return
	}

	host := gin.H(nil)
	user, organizerDetail, _, userErr := h.userUsecase.GetProfile(organizerID)
	if userErr == nil && user != nil {
		orgName := ""
		if organizerDetail != nil && organizerDetail.OrganizationName != "" {
			orgName = organizerDetail.OrganizationName
		}
		host = gin.H{
			"name":         user.Name,
			"organization": orgName,
			"role":         "Event Organizer",
			"avatar":       user.ProfileImage,
			"email":        user.Email,
			"mobile":       user.Mobile,
			"address":      "",
		}
		if organizerDetail != nil {
			host["address"] = organizerDetail.Address
		}
	}

	response.Success(c, "Event fetched successfully", gin.H{
		"event":        event,
		"details":      details,
		"personnels":   personnels,
		"ticket_types": tickets,
		"host":         host,
	})
}

func (h *OrganizerHandler) DeleteEvent(c *gin.Context) {
	organizerID := c.GetString("user_id")
	eventID := c.Param("event_id")

	err := h.eventUsecase.DeleteEvent(c.Request.Context(), organizerID, eventID)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "EVT_004") {
			response.AppError(c, apiErrors.New(403, apiErrors.ForbiddenAction, "Forbidden action"))
			return
		}
		response.AppError(c, apiErrors.New(400, apiErrors.InvalidRequestBody, errMsg))
		return
	}

	response.Success(c, "Event deleted successfully", nil)
}

func (h *OrganizerHandler) RequestEventCancellation(c *gin.Context) {
	organizerID := c.GetString("user_id")
	eventID := c.Param("event_id")
	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.AppError(c, apiErrors.ErrInvalidRequestBody)
		return
	}

	err := h.eventUsecase.RequestEventCancellation(c.Request.Context(), organizerID, eventID, req.Reason)
	if err != nil {
		response.AppError(c, err)
		return
	}

	response.Success(c, "Event cancellation request submitted for admin approval", nil)
}
