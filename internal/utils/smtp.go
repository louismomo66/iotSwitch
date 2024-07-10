package utils

import (
	"fmt"
	"iot_switch/internal/config"
	"log"
	"net/smtp"
)

func SendEmail(to, subject, body string) error {
	conf := config.LoadConfig()
	from := conf.SMTPUser
	password := conf.SMTPPass

	// Setup the SMTP configuration
	smtpHost := conf.SMTPHost
	smtpPort := conf.SMTPPort

	msg := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s\r\n", from, to, subject, body))

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		log.Println("Error sending email:", err)
		return err
	}

	log.Println("Email sent successfully to", to)
	return nil
}
