package repository

import "github.com/aswinsreeraj/evntx/internal/domain"

type CategoryRepository interface {
	Create(category *domain.Category) error
	GetAll() ([]domain.Category, error)
	GetByID(id string) (*domain.Category, error)
	Update(category *domain.Category) error
	Delete(id string) error
}
