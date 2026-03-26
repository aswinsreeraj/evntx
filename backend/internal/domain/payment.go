package domain

import (
	"encoding/json"
	"time"
)

const (
	PaymentStatusInitiated = "initiated"
	PaymentStatusSuccess   = "success"
	PaymentStatusFailed    = "failed"
	PaymentStatusRefunded  = "refunded"
)

type Payment struct {
	ID                string          `json:"id"`
	BookingID         string          `json:"booking_id"`
	Provider          string          `json:"provider"`
	ProviderReference string          `json:"provider_reference"`
	Amount            float64         `json:"amount"`
	Status            string          `json:"status"`
	RawResponse       json.RawMessage `json:"raw_response"`
	CreatedAt         time.Time       `json:"created_at"`
}
