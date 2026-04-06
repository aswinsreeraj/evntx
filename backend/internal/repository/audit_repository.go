package repository

import "github.com/aswinsreeraj/evntx/internal/domain"

type AuditRepository interface {
	Create(log *domain.AuditLog) error
	GetLogs(page, limit int) ([]domain.AuditLog, int64, error)
}
