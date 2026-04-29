package domain

import "time"

type EventDetails struct {
	EventID			string		`json:"event_id"`
	Description		string		`json:"description"`
	VenueAddress		string		`json:"venue_address"`
	MapURL			string		`json:"map_url"`
	TotalCapacity		int		`json:"total_capacity"`
	AvailableCapacity	int		`json:"available_capacity"`
	Rating			float64		`json:"rating"`
	TotalReviews		int		`json:"total_reviews"`
	TermsAndConditions	string		`json:"terms_and_conditions"`
	CreatedAt		time.Time	`json:"created_at"`
	UpdatedAt		time.Time	`json:"updated_at"`
}
