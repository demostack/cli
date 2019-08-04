package sshman

import (
	"fmt"

	"github.com/demostack/cli/pkg/secure"
	"github.com/demostack/cli/pkg/validate"

	"github.com/manifoldco/promptui"
)

// SSHFile is an SSH config file.
type SSHFile struct {
	Arr []SSHEntry `json:"entries"`
}

// SSHEntry represents an SSH entry.
type SSHEntry struct {
	Name       string `json:"name"`
	Hostname   string `json:"hostname"`
	User       string `json:"user"`
	PrivateKey string `json:"private_key"`
}

// DecryptValue will verify the password is correct and return the password
// and the decrypted data.
func DecryptValue(encryptedValue string) (string, string) {
	for true {
		// If a password already exists, verify it.
		prompt := promptui.Prompt{
			Label:    "Password required (secure)",
			Default:  "",
			Mask:     '*',
			Validate: validate.EncryptionKey,
		}
		password := validate.Must(prompt.Run())

		dec, err := secure.Decrypt(encryptedValue, password)

		if err != nil {
			fmt.Println("Password is not correct, please try again.")
		} else {
			fmt.Println("Password correct.")
			return password, dec
		}
	}

	return "", ""
}
