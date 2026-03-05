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
