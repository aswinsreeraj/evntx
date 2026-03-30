package http

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aswinsreeraj/evntx/pkg/logger"
	"github.com/google/uuid"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/internal/usecase"
	apiErrors "github.com/aswinsreeraj/evntx/pkg/errors"
	"github.com/aswinsreeraj/evntx/pkg/response"
	"github.com/gin-gonic/gin"
)

type OrganizerHandler struct {
	eventUsecase  *usecase.EventUsecase
	userUsecase   *usecase.UserUsecase
	walletUsecase *usecase.WalletUsecase
}

func NewOrganizerHandler(eu *usecase.EventUsecase, uu *usecase.UserUsecase, wu *usecase.WalletUsecase) *OrganizerHandler {
	return &OrganizerHandler{eventUsecase: eu, userUsecase: uu, walletUsecase: wu}
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
		"total_credited":    wallet.TotalCredited,
		"total_debited":     wallet.TotalDebited,
	})
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

	if err := h.walletUsecase.RequestPayout(userID, req.Amount); err != nil {
		response.AppError(c, err)
		return
	}

	response.Success(c, "Payout request submitted", gin.H{
		"amount": req.Amount,
		"status": "completed",
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

	var personnels []domain.EventPersonnel
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

	dirPath := "assets/events/" + organizerID
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		response.AppError(c, apiErrors.ErrInternalServerError)
		return
	}

	fileID := uuid.NewString()
	ext := filepath.Ext(file.Filename)
	filename := fileID + ext

	filePath := dirPath + "/" + filename
	imageURL := "/" + filePath

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		response.AppError(c, apiErrors.New(500, apiErrors.InternalServerError, "Failed to save image"))
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
		hostName := user.Name
		if organizerDetail != nil && organizerDetail.OrganizationName != "" {
			hostName = organizerDetail.OrganizationName
		}
		host = gin.H{
			"name":   hostName,
			"role":   "Event Organizer",
			"avatar": user.ProfileImage,
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

func (h *OrganizerHandler) CancelLiveEvent(c *gin.Context) {
	organizerID := c.GetString("user_id")
	eventID := c.Param("event_id")

	err := h.eventUsecase.CancelLiveEvent(c.Request.Context(), organizerID, eventID)
	if err != nil {
		response.AppError(c, err)
		return
	}

	response.Success(c, "Live event cancelled successfully. All users refunded.", nil)
}
