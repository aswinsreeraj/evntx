package repository

import (
	"encoding/json"
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
	apiErrors "github.com/aswinsreeraj/evntx/pkg/errors"
	"gorm.io/gorm"
)

type NotificationModel struct {
	ID        string          `gorm:"type:uuid;primaryKey"`
	UserID    string          `gorm:"type:uuid;index;not null"`
	Type      string          `gorm:"not null"`
	Title     string          `gorm:"not null"`
	Message   string          `gorm:"type:text;not null"`
	IsRead    bool            `gorm:"default:false;not null"`
	Metadata  json.RawMessage `gorm:"type:jsonb"`
	CreatedAt time.Time       `gorm:"not null"`
}

type notificationGormRepository struct {
	db *gorm.DB
}

func NewNotificationGormRepository(db *gorm.DB) *notificationGormRepository {
	return &notificationGormRepository{db: db}
}

func (r *notificationGormRepository) CreateNotification(notification *domain.Notification) error {
	model := NotificationModel{
		ID:        notification.ID,
		UserID:    notification.UserID,
		Type:      notification.Type,
		Title:     notification.Title,
		Message:   notification.Message,
		IsRead:    notification.IsRead,
		Metadata:  notification.Metadata,
		CreatedAt: notification.CreatedAt,
	}

	return r.db.Create(&model).Error
}

func (r *notificationGormRepository) GetNotificationsByUser(
	userID string,
	page int,
	limit int,
) ([]domain.Notification, int64, int64, error) {
	var models []NotificationModel
	var total int64
	var unreadCount int64

	query := r.db.Model(&NotificationModel{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}

	if err := r.db.Model(&NotificationModel{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&unreadCount).Error; err != nil {
		return nil, 0, 0, err
	}

	offset := (page - 1) * limit
	if err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&models).Error; err != nil {
		return nil, 0, 0, err
	}

	notifications := make([]domain.Notification, 0, len(models))
	for _, model := range models {
		notifications = append(notifications, domain.Notification{
			ID:        model.ID,
			UserID:    model.UserID,
			Type:      model.Type,
			Title:     model.Title,
			Message:   model.Message,
			IsRead:    model.IsRead,
			Metadata:  model.Metadata,
			CreatedAt: model.CreatedAt,
		})
	}

	return notifications, total, unreadCount, nil
}

func (r *notificationGormRepository) MarkAsRead(notificationID string, userID string) error {
	result := r.db.Model(&NotificationModel{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Update("is_read", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return apiErrors.ErrResourceNotFound
	}

	return nil
}

func (r *notificationGormRepository) MarkAllAsRead(userID string) error {
	return r.db.Model(&NotificationModel{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true).Error
}
