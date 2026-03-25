package domain

import "time"

type TicketType struct {
	ID                string    `json:"id"`
	EventID           string    `json:"event_id"`
	Name              string    `json:"name"`
	Price             float64   `json:"price"`
	TotalQuantity     int       `json:"total_quantity"`
	AvailableQuantity int       `json:"available_quantity"`
	Version           int       `json:"version"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
