package usecase

import (
	"encoding/json"
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/internal/repository"
	apiErrors "github.com/aswinsreeraj/evntx/pkg/errors"
	"github.com/google/uuid"
)

type NotificationUsecase struct {
	repo repository.NotificationRepository
}

func NewNotificationUsecase(repo repository.NotificationRepository) *NotificationUsecase {
	return &NotificationUsecase{repo: repo}
}

func (u *NotificationUsecase) SendNotification(
	userID string,
	notificationType string,
	title string,
	message string,
	metadata interface{},
) error {
	if userID == "" {
		return apiErrors.New(400, apiErrors.InvalidRequestBody, "User ID is required")
	}

	var rawMetadata json.RawMessage
	if metadata != nil {
		payload, err := json.Marshal(metadata)
		if err != nil {
			return err
		}
		rawMetadata = payload
	}

	notification := &domain.Notification{
		ID:		uuid.NewString(),
		UserID:		userID,
		Type:		notificationType,
		Title:		title,
		Message:	message,
		IsRead:		false,
		Metadata:	rawMetadata,
		CreatedAt:	time.Now(),
	}

	return u.repo.CreateNotification(notification)
}

func (u *NotificationUsecase) GetNotificationsByUser(
	userID string,
	page int,
	limit int,
) ([]domain.Notification, int64, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	return u.repo.GetNotificationsByUser(userID, page, limit)
}

func (u *NotificationUsecase) MarkAsRead(userID string, notificationID string) error {
	if notificationID == "" {
		return apiErrors.New(400, apiErrors.InvalidRequestBody, "Notification ID is required")
	}

	return u.repo.MarkAsRead(notificationID, userID)
}

func (u *NotificationUsecase) MarkAllAsRead(userID string) error {
	return u.repo.MarkAllAsRead(userID)
}

func (u *NotificationUsecase) ClearAll(userID string) error {
	return u.repo.ClearAll(userID)
}
