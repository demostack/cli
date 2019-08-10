package validate

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/demostack/cli/pkg/secure"
)

// Must fails or returns a string from a prompt.
func Must(result string, err error) string {
	if err != nil {
		log.Fatalf("Prompt cancelled %v\n", err)
		return ""
	}

	return strings.TrimSpace(result)
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
	if len(strings.TrimSpace(input)) < 1 {
		return errors.New("Value required")
	}
	return nil
}

// RequireInt ensures the input is a number.
func RequireInt(input string) error {
	_, err := strconv.Atoi(input)
	if err != nil {
		return errors.New("Int required")
	}
	return nil
}

// RequireAWSSessionInt ensures the input is a number.
func RequireAWSSessionInt(input string) error {
	n, err := strconv.Atoi(input)
	if err != nil {
		return errors.New("Int required")
	}

	if n < 15 || n > 2160 {
		return errors.New("Must be between 15 and 2160")
	}

	return nil
}

// EncryptionKey ensures the string can be encrypted with the password.
func EncryptionKey(input string) error {
	_, err := secure.Encrypt("This is the test to encrypt.", input)
	return err
}

// ExpandPath will replace the tilda with the user's home directory.
func ExpandPath(relpath string) string {
	if strings.HasPrefix(relpath, "~/") {
		u, err := user.Current()
		if err != nil {
			return relpath
		}

		return filepath.Join(u.HomeDir, relpath[1:])
	}
	return relpath
}

// RequireFile ensures a file exists.
func RequireFile(input string) error {
	info, err := os.Stat(ExpandPath(input))
	if os.IsNotExist(err) {
		return err
	}

	if !info.IsDir() {
		return nil
	}

	return errors.New("found not found")
}

// RequirePEM ensures a file exists and is a PEM.
func RequirePEM(input string) error {
	b, err := ioutil.ReadFile(ExpandPath(input))
	if err != nil {
		return err
	}

	_, err = secure.ParsePrivatePEM(string(b))
	if err != nil {
		return err
	}

	return nil
}
