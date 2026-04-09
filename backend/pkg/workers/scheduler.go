package workers

import (
	"context"
	"time"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/internal/repository"
	"github.com/aswinsreeraj/evntx/pkg/logger"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
)

type JobFunc func(ctx context.Context) error

type CronScheduler struct {
	cron *cron.Cron
	repo repository.JobRepository
}

func NewCronScheduler(repo repository.JobRepository) *CronScheduler {
	return &CronScheduler{
		cron: cron.New(),
		repo: repo,
	}
}

func (s *CronScheduler) Start() {
	s.cron.Start()
	logger.Log.Info().Msg("Cron Scheduler started")
}

func (s *CronScheduler) Stop() {
	s.cron.Stop()
	logger.Log.Info().Msg("Cron Scheduler stopped")
}

func (s *CronScheduler) RegisterJob(name string, schedule string, job JobFunc, maxAttempts int) error {
	_, err := s.cron.AddFunc(schedule, func() {
		s.executeJobWithRetries(name, job, maxAttempts)
	})
	if err != nil {
		logger.Log.Error().Err(err).Str("job", name).Msg("Failed to register job")
		return err
	}
	logger.Log.Info().Str("job", name).Str("schedule", schedule).Msg("Job registered successfully")
	return nil
}

func (s *CronScheduler) executeJobWithRetries(name string, job JobFunc, maxAttempts int) {
	startedAt := time.Now()
	var finalErr error
	var attempts int

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		attempts = attempt
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		
		err := job(ctx)
		cancel()

		if err == nil {
			finalErr = nil
			break
		}

		finalErr = err
		logger.Log.Warn().Err(err).Str("job", name).Int("attempt", attempt).Msg("Job failed")
		
		if attempt < maxAttempts {
			
			time.Sleep(60 * time.Second)
		}
	}

	status := domain.JobStatusSuccess
	errMsg := ""
	if finalErr != nil {
		status = domain.JobStatusFailed
		errMsg = finalErr.Error()
		logger.Log.Error().Err(finalErr).Str("job", name).Msg("Job exhausted attempts and failed")
	} else {
		logger.Log.Info().Str("job", name).Int("attempts", attempts).Msg("Job completed successfully")
	}

	jobLog := &domain.JobLog{
		ID:           uuid.NewString(),
		JobName:      name,
		Status:       status,
		Attempts:     attempts,
		ErrorMessage: errMsg,
		StartedAt:    startedAt,
		EndedAt:      time.Now(),
	}

	if s.repo != nil {
		if err := s.repo.LogJob(jobLog); err != nil {
			logger.Log.Error().Err(err).Str("job", name).Msg("Failed to log job to database")
		}
	}
}
