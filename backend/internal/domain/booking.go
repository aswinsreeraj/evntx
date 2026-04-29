package domain

import "time"

type Booking struct {
	ID			string
	UserID			string
	EventID			string
	Status			string
	TotalAmount		float64
	PlatformFeeValue	float64
	PlatformFeeType		string
	ExpiresAt		time.Time
	CreatedAt		time.Time
}

type BookingTicket struct {
	BookingID	string
	TicketTypeID	string
	Quantity	int
}

type TicketRequest struct {
	TicketTypeID	string	`json:"ticket_type_id" binding:"required"`
	Quantity	int	`json:"quantity" binding:"required,gt=0"`
}

type BookingWithEvent struct {
	BookingID	string
	EventID		string
	EventTitle	string
	EventCity	string
	EventStartTime	time.Time
	Status		string
	TotalAmount	float64
	TicketCount	int
	CreatedAt	time.Time
	CoverImageURL	string
	VenueName	string
	Tags		string
	EventStatus	string
}

type TicketWithEvent struct {
	TicketID	string
	TicketCode	string
	EventID		string
	EventTitle	string
	TicketType	string
	Status		string
	CheckedInAt	*time.Time
}

type TicketCheckIn struct {
	TicketID	string
	TicketCode	string
	Status		string
	CheckedInAt	time.Time
}
