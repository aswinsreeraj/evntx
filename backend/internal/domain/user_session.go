package domain

import "time"

type UserSession struct {
	ID               string
	UserID           string
	RefreshTokenHash string
	UserAgent        string
	IPAddress        string
	ExpiresAt        time.Time
	Revoked          bool
	CreatedAt        time.Time
}
