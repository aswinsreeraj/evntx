package workers

import (
	"context"

	"github.com/aswinsreeraj/evntx/internal/usecase"
	"github.com/aswinsreeraj/evntx/pkg/logger"
)

func ProcessPayoutSettlementsJob(walletUsecase *usecase.WalletUsecase) JobFunc {
	return func(ctx context.Context) error {
		logger.Log.Info().Msg("Executing ProcessPayoutSettlementsJob...")

		err := walletUsecase.AutoProcessApprovedPayouts(ctx)
		if err != nil {
			return err
		}

		return nil
	}
}
