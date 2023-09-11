package utils

import (
	"crypto/tls"
	"fmt"
	"os"

	gomail "gopkg.in/mail.v2"
)

func SendEmail(to []string, subject string, body string) {
	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", "udoncodestudios@gmail.com")

	// Set E-Mail receivers
	m.SetHeader("To", to...)

	// Set E-Mail subject
	m.SetHeader("Subject", subject)

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", body)

	// Settings for SMTP server
	emailPw := os.Getenv("EMAIL_PW")
	d := gomail.NewDialer("smtp.gmail.com", 587, "udoncodestudios@gmail.com", emailPw)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println("[ ERROR ] Error when sending email:", err)
		panic(err)
	}
}
