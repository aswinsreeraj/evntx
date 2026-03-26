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

type RazorpayOrder struct {
	ID          string          `json:"id"`
	Entity      string          `json:"entity"`
	Amount      int64           `json:"amount"`
	Currency    string          `json:"currency"`
	Receipt     string          `json:"receipt"`
	Status      string          `json:"status"`
	CreatedAt   int64           `json:"created_at"`
	RawResponse json.RawMessage `json:"-"`
}

type PaymentOrderResponse struct {
	OrderID     string `json:"order_id"`
	Amount      int64  `json:"amount"`
	Currency    string `json:"currency"`
	RazorpayKey string `json:"razorpay_key"`
}
