package validate

import (
	"fmt"

	"github.com/demostack/cli/pkg/secure"

	"github.com/manifoldco/promptui"
)

// DecryptValue will verify the password is correct and return the password
// and the decrypted data.
func DecryptValue(encryptedValue string) (string, string) {
	for true {
		// If a password already exists, verify it.
		prompt := promptui.Prompt{
			Label:    "Password required (secure)",
			Default:  "",
			Mask:     '*',
			Validate: EncryptionKey,
		}
		password := Must(prompt.Run())

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
