package validate

import (
	"fmt"
	"log"

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

// GeneratePassphrase will return a new passphrase.
func GeneratePassphrase() *Passphrase {
	password := ""

	for true {
		// If no password is set, create a new one.
		prompt := promptui.Prompt{
			Label:    "New password (secure)",
			Default:  "",
			Mask:     '*',
			Validate: EncryptionKey,
		}
		password = Must(prompt.Run())

		// If no password is set, create a new one.
		prompt = promptui.Prompt{
			Label:    "Verify new password (password)",
			Default:  "",
			Mask:     '*',
			Validate: EncryptionKey,
		}
		verifyPassword := Must(prompt.Run())

		if password == verifyPassword {
			break
		}

		fmt.Println("Password don't match, please try again.")
	}

	id, err := secure.UUID()
	if err != nil {
		log.Fatalln("cannot generate UUID: " + err.Error())
	}

	enc, err := secure.Encrypt(id, password)
	if err != nil {
		log.Fatalln("cannot encrypt UUID: " + err.Error())
	}

	p := NewPassphrase(enc)
	p.password = password
	p.saved = true

	return p
}

// ID returns the encrypted value.
func (p *Passphrase) ID() string {
	return p.encryptedValue
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
			break
		}
	}

	return p.password
}
