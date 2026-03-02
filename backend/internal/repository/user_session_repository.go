package repository

import "github.com/aswinsreeraj/evntx/internal/domain"

type UserSessionRepository interface {
	Create(session *domain.UserSession) error
}
