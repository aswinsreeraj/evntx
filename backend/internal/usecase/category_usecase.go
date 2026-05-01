package usecase

import (
	"errors"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/internal/repository"
)

type CategoryUsecase struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryUsecase(categoryRepo repository.CategoryRepository) *CategoryUsecase {
	return &CategoryUsecase{
		categoryRepo: categoryRepo,
	}
}

func (u *CategoryUsecase) CreateCategory(name string) (*domain.Category, error) {
	if name == "" {
		return nil, errors.New("category name is required")
	}

	category := &domain.Category{
		Name: name,
	}

	err := u.categoryRepo.Create(category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (u *CategoryUsecase) GetAllCategories() ([]domain.Category, error) {
	return u.categoryRepo.GetAll()
}

func (u *CategoryUsecase) UpdateCategory(id string, name string) (*domain.Category, error) {
	if id == "" {
		return nil, errors.New("category ID is required")
	}
	if name == "" {
		return nil, errors.New("category name is required")
	}

	cat, err := u.categoryRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if cat == nil {
		return nil, errors.New("category not found")
	}

	cat.Name = name
	err = u.categoryRepo.Update(cat)
	if err != nil {
		return nil, err
	}

	return cat, nil
}

func (u *CategoryUsecase) DeleteCategory(id string) error {
	if id == "" {
		return errors.New("category ID is required")
	}
	return u.categoryRepo.Delete(id)
}
