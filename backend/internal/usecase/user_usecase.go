package usecase

import (
	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/internal/repository"
	"github.com/google/uuid"
)

type UserUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(r repository.UserRepository) *UserUsecase {
	return &UserUsecase{repo: r}
}

func (u *UserUsecase) Register(email string) (*domain.User, error) {
	user := &domain.User{
		ID:            uuid.NewString(),
		Email:         email,
		IsActive:      true,
		EmailVerified: false,
	}

	err := u.repo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
