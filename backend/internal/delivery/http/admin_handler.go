package http

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/aswinsreeraj/evntx/internal/usecase"
	"github.com/aswinsreeraj/evntx/pkg/errors"
	"github.com/aswinsreeraj/evntx/pkg/response"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	eventUsecase *usecase.EventUsecase
}

func NewAdminHandler(eventUsecase *usecase.EventUsecase) *AdminHandler {
	return &AdminHandler{
		eventUsecase: eventUsecase,
	}
}

func (h *AdminHandler) AdminListEvents(c *gin.Context) {
	search := c.Query("search")
	status := c.Query("status")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	// Convert page to integer, default to 1 if empty/invalid
	page := 1
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}

	// Convert limit to integer, default to 10 if empty/invalid
	limit := 10
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	events, total, err := h.eventUsecase.AdminSearchEvents(search, status, page, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, errors.InternalServerError, "Failed to retrieve events")
		return
	}

	eventResponses := make([]gin.H, 0, len(events))
	for _, evt := range events {
		eventResponses = append(eventResponses, gin.H{
			"id":             evt.ID,
			"title":          evt.Title,
			"slug":           evt.Slug,
			"city":           evt.City,
			"venue_name":     evt.VenueName,
			"category":       evt.Category,
			"start_time":     evt.StartTime,
			"end_time":       evt.EndTime,
			"tags":           evt.Tags,
			"status":         evt.Status,
			"organizer_name": evt.OrganizerName,
			"tickets_sold":   evt.TicketsSold,
			"revenue":        evt.Revenue,
			"created_at":     evt.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"events": eventResponses,
			"pagination": gin.H{
				"total": total,
				"page":  page,
				"limit": limit,
			},
		},
	})
}

func (h *AdminHandler) ApproveEventHandler(c *gin.Context) {
	adminID := c.GetString("user_id")
	eventID := c.Param("event_id")

	err := h.eventUsecase.ApproveEvent(c.Request.Context(), adminID, eventID)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "EVT_004") {
			response.Error(c, http.StatusForbidden, "EVT_004", "Insufficient admin privileges")
			return
		} else if strings.Contains(errMsg, "EVT_006") {
			response.Error(c, http.StatusConflict, "EVT_006", "Event cannot be approved in current state")
			return
		}
		response.Error(c, http.StatusBadRequest, errors.InvalidRequestBody, errMsg)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Event approved successfully",
		"data": gin.H{
			"event_id": eventID,
			"status":   "approved",
		},
	})
}

func (h *AdminHandler) RejectEventHandler(c *gin.Context) {
	adminID := c.GetString("user_id")
	eventID := c.Param("event_id")

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errors.InvalidRequestBody, "Invalid request body")
		return
	}

	err := h.eventUsecase.RejectEvent(c.Request.Context(), adminID, eventID, req.Reason)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "EVT_004") {
			response.Error(c, http.StatusForbidden, "EVT_004", "Insufficient admin privileges")
			return
		} else if strings.Contains(errMsg, "EVT_006") {
			response.Error(c, http.StatusConflict, "EVT_006", "Event cannot be rejected in current state")
			return
		}
		response.Error(c, http.StatusBadRequest, errors.InvalidRequestBody, errMsg)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Event rejected successfully",
		"data": gin.H{
			"event_id": eventID,
			"status":   "rejected",
		},
	})
}
