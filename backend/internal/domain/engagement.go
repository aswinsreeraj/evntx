package domain

import (
	"time"
)

type EngagementEventType string

const (
	EngagementEventPageView			EngagementEventType	= "page_view"
	EngagementEventEventView		EngagementEventType	= "event_view"
	EngagementEventTicketSelected		EngagementEventType	= "ticket_selected"
	EngagementEventCheckoutStarted		EngagementEventType	= "checkout_started"
	EngagementEventSuccessfulBooking	EngagementEventType	= "successful_booking"

	PlatformEventID	= "00000000-0000-0000-0000-000000000000"
)

type VisitorSession struct {
	ID		string		`json:"id"`
	UserID		*string		`json:"user_id,omitempty"`
	IPAddress	string		`json:"ip_address"`
	UserAgent	string		`json:"user_agent"`
	CreatedAt	time.Time	`json:"created_at"`
	LastSeenAt	time.Time	`json:"last_seen_at"`
}

type EngagementEvent struct {
	ID		string			`json:"id"`
	UserID		*string			`json:"user_id,omitempty"`
	SessionID	string			`json:"session_id"`
	EventID		*string			`json:"event_id,omitempty"`
	EventType	EngagementEventType	`json:"event_type"`
	Metadata	string			`json:"metadata"`
	IPAddress	string			`json:"ip_address"`
	UserAgent	string			`json:"user_agent"`
	CreatedAt	time.Time		`json:"created_at"`
}

type EventEngagementDaily struct {
	ID			string		`json:"id"`
	EventID			string		`json:"event_id"`
	Date			time.Time	`json:"date"`
	Visitors		int		`json:"visitors"`
	PageViews		int		`json:"page_views"`
	EventViews		int		`json:"event_views"`
	TicketsSelected		int		`json:"tickets_selected"`
	CheckoutStarted		int		`json:"checkout_started"`
	SuccessfulBookings	int		`json:"successful_bookings"`
	CreatedAt		time.Time	`json:"created_at"`
}
