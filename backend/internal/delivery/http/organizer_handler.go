package http

import (
	"net/http"
	"time"

	"github.com/aswinsreeraj/evntx/pkg/logger"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/internal/usecase"
	apiErrors "github.com/aswinsreeraj/evntx/pkg/errors"
	"github.com/aswinsreeraj/evntx/pkg/response"
	"github.com/gin-gonic/gin"
)

type OrganizerHandler struct {
	eventUsecase *usecase.EventUsecase
}

func NewOrganizerHandler(u *usecase.EventUsecase) *OrganizerHandler {
	return &OrganizerHandler{eventUsecase: u}
}

type detailsInput struct {
	Description        string  `json:"description" binding:"required"`
	VenueAddress       string  `json:"venue_address" binding:"required"`
	MapURL             string  `json:"map_url"`
	TotalCapacity      int     `json:"total_capacity" binding:"required,gt=0"`
	TermsAndConditions string  `json:"terms_and_conditions"`
}

type ticketInput struct {
	Name          string  `json:"name" binding:"required"`
	Price         float64 `json:"price" binding:"gte=0"`
	TotalQuantity int     `json:"total_quantity" binding:"required,gt=0"`
}

type createEventRequest struct {
	Title         string         `json:"title" binding:"required"`
	City          string         `json:"city" binding:"required"`
	VenueName     string         `json:"venue_name" binding:"required"`
	Category      string         `json:"category"`
	StartTime     time.Time      `json:"start_time" binding:"required"`
	EndTime       time.Time      `json:"end_time" binding:"required"`
	Tags          []string       `json:"tags"`
	CoverImageURL string         `json:"cover_image_url"`
	Details       detailsInput   `json:"details" binding:"required"`
	TicketTypes   []ticketInput  `json:"ticket_types" binding:"required,min=1,dive"`
}

func (h *OrganizerHandler) CreateEvent(c *gin.Context) {
	organizerID := c.GetString("user_id")

	var req createEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Error().Err(err).Msg("Invalid JSON or validation error")
		response.Error(c, http.StatusBadRequest, apiErrors.InvalidRequestBody, "Invalid request body")
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

	eventID, err := h.eventUsecase.CreateEvent(c.Request.Context(), organizerID, event, details, tickets)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apiErrors.InvalidStateTransition, err.Error())
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
