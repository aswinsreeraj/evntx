package repository

import "github.com/aswinsreeraj/evntx/internal/domain"

type UserRepository interface {
	Create(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
	FindByID(id string) (*domain.User, error)
	Update(user *domain.User) error
}
