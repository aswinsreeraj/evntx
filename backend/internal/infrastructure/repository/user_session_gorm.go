package repository

import (
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"gorm.io/gorm"
)

type UserSessionModel struct {
	ID               string `gorm:"type:uuid;primaryKey"`
	UserID           string `gorm:"index"`
	RefreshTokenHash string
	UserAgent        string
	IPAddress        string
	ExpiresAt        time.Time
	Revoked          bool
	CreatedAt        time.Time
}

type userSessionGormRepository struct {
	db *gorm.DB
}

func NewUserSessionGormRepository(db *gorm.DB) *userSessionGormRepository {
	return &userSessionGormRepository{db: db}
}

func (r *userSessionGormRepository) Create(session *domain.UserSession) error {
	model := UserSessionModel{
		ID:               session.ID,
		UserID:           session.UserID,
		RefreshTokenHash: session.RefreshTokenHash,
		UserAgent:        session.UserAgent,
		IPAddress:        session.IPAddress,
		ExpiresAt:        session.ExpiresAt,
		Revoked:          session.Revoked,
	}

	return r.db.Create(&model).Error
}
