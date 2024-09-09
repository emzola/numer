package email

import (
	"log"
	"net/smtp"
)

type EmailSender struct {
	host     string
	port     string
	username string
	password string
}

func NewEmailSender(host, port, username, password string) *EmailSender {
	return &EmailSender{
		host:     host,
		port:     port,
		username: username,
		password: password,
	}
}

func (e *EmailSender) SendEmail(to string, subject string, body string) error {
	from := e.username
	pass := e.password
	smtpHost := e.host
	smtpPort := e.port

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	err := smtp.SendMail(smtpHost+":"+smtpPort,
		smtp.PlainAuth("", from, pass, smtpHost),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return err
	}

	log.Print("Email sent successfully to ", to)
	return nil
}
