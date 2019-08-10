package email

import (
	"fmt"
	"log"
	"strconv"

	"github.com/demostack/cli/pkg/sendmail"
	"github.com/demostack/cli/pkg/validate"

	"github.com/manifoldco/promptui"
)

// SetSMTP will set the SMTP info.
func (c Config) SetSMTP(passphrase *validate.Passphrase) {
	fmt.Println("Set the SMTP connection information.")

	// Load the vars.
	f := new(File)
	err := c.store.Load(f, c.Prefix)
	if err != nil {
		log.Fatalln(err)
	}

	// Decrypted struct.
	dec := sendmail.SMTP{}

	prompt := promptui.Prompt{
		Label:    "Host (string)",
		Default:  "",
		Validate: validate.RequireString,
	}
	dec.Host = validate.Must(prompt.Run())

	prompt = promptui.Prompt{
		Label:    "Port (int)",
		Default:  "",
		Validate: validate.RequireInt,
	}
	rawPort := validate.Must(prompt.Run())
	dec.Port, err = strconv.Atoi(rawPort)
	if err != nil {
		log.Fatalln(err)
	}

	pSelect := promptui.Select{
		Label: "Skip SSL verification (select)",
		Items: []string{
			"false",
			"true",
		},
	}
	confirm := validate.MustSelect(pSelect.Run())
	dec.SkipVerify, err = strconv.ParseBool(confirm)
	if err != nil {
		log.Fatalln(err)
	}

	prompt = promptui.Prompt{
		Label:    "From Address (email)",
		Default:  "",
		Validate: validate.RequireString,
	}
	dec.From = validate.Must(prompt.Run())

	prompt = promptui.Prompt{
		Label:    "Username (string)",
		Default:  "",
		Validate: validate.RequireString,
	}
	dec.Username = validate.Must(prompt.Run())

	prompt = promptui.Prompt{
		Label:   "Password (secret, optional)",
		Default: "",
		Mask:    '*',
	}
	dec.Password = validate.Must(prompt.Run())

	pSelect = promptui.Select{
		Label: "Send test email (select)",
		Items: []string{
			"yes",
			"no",
		},
	}
	confirm = validate.MustSelect(pSelect.Run())

	if confirm == "yes" {
		prompt = promptui.Prompt{
			Label:    "To Address (email)",
			Default:  "",
			Validate: validate.RequireString,
		}
		to := validate.Must(prompt.Run())

		err := dec.SendMail([]string{to}, "demostack Confirmation Email.",
			"The SMTP settings are configured properly.", dec.SkipVerify)
		if err != nil {
			log.Println(err)
		}
	}

	enc, err := dec.Encrypted(passphrase)
	if err != nil {
		log.Fatalln(err)
	}

	f.SMTP = enc

	err = c.store.Save(f, c.Prefix)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Saved SMTP connection information.")
}
