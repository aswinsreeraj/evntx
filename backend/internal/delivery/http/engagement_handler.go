package http

import (
	"strings"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/internal/usecase"
	jwtutil "github.com/aswinsreeraj/evntx/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type EngagementHandler struct {
	usecase *usecase.EngagementUsecase
}

func NewEngagementHandler(u *usecase.EngagementUsecase) *EngagementHandler {
	return &EngagementHandler{usecase: u}
}

// @Summary Initialize Session
// @Description Creates a new visitor session
// @Tags Engagement
// @Produce json
// @Success 200 {object} domain.VisitorSession
// @Router /api/v1/engagement/session [post]
func (h *EngagementHandler) InitializeSession(c *gin.Context) {
	var userID *string
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if id, err := jwtutil.ParseAccessToken(tokenString); err == nil {
			userID = &id
		}
	}

	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()

	session, err := h.usecase.InitializeSession(c.Request.Context(), userID, ipAddress, userAgent)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to initialize session"})
		return
	}

	c.JSON(200, session)
}

type trackEventRequest struct {
	SessionID string                     `json:"session_id" binding:"required"`
	EventType domain.EngagementEventType `json:"event_type" binding:"required"`
	EventID   *string                    `json:"event_id,omitempty"`
	Metadata  string                     `json:"metadata,omitempty"`
}

// @Summary Track Event
// @Description Logs an engagement action
// @Tags Engagement
// @Param request body trackEventRequest true "Event Details"
// @Produce json
// @Success 204
// @Router /api/v1/engagement/track [post]
func (h *EngagementHandler) TrackEvent(c *gin.Context) {
	var req trackEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	var userID *string
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if id, err := jwtutil.ParseAccessToken(tokenString); err == nil {
			userID = &id
		}
	}

	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()

	err := h.usecase.TrackEvent(
		c.Request.Context(),
		req.SessionID,
		userID,
		req.EventType,
		req.EventID,
		req.Metadata,
		ipAddress,
		userAgent,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to track event"})
		return
	}

	c.Status(204)
}
