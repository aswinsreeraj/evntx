package domain

import "time"

const (
	PlatformWalletID = "00000000-0000-0000-0000-000000000001"

	PlatformRefTypePayment  = "payment"
	PlatformRefTypeRefund   = "refund"
	PlatformRefTypePayout   = "payout"
	PlatformRefTypeEarning  = "earning"
)

type PlatformWallet struct {
	ID               string    `json:"id"`
	AvailableBalance float64   `json:"available_balance"`
	PendingBalance   float64   `json:"pending_balance"`
	TotalCredited    float64   `json:"total_credited"`
	TotalDebited     float64   `json:"total_debited"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type PlatformWalletTransaction struct {
	ID            string    `json:"id"`
	WalletID      string    `json:"wallet_id"`
	Type          string    `json:"type"`
	Amount        float64   `json:"amount"`
	ReferenceType string    `json:"reference_type"`
	ReferenceID   string    `json:"reference_id"`
	CreatedAt     time.Time `json:"created_at"`
}
