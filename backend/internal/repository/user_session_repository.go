package repository

import "github.com/aswinsreeraj/evntx/internal/domain"

type UserSessionRepository interface {
	Create(session *domain.UserSession) error
	FindByUserID(userID string) (*domain.UserSession, error)
	Revoke(sessionID string) error
}
