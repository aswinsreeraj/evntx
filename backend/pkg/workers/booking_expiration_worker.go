package workers

import (
	"context"
	"time"

	"github.com/aswinsreeraj/evntx/internal/usecase"
	"github.com/aswinsreeraj/evntx/pkg/logger"
)

type BookingExpirationWorker struct {
	bookingUsecase *usecase.BookingUsecase
	ticker         *time.Ticker
	quit           chan struct{}
}

func NewBookingExpirationWorker(bookingUsecase *usecase.BookingUsecase) *BookingExpirationWorker {
	return &BookingExpirationWorker{
		bookingUsecase: bookingUsecase,
		quit:           make(chan struct{}),
	}
}

func (w *BookingExpirationWorker) Start() {
	// Tick every 1 minute
	w.ticker = time.NewTicker(1 * time.Minute)

	logger.Log.Info().Msg("BookingExpirationWorker started")

	for {
		select {
		case <-w.ticker.C:
			// Execute processing
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			err := w.bookingUsecase.ProcessExpiredBookings(ctx)
			cancel()

			if err != nil {
				logger.Log.Error().Err(err).Msg("failed to process expired bookings")
			}
		case <-w.quit:
			w.ticker.Stop()
			logger.Log.Info().Msg("BookingExpirationWorker stopped")
			return
		}
	}
}

func (w *BookingExpirationWorker) Stop() {
	close(w.quit)
}
