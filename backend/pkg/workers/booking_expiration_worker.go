package workers

import (
	"context"

	"github.com/aswinsreeraj/evntx/internal/usecase"
	"github.com/aswinsreeraj/evntx/pkg/logger"
)

func ProcessExpiredBookingsJob(bookingUsecase *usecase.BookingUsecase) JobFunc {
	return func(ctx context.Context) error {
		logger.Log.Info().Msg("Executing BookingExpirationJob...")
		
		err := bookingUsecase.ProcessExpiredBookings(ctx)
		if err != nil {
			return err
		}
		
		return nil
	}
}
