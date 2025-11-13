package mail

import (
	"errors"
	"fmt"
	"net/smtp"
	"os"
)

// mailSender is responsible for sending emails using SMTP.
type MailSender struct {
	from     string
	password string
	host     string
	port     string
}

// InitMailSender initializes a MailSender with SMTP configuration from environment variables.
func InitMailSender() *MailSender {
	return &MailSender{
		from:     os.Getenv("SMTP_FROM"),
		password: os.Getenv("SMTP_PASSWORD"),
		host:     os.Getenv("SMTP_HOST"),
		port:     os.Getenv("SMTP_PORT"),
	}
}

// SendEmail sends an email using passed SendMailCommand.
func (mailSender *MailSender) SendEmail(cmd SendMailCommand) error {
	from, password, host, port := mailSender.from, mailSender.password, mailSender.host, mailSender.port

	if from == "" || password == "" || host == "" || port == "" {
		return errors.New("SMTP configuration missing. Please set SMTP_FROM, SMTP_PASSWORD, SMTP_HOST, SMTP_PORT")
	}

	to := cmd.To
	subject := cmd.Title
	body := cmd.Body

	// Build message
	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=utf-8\r\n\r\n%s",
		from, to, subject, body,
	)

	addr := fmt.Sprintf("%s:%s", host, port)
	auth := smtp.PlainAuth("", from, password, host)

	if err := smtp.SendMail(addr, auth, from, []string{to}, []byte(msg)); err != nil {
		return fmt.Errorf("failed to send email to %s: %w", to, err)
	}

	return nil
}
