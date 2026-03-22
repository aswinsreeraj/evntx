package http

import (
	"strconv"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/internal/usecase"
	apiErrors "github.com/aswinsreeraj/evntx/pkg/errors"
	"github.com/aswinsreeraj/evntx/pkg/response"
	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	usecase     *usecase.EventUsecase
	userUsecase *usecase.UserUsecase
}

func NewEventHandler(u *usecase.EventUsecase, uu *usecase.UserUsecase) *EventHandler {
	return &EventHandler{usecase: u, userUsecase: uu}
}

func (h *EventHandler) ListEvents(c *gin.Context) {

	city := c.Query("city")
	category := c.Query("category")
	search := c.Query("search")
	sort := c.Query("sort")
	minPrice := c.Query("min_price")
	maxPrice := c.Query("max_price")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	events, total, min_price_val, max_price_val, err := h.usecase.ListEvents(city, category, search, sort, minPrice, maxPrice, startDate, endDate, page, limit)

	if err != nil {
		response.AppError(c, apiErrors.ErrInternalServerError)
		return
	}

	response.Success(c, "Events fetched successfully", gin.H{
		"events": events,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
		"price_range": gin.H{
			"min": min_price_val,
			"max": max_price_val,
		},
	})
}

func (h *EventHandler) GetEvent(c *gin.Context) {

	slug := c.Param("slug")

	event, details, personnels, tickets, err := h.usecase.GetEvent(slug)

	if err != nil {
		response.AppError(c, apiErrors.ErrResourceNotFound)
		return
	}

	host := gin.H(nil)
	if evt, ok := event.(*domain.Event); ok && h.userUsecase != nil {
		user, organizerDetail, _, userErr := h.userUsecase.GetProfile(evt.OrganizerID)
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
	}

	response.Success(c, "Event fetched successfully", gin.H{
		"event":        event,
		"details":      details,
		"personnels":   personnels,
		"ticket_types": tickets,
		"host":         host,
	})
}
