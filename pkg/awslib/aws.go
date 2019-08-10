package awslib

import (
	"time"

	"github.com/demostack/cli/pkg/secure"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

// Storage is S3 bucket storage.
type Storage struct {
	AccessKeyID     string    `json:"id"`
	SecretAccessKey string    `json:"secret"`
	SessionToken    string    `json:"session"`
	Region          string    `json:"region"`
	Bucket          string    `json:"bucket"`
	Expiration      time.Time `json:"expiration"`
}

// Encrypted returns the storage object encrypted.
func (c Storage) Encrypted(password string) (Storage, error) {
	var err error
	c.SecretAccessKey, err = secure.Encrypt(c.SecretAccessKey, password)
	if err != nil {
		return c, err
	}

	c.SessionToken, err = secure.Encrypt(c.SessionToken, password)
	if err != nil {
		return c, err
	}

	return c, nil
}

// Decrypted returns the storage object decrypted.
func (c Storage) Decrypted(password string) (Storage, error) {
	if c.SecretAccessKey == "" {
		return c, nil
	}

	var err error
	c.SecretAccessKey, err = secure.Decrypt(c.SecretAccessKey, password)
	if err != nil {
		return c, err
	}

	c.SessionToken, err = secure.Decrypt(c.SessionToken, password)
	if err != nil {
		return c, err
	}

	return c, nil
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
