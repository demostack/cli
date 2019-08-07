package awslib

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
)

// SendMessage .
func SendMessage(c Storage, phone, message string) error {
	// Create a new session.
	sess := Session(c)

	// Create an SNS client.
	svc := sns.New(sess)

	// Send the text message.
	_, err := svc.Publish(&sns.PublishInput{
		Message:     aws.String(message),
		PhoneNumber: aws.String(phone),
	})

	return err
}
