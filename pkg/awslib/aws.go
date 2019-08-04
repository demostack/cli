package awslib

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sts"
)

// Storage is S3 bucket storage.
type Storage struct {
	AccessKeyID     string `json:"id"`
	SecretAccessKey string `json:"secret"`
	SessionToken    string `json:"session"`
	Region          string `json:"region"`
	Bucket          string `json:"bucket"`
}

// Session returns an AWS session.
func Session(c Storage) *session.Session {
	return session.Must(session.NewSession(&aws.Config{
		Region: aws.String(c.Region),
		Credentials: credentials.NewStaticCredentialsFromCreds(credentials.Value{
			AccessKeyID:     c.AccessKeyID,
			SecretAccessKey: c.SecretAccessKey,
			SessionToken:    c.SessionToken, // This can be an empty string.
		}),
	}))
}

// AccountNumber .
func AccountNumber(c Storage) (string, error) {
	sess := Session(c)

	// Get the identity of the current user.
	svc := sts.New(sess)
	out, err := svc.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return "", err
	}

	return aws.StringValue(out.Account), nil
}

// CreateBucket .
func CreateBucket(c Storage) error {
	sess := Session(c)
	svc := s3.New(sess)

	_, err := svc.GetBucketVersioning(&s3.GetBucketVersioningInput{
		Bucket: aws.String(c.Bucket),
	})
	if err != nil {
		fmt.Println(err)

		_, err = svc.CreateBucket(&s3.CreateBucketInput{
			Bucket: aws.String(c.Bucket),
		})
		if err != nil {
			return err
		}
		fmt.Println("Bucket created:", c.Bucket)
	} else {
		fmt.Println("Bucket found:", c.Bucket)
	}

	return nil
}
