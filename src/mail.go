package gocontact

import (
	"fmt"
	"net/smtp"
)

// M stores required mail credentials and fields on initialize
type M struct {
	sender   string
	password string
	smtp     string
	port     int
	to       string
}

var m M

// InitMail should be called in main to load mail credentials
func InitMail(sender string, password string, smtp string, port int, to string) {
	m = M{sender, password, smtp, port, to}
}

// SendContactMail sends an email given subject and message body
func SendContactMail(subject string, name string, email string, body string) error {
	subject = formatSubject(subject)
	body = formatBody(name, email, body)
	msg := formatMessage(subject, body)

	err := smtp.SendMail(fmt.Sprintf("%s:%d", m.smtp, m.port),
		smtp.PlainAuth("", m.sender, m.password, m.smtp),
		m.sender, []string{m.to}, []byte(msg))

	if err != nil {
		return err
	}

	return nil
}

func formatSubject(subject string) string {
	return "Hello: " + subject
}

func formatBody(name string, email string, body string) string {
	return fmt.Sprintf("Name: %s\nEmail: %s\nMessage:\n%s\n", name, email, body)
}

func formatMessage(subject string, body string) []byte {
	message := fmt.Sprintf("From: %s\nTo: %s\n Subject: %s\n\n%s", m.sender, m.to, subject, body)
	return []byte(message)
}
