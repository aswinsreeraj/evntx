package repository

import "github.com/aswinsreeraj/evntx/internal/domain"

type UserRoleRepository interface {
	GetRolesByUserID(userID string) ([]domain.UserRole, error)
}
