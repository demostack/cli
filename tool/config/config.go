package config

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/demostack/cli/pkg/awslib"
	"github.com/demostack/cli/pkg/secure"
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
func (c Config) Load() (File, *validate.Passphrase, error) {
	f := File{}
	err := c.store.Load(&f, c.Prefix)

	var passphrase *validate.Passphrase

	if err != nil {
		fmt.Println("Initialization - Please set a password.")

		// Generate the passphrase and the ID.
		passphrase = validate.GeneratePassphrase()
		f.ID = passphrase.ID()

		// Set the default.
		f.Storage.Current = "filesystem"

		err = c.store.Save(f, c.Prefix)
		if err != nil {
			return f, nil, errors.New("config save error: " + err.Error())
		}
	} else {
		fmt.Printf("Current storage provider: %v.\n", f.Storage.Current)

		if f.Storage.Current == "aws" {
			fmt.Printf(TimeRemaining(f.Storage.AWS.Expiration))
		}
	}

	return f, passphrase, nil
}

// TimeRemaining .
func TimeRemaining(t time.Time) string {
	now := time.Now()
	difference := t.Sub(now)

	if difference < 0 {
		return fmt.Sprintf("AWS access keys expired. Please login again.\n")
	}

	total := int(difference.Seconds())
	days := int(total / (60 * 60 * 24))
	hours := int(total / (60 * 60) % 24)
	minutes := int(total/60) % 60
	seconds := int(total % 60)

	return fmt.Sprintf("Expiration: %d day(s), %d hour(s), %d minute(s), %d second(s)\n", days, hours, minutes, seconds)
}

// File is the demostack config file.
type File struct {
	ID      string  `json:"id"`
	Login   Login   `json:"login"`
	Storage Storage `json:"storage"`
}

// Storage is the storage of the config files.
type Storage struct {
	// Current supports the following values: filesystem, aws.
	Current    string         `json:"current"`
	AWS        awslib.Storage `json:"aws"`
	Filesystem Filesystem     `json:"filesystem"`
}

// Login holds the login credentials.
type Login struct {
	Host     string `json:"host"`
	Username string `json:"username"`
}

// Filesystem is for the local filesystem.
type Filesystem struct{}

// Encrypted .
func (f File) Encrypted(passphrase *validate.Passphrase) File {
	var err error

	f.ID, err = secure.Encrypt(f.ID, passphrase.Password())
	if err != nil {
		log.Fatalln(err)
	}

	f.Storage.AWS, err = f.Storage.AWS.Encrypted(passphrase.Password())
	if err != nil {
		log.Fatalln(err)
	}

	return f
}

// Decrypted .
func (f File) Decrypted(passphrase *validate.Passphrase) File {
	var err error

	f.ID, err = secure.Decrypt(f.ID, passphrase.Password())
	if err != nil {
		log.Fatalln(err)
	}

	f.Storage.AWS, err = f.Storage.AWS.Decrypted(passphrase.Password())
	if err != nil {
		log.Fatalln(err)
	}

	return f
}
