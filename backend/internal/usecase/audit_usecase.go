package usecase

import (
	"encoding/json"
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/internal/repository"
	"github.com/google/uuid"
)

type AuditUsecase struct {
	auditRepo repository.AuditRepository
	userRepo  repository.UserRepository
}

func NewAuditUsecase(auditRepo repository.AuditRepository, userRepo repository.UserRepository) *AuditUsecase {
	return &AuditUsecase{
		auditRepo: auditRepo,
		userRepo:  userRepo,
	}
}

func (u *AuditUsecase) LogAction(adminID, action string, tag domain.AuditActionTag, detailsMap map[string]interface{}, ipAddress string) error {
	adminName := "System/Unknown"
	
	if adminID != "" {
		admin, err := u.userRepo.FindByID(adminID)
		if err == nil && admin != nil {
			adminName = admin.Name
		}
	}

	detailsJSON := "{}"
	if detailsMap != nil {
		if b, err := json.Marshal(detailsMap); err == nil {
			detailsJSON = string(b)
		}
	}

	log := &domain.AuditLog{
		ID:        uuid.NewString(),
		AdminID:   adminID,
		AdminName: adminName,
		Action:    action,
		ActionTag: tag,
		Details:   detailsJSON,
		IPAddress: ipAddress,
		Timestamp: time.Now(),
	}

	return u.auditRepo.Create(log)
}

func (u *AuditUsecase) GetLogs(page, limit int) ([]domain.AuditLog, int64, error) {
	return u.auditRepo.GetLogs(page, limit)
}
