package repository

import "github.com/aswinsreeraj/evntx/internal/domain"

type UserRoleRepository interface {
	GetRolesByUserID(userID string) ([]domain.UserRole, error)
	AddRole(userID string, role domain.UserRole) error
	RemoveRole(userID string, role domain.UserRole) error
}
