package utils

import (
	"log"
	"net/smtp"
)

// SendMail sends email to the provided arguments.
func SendMail() {
	auth := smtp.CRAMMD5Auth("", "")

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	to := []string{"recipient@example.net"}
	msg := []byte("To: recipient@example.net\r\n" +
		"Subject: discount Gophers!\r\n" +
		"\r\n" +
		"This is the email body.\r\n")
	err := smtp.SendMail("mailhog:1025", auth, "sender@example.org", to, msg)
	if err != nil {
		log.Fatal(err)
	}
}
