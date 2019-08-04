package config

import (
	"fmt"
	"log"

	"github.com/demostack/cli/pkg/awslib"
	"github.com/demostack/cli/pkg/secure"
	"github.com/demostack/cli/pkg/validate"

	"github.com/manifoldco/promptui"
)

// SetStorageAWS will set the AWS storage.
func (c Config) SetStorageAWS(f File, password string) {
	fmt.Println("Set the storage provider to AWS.")

	key := awslib.Storage{}

	account := ""

	var err error
	for {
		prompt := promptui.Prompt{
			Label:    "AWS Access Key ID (string)",
			Default:  "",
			Validate: validate.RequireString,
		}
		key.AccessKeyID = validate.Must(prompt.Run())

		prompt = promptui.Prompt{
			Label:    "AWS Access Key Secret (secret)",
			Default:  "",
			Mask:     '*',
			Validate: validate.RequireString,
		}
		key.SecretAccessKey = validate.Must(prompt.Run())

		prompt = promptui.Prompt{
			Label:   "AWS Access Session Token (secret, optional)",
			Default: "",
			Mask:    '*',
		}
		key.SessionToken = validate.Must(prompt.Run())

		prompt = promptui.Prompt{
			Label:    "AWS Region (string)",
			Default:  "us-east-1",
			Validate: validate.RequireString,
		}
		key.Region = validate.Must(prompt.Run())

		account, err = awslib.AccountNumber(key)
		if err != nil {
			fmt.Println(err)
			continue
		}

		key.Bucket = fmt.Sprintf("%v-demostack-config", account)

		prompt = promptui.Prompt{
			Label:    "S3 Bucket to store config (string)",
			Default:  key.Bucket,
			Validate: validate.RequireString,
		}
		key.Bucket = validate.Must(prompt.Run())

		break
	}

	// Create the S3 bucket.
	err = awslib.CreateBucket(key)
	if err != nil {
		log.Fatalln(err)
	}

	// Encrypt the secret access key.
	key.SecretAccessKey, err = secure.Encrypt(key.SecretAccessKey, password)
	if err != nil {
		log.Fatalln(err)
	}

	// Encrypt the session token.
	key.SessionToken, err = secure.Encrypt(key.SessionToken, password)
	if err != nil {
		log.Fatalln(err)
	}

	f.Storage.Current = "aws"
	f.Storage.AWS = key

	err = c.store.Save(f, c.Prefix)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Set provider to:", "aws")
}
