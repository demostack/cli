package secureenv

import (
	"fmt"

	"github.com/demostack/cli/pkg/secure"
	"github.com/demostack/cli/pkg/validate"

	"github.com/manifoldco/promptui"
)

// EnvFile is an environment config file.
type EnvFile struct {
	App string   `json:"app"`
	Arr []EnvVar `json:"vars"`
}

// Strings returns a string array of environment variables.
func (ef EnvFile) Strings(password string) []string {
	arr := make([]string, 0)
	for _, v := range ef.Arr {
		s := v.String(password)
		if len(s) > 0 {
			arr = append(arr, s)
		}
	}
	return arr
}

// HasEncryptedValues returns true if there are encrypted values and returns
// one of the encrypted values.
func (ef EnvFile) HasEncryptedValues() (bool, string) {
	for _, v := range ef.Arr {
		if v.Encrypted {
			return true, v.Value
		}
	}
	return false, ""
}

// EnvVar represents an environment variable.
type EnvVar struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	Encrypted bool   `json:"encrypted"`
}

// String returns the name and value in this format: name=value.
func (ev EnvVar) String(password string) string {
	if ev.Encrypted {
		v, err := secure.Decrypt(ev.Value, password)
		if err != nil {
			fmt.Println("Could not decrypt var:", ev.Name)
			return ""
		}
		return fmt.Sprintf("%v=%v", ev.Name, v)
	}
	return fmt.Sprintf("%v=%v", ev.Name, ev.Value)
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
