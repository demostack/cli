package config

import (
	"errors"
	"fmt"

	"github.com/demostack/cli/pkg/validate"
	"github.com/demostack/cli/tool"
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

		// Generate the passphrase and the ID.
		passphrase := validate.GeneratePassphrase()
		f.ID = passphrase.ID()

		err = c.store.Save(f, c.Prefix)
		if err != nil {
			return f, errors.New("config save error: " + err.Error())
		}
	} else {
		fmt.Printf("Current storage provider: %v.\n", f.Storage.Current)
	}

	return f, nil
}
