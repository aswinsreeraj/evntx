package http

import (
	"net/http"
	"strings"
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

type updateEventRequest struct {
	Title         *string             `json:"title"`
	City          *string             `json:"city"`
	VenueName     *string             `json:"venue_name"`
	Category      *string             `json:"category"`
	StartTime     *time.Time          `json:"start_time"`
	EndTime       *time.Time          `json:"end_time"`
	Tags          []string            `json:"tags"`
	CoverImageURL *string             `json:"cover_image_url"`
	Details       *detailsUpdateInput `json:"details"`
	TicketTypes   []ticketUpdateInput `json:"ticket_types" binding:"omitempty,dive"`
}

func (h *OrganizerHandler) UpdateEvent(c *gin.Context) {
	organizerID := c.GetString("user_id")
	eventID := c.Param("event_id")

	var req updateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Error().Err(err).Msg("Invalid JSON or validation error")
		response.Error(c, http.StatusBadRequest, apiErrors.InvalidRequestBody, "Invalid request body")
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

	err := h.eventUsecase.UpdateEvent(c.Request.Context(), organizerID, eventID, eventUpdates, detailsUpdates, tickets)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "EVT_004") {
			response.Error(c, http.StatusForbidden, "EVT_004", "Forbidden action")
			return
		} else if strings.Contains(errMsg, "EVT_006") {
			response.Error(c, http.StatusConflict, "EVT_006", "Event cannot be updated in current state")
			return
		}
		// Generic Bad Request for others (e.g. constraints)
		response.Error(c, http.StatusBadRequest, apiErrors.InvalidRequestBody, errMsg)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Event updated successfully",
		"data": gin.H{
			"event_id": eventID,
			"status":   "draft",
		},
	})
}
