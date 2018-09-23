package utils

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

// SendMail sends email to the provided arguments.
func SendMail(sub, msg, toEmail string) {
	auth := smtp.CRAMMD5Auth("", "")
	fromEmail := os.Getenv("FROM_EMAIL")

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	to := []string{toEmail}
	content := fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", toEmail, sub, msg,
	)

	emailMsg := []byte(content)
	err := smtp.SendMail("mailhog:1025", auth, fromEmail, to, emailMsg)
	if err != nil {
		log.Fatal(err)
	}
}
