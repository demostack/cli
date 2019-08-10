package email

import (
	"fmt"
	"log"

	"github.com/demostack/cli/pkg/validate"

	"github.com/manifoldco/promptui"
)

// SendSMTP will send an email via SMTP.
func (c Config) SendSMTP(passphrase *validate.Passphrase) {
	fmt.Println("Send an email via SMTP.")

	// Load the vars.
	f := new(File)
	err := c.store.Load(f, c.Prefix)
	if err != nil {
		log.Fatalln(err)
	}

	if f.SMTP.Host == "" {
		log.Fatalln("No SMTP connection info found. Please set with: demostack config smtp")
	}

	prompt := promptui.Prompt{
		Label:    "To Address (string)",
		Default:  "",
		Validate: validate.RequireString,
	}
	to := validate.Must(prompt.Run())

	prompt = promptui.Prompt{
		Label:   "Subject (string, optional)",
		Default: "",
	}
	subject := validate.Must(prompt.Run())

	prompt = promptui.Prompt{
		Label:   "Body (string, optional)",
		Default: "",
	}
	body := validate.Must(prompt.Run())

	file := ""
	for true {
		prompt = promptui.Prompt{
			Label:   "File to Attach (string, optional)",
			Default: "",
		}
		file = validate.Must(prompt.Run())

		if len(file) == 0 {
			fmt.Println("Skipping file.")
			break
		}

		err := validate.RequireFile(file)
		if err == nil {
			file = validate.ExpandPath(file)
			break
		}

		fmt.Println("File cannot be found.")
	}

	dec, err := f.SMTP.Decrypted(passphrase)
	if err != nil {
		log.Println(err)
	}

	err = dec.SendMail([]string{to}, subject, body, dec.SkipVerify, file)
	if err != nil {
		log.Println(err)
	}

	fmt.Println("Email sent.")
}
