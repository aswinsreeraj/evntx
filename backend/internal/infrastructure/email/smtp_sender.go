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
	body := fmt.Sprintf("Your OTP is: %s\nIt expires in 5 minutes.", otp)

	message := []byte(
		fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body),
	)

	addr := fmt.Sprintf("%s:%s", host, port)

	return smtp.SendMail(addr, auth, username, []string{email}, message)
}
