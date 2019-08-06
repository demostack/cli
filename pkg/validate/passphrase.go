package validate

import (
	"fmt"

	"github.com/demostack/cli/pkg/secure"

	"github.com/manifoldco/promptui"
)

// Passphrase .
type Passphrase struct {
	saved          bool
	password       string
	encryptedValue string
}

// NewPassphrase .
func NewPassphrase(encryptedValue string) *Passphrase {
	return &Passphrase{
		encryptedValue: encryptedValue,
	}
}

// Password returns the password if saved or prompted the user to enter it in.
func (p *Passphrase) Password() string {
	if p.saved {
		return p.password
	}

	for true {
		// If a password already exists, verify it.
		prompt := promptui.Prompt{
			Label:    "Password required (secure)",
			Default:  "",
			Mask:     '*',
			Validate: EncryptionKey,
		}
		p.password = Must(prompt.Run())

		_, err := secure.Decrypt(p.encryptedValue, p.password)

		if err != nil {
			fmt.Println("Password is not correct, please try again.")
		} else {
			fmt.Println("Password correct.")
			p.saved = true
		}
	}

	return p.password
}
