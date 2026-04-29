package domain

import "time"

type EmailOTP struct {
	ID		string
	Email		string
	OTPHash		string
	ExpiresAt	time.Time
	Consumed	bool
	CreatedAt	time.Time
}
