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
			"updated_at":     time.Now(),
		}).Error
}

func (r *userGormRepository) Search(
	search string,
	page int,
	limit int,
) ([]domain.User, int64, error) {

	var models []UserModel
	var total int64

	query := r.db.Model(&UserModel{})

	if search != "" {
		query = query.Where(
			"name ILIKE ? OR email ILIKE ?",
			"%"+search+"%",
			"%"+search+"%",
		)
	}

	query.Count(&total)

	offset := (page - 1) * limit

	err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&models).Error

	if err != nil {
		return nil, 0, err
	}

	users := make([]domain.User, 0)
	for _, m := range models {
		users = append(users, domain.User{
			ID:            m.ID,
			Name:          m.Name,
			Email:         m.Email,
			IsActive:      m.IsActive,
			EmailVerified: m.EmailVerified,
			CreatedAt:     m.CreatedAt,
			UpdatedAt:     m.UpdatedAt,
		})
	}

	return users, total, nil
}

func (r *userGormRepository) UpdateStatus(userID string, isActive bool) error {
	return r.db.Model(&UserModel{}).
		Where("id = ?", userID).
		Update("is_active", isActive).Error
}
