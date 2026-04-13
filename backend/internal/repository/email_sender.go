package repository

type EmailSender interface {
	SendOTP(email string, otp string) error
	SendOrganizerApproval(email string, name string) error
}
