package email

import (
	"log"

	"github.com/demostack/cli/pkg/sendmail"
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
		Prefix: "email",
	}
}

// File is the demostack config file.
type File struct {
	SMTP sendmail.SMTP `json:"smtp"`
}

// Encrypted .
func (f File) Encrypted(passphrase *validate.Passphrase) File {
	var err error

	f.SMTP, err = f.SMTP.Encrypted(passphrase)
	if err != nil {
		log.Fatalln(err)
	}

	return f
}

// Decrypted .
func (f File) Decrypted(passphrase *validate.Passphrase) File {
	var err error

	f.SMTP, err = f.SMTP.Decrypted(passphrase)
	if err != nil {
		log.Fatalln(err)
	}

	return f
}
