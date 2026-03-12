package repository

type EmailSender interface {
	SendOTP(email string, otp string) error
}
