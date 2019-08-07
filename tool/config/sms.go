package config

import (
	"fmt"
	"log"

	"github.com/demostack/cli/pkg/awslib"
	"github.com/demostack/cli/pkg/validate"

	"github.com/manifoldco/promptui"
)

// SendSMS will send a text message via AWS SNS.
func (c Config) SendSMS(f File, passphrase *validate.Passphrase) {
	fmt.Println("Send a text message via AWS SNS.")

	if f.Storage.AWS.AccessKeyID == "" {
		log.Fatalln("No AWS access keys found. Please set with: demostack config storage aws")
	}

	prompt := promptui.Prompt{
		Label:     "Phone Number (string, format: +12225559999)",
		Default:   "+1",
		AllowEdit: true,
		Validate:  validate.RequireString,
	}
	to := validate.Must(prompt.Run())

	prompt = promptui.Prompt{
		Label:   "Message (string, optional)",
		Default: "",
	}
	message := validate.Must(prompt.Run())

	dec, err := f.Storage.AWS.Decrypted(passphrase.Password())
	if err != nil {
		log.Println(err)
	}

	err = awslib.SendMessage(dec, to, message)
	if err != nil {
		log.Println(err)
	}

	fmt.Println("Text message sent.")
}
