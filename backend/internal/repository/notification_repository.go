package repository

import "github.com/aswinsreeraj/evntx/internal/domain"

type NotificationRepository interface {
	CreateNotification(notification *domain.Notification) error
	GetNotificationsByUser(userID string, page int, limit int) ([]domain.Notification, int64, int64, error)
	MarkAsRead(notificationID string, userID string) error
	MarkAllAsRead(userID string) error
	ClearAll(userID string) error
}
