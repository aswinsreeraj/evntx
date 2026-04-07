package workers

import (
	"context"

	"github.com/aswinsreeraj/evntx/internal/usecase"
	"github.com/aswinsreeraj/evntx/pkg/logger"
)

// ProcessPayoutSettlementsJob processes all "Approved" (but unpaid) payout requests automatically 
// or automates transfers. Here we will define a generic processor on the wallet usecase.
func ProcessPayoutSettlementsJob(walletUsecase *usecase.WalletUsecase) JobFunc {
	return func(ctx context.Context) error {
		logger.Log.Info().Msg("Executing ProcessPayoutSettlementsJob...")

		// Note: The wallet usecase should have an exposed method to process approved payouts
		// that haven't been transferred yet.
		err := walletUsecase.AutoProcessApprovedPayouts(ctx)
		if err != nil {
			return err
		}
		
		return nil
	}
}
