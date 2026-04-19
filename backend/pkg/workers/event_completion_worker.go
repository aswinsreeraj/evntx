package workers

import (
	"context"

	"github.com/aswinsreeraj/evntx/internal/usecase"
	"github.com/aswinsreeraj/evntx/pkg/logger"
)

func AutoProcessCompletedEventsJob(eventUsecase *usecase.EventUsecase) JobFunc {
	return func(ctx context.Context) error {
		logger.Log.Info().Msg("Executing AutoProcessCompletedEventsJob...")
		
		err := eventUsecase.AutoProcessCompletedEvents(ctx)
		if err != nil {
			return err
		}
		
		return nil
	}
}
