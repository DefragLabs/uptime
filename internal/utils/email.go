package utils

import (
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

// SendMail sends email to the provided arguments.
func SendMail(sub, msg, toEmail string) {
	auth := smtp.CRAMMD5Auth("", "")
	fromEmail := os.Getenv("FROM_EMAIL")
	env := os.Getenv("ENV")

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	to := []string{toEmail}
	content := fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", toEmail, sub, msg,
	)

	emailMsg := []byte(content)
	if env == "local" {
		err := smtp.SendMail("mailhog:1025", auth, fromEmail, to, emailMsg)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		awsSESRegion := os.Getenv("AWS_SES_REGION")
		awsSESAccessKeyID := os.Getenv("AWS_SES_ACCESS_KEY")
		awsSESAccessKeySecret := os.Getenv("AWS_SES_ACCESS_SECRET")

		if awsSESRegion == "" || awsSESAccessKeyID == "" || awsSESAccessKeySecret == "" {
			log.Fatal("AWS SES configuration invalid.")
		}

		awsSession := session.New(&aws.Config{
			Region:      aws.String(awsSESRegion),
			Credentials: credentials.NewStaticCredentials(awsSESAccessKeyID, awsSESAccessKeySecret, ""),
		})

		sesSession := ses.New(awsSession)

		sesEmailInput := &ses.SendEmailInput{
			Destination: &ses.Destination{
				ToAddresses: []*string{aws.String(toEmail)},
			},
			Message: &ses.Message{
				Body: &ses.Body{
					Html: &ses.Content{
						Data: aws.String(msg)},
				},
				Subject: &ses.Content{
					Data: aws.String(sub),
				},
			},
			Source: aws.String(fromEmail),
			ReplyToAddresses: []*string{
				aws.String(fromEmail),
			},
		}

		_, err := sesSession.SendEmail(sesEmailInput)
		if err != nil {
			log.Fatal(err)
		}
	}
}
