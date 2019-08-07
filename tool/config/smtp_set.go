package config

import (
	"fmt"
	"log"
	"strconv"

	"github.com/demostack/cli/pkg/sendmail"
	"github.com/demostack/cli/pkg/validate"

	"github.com/manifoldco/promptui"
)

// SetSMTP will set the SMTP info.
func (c Config) SetSMTP(f File, passphrase *validate.Passphrase) {
	fmt.Println("Set the SMTP connection information.")

	var err error
	key := sendmail.SMTP{}

	prompt := promptui.Prompt{
		Label:    "Host (string)",
		Default:  "",
		Validate: validate.RequireString,
	}
	key.Host = validate.Must(prompt.Run())

	prompt = promptui.Prompt{
		Label:    "Port (int)",
		Default:  "",
		Validate: validate.RequireInt,
	}
	rawPort := validate.Must(prompt.Run())
	key.Port, err = strconv.Atoi(rawPort)
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
	key.SkipVerify, err = strconv.ParseBool(confirm)
	if err != nil {
		log.Fatalln(err)
	}

	prompt = promptui.Prompt{
		Label:    "From Address (email)",
		Default:  "",
		Validate: validate.RequireString,
	}
	key.From = validate.Must(prompt.Run())

	prompt = promptui.Prompt{
		Label:    "Username (string)",
		Default:  "",
		Validate: validate.RequireString,
	}
	key.Username = validate.Must(prompt.Run())

	prompt = promptui.Prompt{
		Label:   "Password (secret, optional)",
		Default: "",
		Mask:    '*',
	}
	key.Password = validate.Must(prompt.Run())

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

		err := key.SendMail([]string{to}, "demostack Confirmation Email.",
			"The SMTP settings are configured properly.", key.SkipVerify)
		if err != nil {
			log.Println(err)
		}
	}

	enc, err := key.Encrypted(passphrase)
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
