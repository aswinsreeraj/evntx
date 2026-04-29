package repository

import (
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"gorm.io/gorm"
)

type JobLogModel struct {
	ID		string		`gorm:"type:uuid;primaryKey"`
	JobName		string		`gorm:"index;not null"`
	Status		string		`gorm:"index;not null"`
	Attempts	int		`gorm:"not null"`
	ErrorMessage	string		`gorm:"type:text"`
	StartedAt	time.Time	`gorm:"not null"`
	EndedAt		time.Time	`gorm:"not null"`
}

type jobGormRepository struct {
	db *gorm.DB
}

func NewJobGormRepository(db *gorm.DB) *jobGormRepository {
	return &jobGormRepository{db: db}
}

func (r *jobGormRepository) LogJob(log *domain.JobLog) error {
	model := JobLogModel{
		ID:		log.ID,
		JobName:	log.JobName,
		Status:		string(log.Status),
		Attempts:	log.Attempts,
		ErrorMessage:	log.ErrorMessage,
		StartedAt:	log.StartedAt,
		EndedAt:	log.EndedAt,
	}
	return r.db.Create(&model).Error
}
