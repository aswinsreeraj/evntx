package repository

import (
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"gorm.io/gorm"
)

type UserModel struct {
	ID            string `gorm:"type:uuid;primaryKey"`
	Name          string
	Email         string `gorm:"uniqueIndex"`
	IsActive      bool
	EmailVerified bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type userGormRepository struct {
	db *gorm.DB
}

func NewUserGormRepository(db *gorm.DB) *userGormRepository {
	return &userGormRepository{db: db}
}

func (r *userGormRepository) Create(user *domain.User) error {
	model := UserModel{
		ID:            user.ID,
		Name:          user.Name,
		Email:         user.Email,
		IsActive:      user.IsActive,
		EmailVerified: user.EmailVerified,
	}

	return r.db.Create(&model).Error
}

func (r *userGormRepository) FindByEmail(email string) (*domain.User, error) {
	var model UserModel

	err := r.db.Where("email = ?", email).First(&model).Error
	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:            model.ID,
		Name:          model.Name,
		Email:         model.Email,
		IsActive:      model.IsActive,
		EmailVerified: model.EmailVerified,
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
	}, nil
}

func (r *userGormRepository) FindByID(id string) (*domain.User, error) {
	var model UserModel

	err := r.db.First(&model, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:            model.ID,
		Name:          model.Name,
		Email:         model.Email,
		IsActive:      model.IsActive,
		EmailVerified: model.EmailVerified,
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
	}, nil
}

func (r *userGormRepository) Update(user *domain.User) error {
	return r.db.Model(&UserModel{}).
		Where("id = ?", user.ID).
		Updates(map[string]interface{}{
			"name":           user.Name,
			"email":          user.Email,
			"is_active":      user.IsActive,
			"email_verified": user.EmailVerified,
		}).Error
}
