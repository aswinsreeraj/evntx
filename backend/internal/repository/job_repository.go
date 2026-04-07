package repository

import "github.com/aswinsreeraj/evntx/internal/domain"

type JobRepository interface {
	LogJob(log *domain.JobLog) error
}
