package domain

import "time"

const (
	WalletTransactionTypeCredit = "cr"
	WalletTransactionTypeDebit  = "dr"

	WalletTransactionStatusPending   = "pending"
	WalletTransactionStatusCompleted = "completed"
	WalletTransactionStatusFailed    = "failed"
)

type Wallet struct {
	ID               string    `json:"id"`
	UserID           string    `json:"user_id"`
	AvailableBalance float64   `json:"available_balance"`
	PendingBalance   float64   `json:"pending_balance"`
	TotalCredited    float64   `json:"total_credited"`
	TotalDebited     float64   `json:"total_debited"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type WalletTransaction struct {
	ID            string    `json:"id"`
	WalletID      string    `json:"wallet_id"`
	Type          string    `json:"type"`
	Amount        float64   `json:"amount"`
	ReferenceType string    `json:"reference_type"`
	ReferenceID   string    `json:"reference_id"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

type WalletTransactionFilter struct {
	Type   string
	Status string
}
