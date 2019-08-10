package config

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/demostack/cli/pkg/awslib"
	"github.com/demostack/cli/pkg/validate"
)

// LoginTest .
func (c Config) LoginTest(f File, passphrase *validate.Passphrase) {
	fmt.Println("Get AWS credentials info")

	if len(f.Storage.AWS.AccessKeyID) == 0 {
		log.Fatalln("AWS credentials not found.")
	}

	dec, err := f.Storage.AWS.Decrypted(passphrase.Password())
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf(TimeRemaining(f.Storage.AWS.Expiration))

	out, err := awslib.GetCallerIdentity(dec)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("User ID:", aws.StringValue(out.UserId))
	fmt.Println("Account:", aws.StringValue(out.Account))
	fmt.Println("ARN:", aws.StringValue(out.Arn))

}
