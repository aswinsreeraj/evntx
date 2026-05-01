package repository

import (
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"gorm.io/gorm"
)

type CategoryModel struct {
	ID        string `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string `gorm:"type:varchar(255);uniqueIndex;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (CategoryModel) TableName() string {
	return "categories"
}

type categoryGormRepository struct {
	db *gorm.DB
}

func NewCategoryGormRepository(db *gorm.DB) *categoryGormRepository {
	return &categoryGormRepository{db: db}
}

func (r *categoryGormRepository) Create(category *domain.Category) error {
	model := &CategoryModel{
		Name: category.Name,
	}

	if err := r.db.Create(model).Error; err != nil {
		return err
	}

	category.ID = model.ID
	category.CreatedAt = model.CreatedAt
	category.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *categoryGormRepository) GetAll() ([]domain.Category, error) {
	var models []CategoryModel
	if err := r.db.Order("name asc").Find(&models).Error; err != nil {
		return nil, err
	}

	categories := make([]domain.Category, len(models))
	for i, m := range models {
		categories[i] = domain.Category{
			ID:        m.ID,
			Name:      m.Name,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		}
	}
	return categories, nil
}

func (r *categoryGormRepository) GetByID(id string) (*domain.Category, error) {
	var m CategoryModel
	if err := r.db.First(&m, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // or return a specific not found error depending on standard
		}
		return nil, err
	}

	return &domain.Category{
		ID:        m.ID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}, nil
}

func (r *categoryGormRepository) Update(category *domain.Category) error {
	var m CategoryModel
	if err := r.db.First(&m, "id = ?", category.ID).Error; err != nil {
		return err
	}

	m.Name = category.Name
	if err := r.db.Save(&m).Error; err != nil {
		return err
	}

	category.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *categoryGormRepository) Delete(id string) error {
	return r.db.Delete(&CategoryModel{}, "id = ?", id).Error
}
