package domain

import (
	"time"
)

type User struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Mobile        string    `json:"mobile"`
	Dob           string    `json:"dob"`
	Gender        string    `json:"gender"`
	ProfileImage  string    `json:"profile_image"`
	Locations     []string  `json:"locations"`
	IsActive      bool      `json:"is_active"`
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type OrganizerDetail struct {
	UserID           string `json:"user_id"`
	OrganizationName string `json:"organization_name"`
	Address          string `json:"address"`
}

type OrganizerDetails struct {
	User
	OrganizerDetail
	TotalBookings int64 `json:"total_bookings"`
	TotalEvents   int64 `json:"total_events"`
	WalletBalance int64 `json:"wallet_balance"`
	TotalRevenue  int64 `json:"total_revenue"`
}
