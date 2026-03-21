package domain

import "time"

type Event struct {
	ID            string
	OrganizerID   string
	Title         string
	Slug          string
	City          string
	VenueName     string
	Category      string
	StartTime     time.Time
	EndTime       time.Time
	Tags          string
	Status        string
	CoverImageURL string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type AdminEventDetails struct {
	Event
	OrganizerName string `json:"organizer_name"`
	TicketsSold   int64  `json:"tickets_sold"`
	Revenue       int64  `json:"revenue"`
}
