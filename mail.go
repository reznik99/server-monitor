package main

import (
	"fmt"
	"os"

	"github.com/wneessen/go-mail"
)

func SendEmailAlert() error {
	var subject = "RaspberrypPI Tor Relay server-monitor alert"
	var body = "It's on fire!"
	var fromEmail = os.Getenv("SOURCE_EMAIL_ADDRESS")
	var toEmail = os.Getenv("TARGET_EMAIL_ADDRESS")

	// Create email object
	// TODO: use HTML template?
	msg := mail.NewMsg()
	if err := msg.From(fromEmail); err != nil {
		return fmt.Errorf("invalid from email address '%s': %s", fromEmail, err)
	}
	if err := msg.To(toEmail); err != nil {
		return fmt.Errorf("invalid to email address '%s': %s", toEmail, err)
	}
	msg.Subject(subject)
	msg.SetBodyString(mail.TypeTextPlain, body)

	// Send email
	if err := msg.WriteToSendmail(); err != nil {
		return fmt.Errorf("send mail err: %s", err)
	}

	return nil
}
