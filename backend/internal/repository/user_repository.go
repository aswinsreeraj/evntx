package repository

import "github.com/aswinsreeraj/evntx/internal/domain"

type UserRepository interface {
	Create(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
	FindByID(id string) (*domain.User, error)
	Update(user *domain.User) error
	Search(
		search string,
		status string,
		page int,
		limit int,
	) ([]domain.User, int64, error)
	SearchOrganizers(
		search string,
		status string,
		page int,
		limit int,
	) ([]domain.OrganizerDetails, int64, error)

	UpdateStatus(userID string, isActive bool) error
}
