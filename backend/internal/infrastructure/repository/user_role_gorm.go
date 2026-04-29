package repository

import (
	"github.com/aswinsreeraj/evntx/internal/domain"
	"gorm.io/gorm"
)

type UserRoleModel struct {
	UserID	string	`gorm:"primaryKey"`
	Role	string	`gorm:"primaryKey"`
}

type userRoleGormRepository struct {
	db *gorm.DB
}

func NewUserRoleGormRepository(db *gorm.DB) *userRoleGormRepository {
	return &userRoleGormRepository{db: db}
}

func (r *userRoleGormRepository) GetRolesByUserID(userID string) ([]domain.UserRole, error) {
	var models []UserRoleModel

	err := r.db.Where("user_id = ?", userID).Find(&models).Error
	if err != nil {
		return nil, err
	}

	roles := make([]domain.UserRole, 0)
	for _, m := range models {
		roles = append(roles, domain.UserRole(m.Role))
	}

	return roles, nil
}

func (r *userRoleGormRepository) AddRole(userID string, role domain.UserRole) error {
	model := UserRoleModel{
		UserID:	userID,
		Role:	string(role),
	}
	return r.db.Save(&model).Error
}

func (r *userRoleGormRepository) RemoveRole(userID string, role domain.UserRole) error {
	return r.db.Delete(&UserRoleModel{}, "user_id = ? AND role = ?", userID, string(role)).Error
}
