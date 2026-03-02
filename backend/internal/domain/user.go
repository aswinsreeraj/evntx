package domain

import (
	"time"
)

type User struct {
	ID            string
	Name          string
	Email         string
	IsActive      bool
	EmailVerified bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
