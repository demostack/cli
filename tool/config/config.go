package config

import (
	"errors"
	"fmt"

	"github.com/demostack/cli/pkg/secure"
	"github.com/demostack/cli/pkg/validate"
	"github.com/demostack/cli/tool"
	"github.com/manifoldco/promptui"
)

// Config .
type Config struct {
	log   tool.ILogger
	store tool.IStorage

	Prefix string
}

// NewConfig .
func NewConfig(l tool.ILogger, store tool.IStorage) Config {
	return Config{
		log:    l,
		store:  store,
		Prefix: "config",
	}
}

// Load the app configuration file.
func (c Config) Load() (File, error) {
	f := File{}
	err := c.store.Load(&f, c.Prefix)
	if err != nil {
		fmt.Println("Initialization - Please set a password.")
		password := ""

		for true {
			// If no password is set, create a new one.
			prompt := promptui.Prompt{
				Label:    "New password (secure)",
				Default:  "",
				Mask:     '*',
				Validate: validate.EncryptionKey,
			}
			password = validate.Must(prompt.Run())

			// If no password is set, create a new one.
			prompt = promptui.Prompt{
				Label:    "Verify new password (password)",
				Default:  "",
				Mask:     '*',
				Validate: validate.EncryptionKey,
			}
			verifyPassword := validate.Must(prompt.Run())

			if password == verifyPassword {
				break
			}

			fmt.Println("Password don't match, please try again.")
		}

		id, err := secure.UUID()
		if err != nil {
			return f, errors.New("cannot generate UUID: " + err.Error())
		}

		f.ID, err = secure.Encrypt(id, password)
		if err != nil {
			return f, errors.New("cannot encrypt UUID: " + err.Error())
		}

		err = c.store.Save(f, c.Prefix)
		if err != nil {
			return f, errors.New("config save error: " + err.Error())
		}
	} else {
		fmt.Printf("Current storage provider: %v.\n", f.Storage.Current)
	}

	return f, nil
}
