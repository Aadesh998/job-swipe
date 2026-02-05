package utils

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func SendEmail(to string, subject string, body string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPortStr := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASSWORD")

	if smtpHost == "" {
		fmt.Printf("Mock Email to %s\nSubject: %s\nBody: %s\n", to, subject, body)
		return nil
	}

	smtpPort, _ := strconv.Atoi(smtpPortStr)

	m := gomail.NewMessage()
	m.SetHeader("From", smtpUser)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)

	return d.DialAndSend(m)
}
