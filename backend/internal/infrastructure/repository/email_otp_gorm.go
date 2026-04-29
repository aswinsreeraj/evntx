package repository

import (
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"gorm.io/gorm"
)

type EmailOTPModel struct {
	ID		string	`gorm:"type:uuid;primaryKey"`
	Email		string	`gorm:"index"`
	OTPHash		string
	ExpiresAt	time.Time
	Consumed	bool
	CreatedAt	time.Time
}

type emailOTPGormRepository struct {
	db *gorm.DB
}

func NewEmailOTPGormRepository(db *gorm.DB) *emailOTPGormRepository {
	return &emailOTPGormRepository{db: db}
}

func (r *emailOTPGormRepository) Create(otp *domain.EmailOTP) error {
	model := EmailOTPModel{
		ID:		otp.ID,
		Email:		otp.Email,
		OTPHash:	otp.OTPHash,
		ExpiresAt:	otp.ExpiresAt,
		Consumed:	otp.Consumed,
	}

	return r.db.Create(&model).Error
}

func (r *emailOTPGormRepository) InvalidatePrevious(email string) error {
	return r.db.Model(&EmailOTPModel{}).
		Where("email = ? AND consumed = ?", email, false).
		Update("consumed", true).Error
}

func (r *emailOTPGormRepository) FindValidOTP(email string) (*domain.EmailOTP, error) {
	var model EmailOTPModel

	err := r.db.
		Where("email = ? AND consumed = ? AND expires_at > ?", email, false, time.Now()).
		Order("created_at DESC").
		First(&model).Error

	if err != nil {
		return nil, err
	}

	return &domain.EmailOTP{
		ID:		model.ID,
		Email:		model.Email,
		OTPHash:	model.OTPHash,
		ExpiresAt:	model.ExpiresAt,
		Consumed:	model.Consumed,
		CreatedAt:	model.CreatedAt,
	}, nil
}

func (r *emailOTPGormRepository) MarkConsumed(id string) error {
	return r.db.Model(&EmailOTPModel{}).
		Where("id = ?", id).
		Update("consumed", true).Error
}
