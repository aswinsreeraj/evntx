package http

import (
	"net/http"
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
