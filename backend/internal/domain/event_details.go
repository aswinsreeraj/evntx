package domain

import "time"

type EventDetails struct {
	EventID            string
	Description        string
	VenueAddress       string
	MapURL             string
	TotalCapacity      int
	AvailableCapacity  int
	Rating             float64
	TotalReviews       int
	TermsAndConditions string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
