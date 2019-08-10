package awslib

import (
	"github.com/aws/aws-sdk-go/service/sts"
)

// GetCallerIdentity .
func GetCallerIdentity(c Storage) (*sts.GetCallerIdentityOutput, error) {
	// Create a new session.
	sess := Session(c)

	// Create a STS client.
	svc := sts.New(sess)
	out, err := svc.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, err
	}

	return out, nil
}
