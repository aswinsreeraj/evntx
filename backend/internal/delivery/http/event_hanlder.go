package http

import (
	"net/http"
	"strconv"

	"github.com/aswinsreeraj/evntx/internal/usecase"
	"github.com/aswinsreeraj/evntx/pkg/response"
	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	usecase *usecase.EventUsecase
}

func NewEventHandler(u *usecase.EventUsecase) *EventHandler {
	return &EventHandler{usecase: u}
}

func (h *EventHandler) ListEvents(c *gin.Context) {

	city := c.Query("city")

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	events, total, err := h.usecase.ListEvents(city, page, limit)

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "EVT_002", "Failed to fetch events")
		return
	}

	response.Success(c, "Events fetched successfully", gin.H{
		"events": events,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (h *EventHandler) GetEvent(c *gin.Context) {

	slug := c.Param("slug")

	event, details, personnels, err := h.usecase.GetEvent(slug)

	if err != nil {
		response.Error(c, http.StatusNotFound, "EVT_002", "Event not found")
		return
	}

	response.Success(c, "Event fetched successfully", gin.H{
		"event":      event,
		"details":    details,
		"personnels": personnels,
	})
}
