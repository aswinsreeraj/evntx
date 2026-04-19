package email

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/aswinsreeraj/evntx/internal/repository"
)

type SMTPSender struct{}

func NewSMTPSender() repository.EmailSender {
	return &SMTPSender{}
}

func (s *SMTPSender) SendOTP(email, otp string) error {

	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	auth := smtp.PlainAuth("", username, password, host)

	subject := "Your EVNTX OTP Code"

	body := fmt.Sprintf(
		"Welcome to EVNTX!\r\n\r\n"+
			"Your One-Time Password (OTP) for verification is: %s\r\n\r\n"+
			"This OTP is valid for 5 minutes. Please do not share this code with anyone.\r\n\r\n"+
			"If you did not request this, please ignore this email.\r\n\r\n"+
			"Best regards,\r\n"+
			"The EVNTX Team",
		otp,
	)

	message := []byte(
		fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body),
	)

	addr := fmt.Sprintf("%s:%s", host, port)

	return smtp.SendMail(addr, auth, username, []string{email}, message)
}

func (s *SMTPSender) SendOrganizerApproval(email, name string) error {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	auth := smtp.PlainAuth("", username, password, host)

	subject := "Your EVNTX Organizer Access Is Approved"
	body := fmt.Sprintf(
		"Hello %s,\r\n\r\n"+
			"Great news! Your organizer account request has been approved by the EVNTX team.\r\n\r\n"+
			"You can now log in and access the organizer dashboard to create and manage events.\r\n\r\n"+
			"If you did not request this, please contact support immediately.\r\n\r\n"+
			"Best regards,\r\n"+
			"The EVNTX Team",
		name,
	)

	message := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body))
	addr := fmt.Sprintf("%s:%s", host, port)
	return smtp.SendMail(addr, auth, username, []string{email}, message)
}
