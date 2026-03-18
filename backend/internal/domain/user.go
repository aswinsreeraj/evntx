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
	Gender           string    `json:"gender"`
	ProfileImage     string    `json:"profile_image"`
	OrganizationName string    `json:"organization_name"`
	Locations        []string  `json:"locations"`
	IsActive         bool      `json:"is_active"`
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
