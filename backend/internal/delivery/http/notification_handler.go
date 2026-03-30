package http

import (
	"strconv"

	"github.com/aswinsreeraj/evntx/internal/usecase"
	apiErrors "github.com/aswinsreeraj/evntx/pkg/errors"
	"github.com/aswinsreeraj/evntx/pkg/response"
	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	notificationUsecase *usecase.NotificationUsecase
}

func NewNotificationHandler(notificationUsecase *usecase.NotificationUsecase) *NotificationHandler {
	return &NotificationHandler{notificationUsecase: notificationUsecase}
}

func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userID := c.GetString("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	notifications, total, unreadCount, err := h.notificationUsecase.GetNotificationsByUser(userID, page, limit)
	if err != nil {
		response.AppError(c, err)
		return
	}

	response.Success(c, "Notifications fetched successfully", gin.H{
		"notifications": notifications,
		"unread_count":  unreadCount,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userID := c.GetString("user_id")
	notificationID := c.Param("id")

	if err := h.notificationUsecase.MarkAsRead(userID, notificationID); err != nil {
		response.AppError(c, err)
		return
	}

	response.Success(c, "Notification marked as read", nil)
}

func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID := c.GetString("user_id")

	if err := h.notificationUsecase.MarkAllAsRead(userID); err != nil {
		response.AppError(c, apiErrors.ErrInternalServerError)
		return
	}

	response.Success(c, "All notifications marked as read", nil)
}

func (h *NotificationHandler) ClearAll(c *gin.Context) {
	userID := c.GetString("user_id")

	if err := h.notificationUsecase.ClearAll(userID); err != nil {
		response.AppError(c, apiErrors.ErrInternalServerError)
		return
	}

	response.Success(c, "All notifications cleared", nil)
}

