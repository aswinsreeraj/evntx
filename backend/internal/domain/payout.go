package domain

import (
	"time"
)

type PayoutStatus string

const (
	PayoutStatusPending    PayoutStatus = "pending"
	PayoutStatusApproved   PayoutStatus = "approved"
	PayoutStatusRejected   PayoutStatus = "rejected"
	PayoutStatusProcessing PayoutStatus = "processing"
	PayoutStatusCompleted  PayoutStatus = "completed"
	PayoutStatusFailed     PayoutStatus = "failed"
)

type PayoutRequest struct {
	ID            string       `json:"id"`
	UserID        string       `json:"user_id"`
	Amount        float64      `json:"amount"`
	Status        PayoutStatus `json:"status"`
	RequestedAt   time.Time    `json:"requested_at"`
	ReviewedAt    *time.Time   `json:"reviewed_at,omitempty"`
	ProcessedAt   *time.Time   `json:"processed_at,omitempty"`
	AdminID       *string      `json:"admin_id,omitempty"`
	FailureReason *string      `json:"failure_reason,omitempty"`
	CreatedAt     time.Time    `json:"created_at"`
}

type PayoutCredential struct {
	ID                     string    `json:"id"`
	UserID                 string    `json:"user_id"`
	AccountHolderName      string    `json:"account_holder_name"`
	AccountNumberEncrypted string    `json:"-"`
	IFSCCodeEncrypted      string    `json:"-"`
	UPIIDEncrypted         *string   `json:"-"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

type AdminPayoutDetail struct {
	PayoutRequest
	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email"`
}
