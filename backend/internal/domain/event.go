package domain

import "time"

type Event struct {
	ID            string    `json:"id"`
	OrganizerID   string    `json:"organizer_id"`
	Title         string    `json:"title"`
	Slug          string    `json:"slug"`
	City          string    `json:"city"`
	VenueName     string    `json:"venue_name"`
	Category      string    `json:"category"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	Tags          string    `json:"tags"`
	Status        string    `json:"status"`
	CoverImageURL string    `json:"cover_image_url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type AdminEventDetails struct {
	Event
	OrganizerName string `json:"organizer_name"`
	TicketsSold   int64  `json:"tickets_sold"`
	Revenue       int64  `json:"revenue"`
}

type TicketCancelRequest struct {
	TicketType string `json:"ticket_type" binding:"required"`
	Quantity   int    `json:"quantity" binding:"required,min=1"`
}
