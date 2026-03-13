package domain

import "time"

// TicketType represent different pricing tiers or ticket categories for an Event.
type TicketType struct {
	ID                string
	EventID           string
	Name              string
	Price             float64
	TotalQuantity     int
	AvailableQuantity int
	Version           int
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
