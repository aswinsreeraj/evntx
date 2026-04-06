package domain

import (
	"time"
)

type RefundStatus string

const (
	RefundStatusPending   RefundStatus = "pending"
	RefundStatusProcessed RefundStatus = "processed"
)

type RefundRequest struct {
	ID          string       `json:"id"`
	UserID      string       `json:"user_id"`
	BookingID   string       `json:"booking_id"`
	Amount      float64      `json:"amount"`
	Status      RefundStatus `json:"status"`
	RequestedAt time.Time    `json:"requested_at"`
	ProcessedAt *time.Time   `json:"processed_at,omitempty"`
	AdminID     *string      `json:"admin_id,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
}

type AdminRefundDetail struct {
	RefundRequest
	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email"`
}
