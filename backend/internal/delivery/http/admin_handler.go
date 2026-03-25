package http

import (
	"strconv"
	"strings"

	"github.com/aswinsreeraj/evntx/internal/usecase"
	"github.com/aswinsreeraj/evntx/internal/domain"
	pkgErrors "github.com/aswinsreeraj/evntx/pkg/errors"
	"github.com/aswinsreeraj/evntx/pkg/response"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	eventUsecase *usecase.EventUsecase
	userUsecase  *usecase.UserUsecase
}

func NewAdminHandler(eventUsecase *usecase.EventUsecase, userUsecase *usecase.UserUsecase) *AdminHandler {
	return &AdminHandler{
		eventUsecase: eventUsecase,
		userUsecase:  userUsecase,
	}
}

func (h *AdminHandler) AdminListEvents(c *gin.Context) {
	search := c.Query("search")
	status := c.Query("status")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page := 1
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}

	limit := 10
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	events, total, err := h.eventUsecase.AdminSearchEvents(search, status, page, limit)
	if err != nil {
		response.AppError(c, pkgErrors.ErrInternalServerError)
		return
	}

	eventResponses := make([]gin.H, 0, len(events))
	for _, evt := range events {
		eventResponses = append(eventResponses, gin.H{
			"id":             evt.ID,
			"slug":           evt.Slug,
			"title":          evt.Title,
			"organizer_name": evt.OrganizerName,
			"start_time":     evt.StartTime,
			"date":           evt.StartTime,
			"tickets_sold":   evt.TicketsSold,
			"revenue":        evt.Revenue,
			"city":           evt.City,
			"status":         evt.Status,
		})
	}

	response.Success(c, "Events retrieved successfully", gin.H{
		"events": eventResponses,
		"pagination": gin.H{
			"total": total,
			"page":  page,
			"limit": limit,
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
			response.AppError(c, pkgErrors.New(403, pkgErrors.ForbiddenAction, "Insufficient admin privileges"))
			return
		} else if strings.Contains(errMsg, "EVT_006") {
			response.AppError(c, pkgErrors.New(409, pkgErrors.InvalidStateTransition, "Event cannot be approved in current state"))
			return
		}
		response.AppError(c, pkgErrors.ErrInternalServerError)
		return
	}

	response.Success(c, "Event approved successfully", gin.H{
		"event_id": eventID,
		"status":   "approved",
	})
}

func (h *AdminHandler) RejectEventHandler(c *gin.Context) {
	adminID := c.GetString("user_id")
	eventID := c.Param("event_id")

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.AppError(c, pkgErrors.ErrInvalidRequestBody)
		return
	}

	err := h.eventUsecase.RejectEvent(c.Request.Context(), adminID, eventID, req.Reason)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "EVT_004") {
			response.AppError(c, pkgErrors.New(403, pkgErrors.ForbiddenAction, "Insufficient admin privileges"))
			return
		} else if strings.Contains(errMsg, "EVT_006") {
			response.AppError(c, pkgErrors.New(409, pkgErrors.InvalidStateTransition, "Event cannot be rejected in current state"))
			return
		}
		response.AppError(c, pkgErrors.ErrInternalServerError)
		return
	}

	response.Success(c, "Event rejected successfully", gin.H{
		"event_id": eventID,
		"status":   "rejected",
	})
}

func (h *AdminHandler) AdminGetEvent(c *gin.Context) {
	slug := c.Param("slug")

	event, details, personnels, tickets, err := h.eventUsecase.AdminGetEvent(slug)

	if err != nil {
		response.AppError(c, pkgErrors.ErrResourceNotFound)
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
