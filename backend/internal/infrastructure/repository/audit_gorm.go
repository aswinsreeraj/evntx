package repository

import (
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"gorm.io/gorm"
)

type AuditLogModel struct {
	ID		string		`gorm:"type:uuid;primaryKey"`
	AdminID		string		`gorm:"type:uuid;index;not null"`
	AdminName	string		`gorm:"not null"`
	Action		string		`gorm:"not null"`
	ActionTag	string		`gorm:"index;not null"`
	Details		string		`gorm:"type:jsonb"`
	IPAddress	string		`gorm:"size:45"`
	Timestamp	time.Time	`gorm:"index;not null"`
}

type auditGormRepository struct {
	db *gorm.DB
}

func NewAuditGormRepository(db *gorm.DB) *auditGormRepository {
	return &auditGormRepository{db: db}
}

func (r *auditGormRepository) Create(log *domain.AuditLog) error {
	model := AuditLogModel{
		ID:		log.ID,
		AdminID:	log.AdminID,
		AdminName:	log.AdminName,
		Action:		log.Action,
		ActionTag:	string(log.ActionTag),
		Details:	log.Details,
		IPAddress:	log.IPAddress,
		Timestamp:	log.Timestamp,
	}

	return r.db.Create(&model).Error
}

func (r *auditGormRepository) GetLogs(page, limit int) ([]domain.AuditLog, int64, error) {
	var models []AuditLogModel
	var total int64

	query := r.db.Model(&AuditLogModel{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Order("timestamp DESC").Offset(offset).Limit(limit).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	var parsedLogs []domain.AuditLog
	for _, m := range models {
		parsedLogs = append(parsedLogs, domain.AuditLog{
			ID:		m.ID,
			AdminID:	m.AdminID,
			AdminName:	m.AdminName,
			Action:		m.Action,
			ActionTag:	domain.AuditActionTag(m.ActionTag),
			Details:	m.Details,
			IPAddress:	m.IPAddress,
			Timestamp:	m.Timestamp,
		})
	}

	return parsedLogs, total, nil
}
