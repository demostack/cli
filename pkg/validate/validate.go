package validate

import (
	"errors"
	"log"

	"github.com/demostack/cli/pkg/secure"
)

// Must fails or returns a string from a prompt.
func Must(result string, err error) string {
	if err != nil {
		log.Fatalf("Prompt cancelled %v\n", err)
		return ""
	}

	return result
}

// MustSelect fails or returns a string from a select prompt.
func MustSelect(i int, result string, err error) string {
	if err != nil {
		log.Fatalf("Prompt cancelled %v\n", err)
		return ""
	}

	return result
}

// RequireString ensures the string is not empty.
func RequireString(input string) error {
	if len(input) < 1 {
		return errors.New("Value required")
	}
	return nil
}

// EncryptionKey ensures the string can be encrypted with the password.
func EncryptionKey(input string) error {
	_, err := secure.Encrypt("This is the test to encrypt.", input)
	return err
}
