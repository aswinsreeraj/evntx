package domain

import (
	"encoding/json"
	"time"
)

const (
	NotificationTypeBookingReserved = "booking_reserved"
	NotificationTypeBookingCancelled = "booking_cancelled"
	NotificationTypePaymentSuccess  = "payment_success"
	NotificationTypePaymentFailed   = "payment_failed"
	NotificationTypeTicketGenerated = "ticket_generated"
	NotificationTypeCheckInSuccess  = "check_in_success"
	NotificationTypeSettlement      = "settlement"
	NotificationTypePayoutRequest   = "payout_request"
)

type Notification struct {
	ID        string          `json:"id"`
	UserID    string          `json:"user_id"`
	Type      string          `json:"type"`
	Title     string          `json:"title"`
	Message   string          `json:"message"`
	IsRead    bool            `json:"is_read"`
	Metadata  json.RawMessage `json:"metadata,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
}
