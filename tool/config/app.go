package config

import (
	"log"

	"github.com/demostack/cli/pkg/awslib"
	"github.com/demostack/cli/pkg/secure"
	"github.com/demostack/cli/pkg/sendmail"
	"github.com/demostack/cli/pkg/validate"
)

// File is the demostack config file.
type File struct {
	ID      string        `json:"id"`
	Storage Storage       `json:"storage"`
	SMTP    sendmail.SMTP `json:"smtp"`
}

// Storage is the storage of the config files.
type Storage struct {
	// Current supports the following values: filesystem, aws.
	Current    string         `json:"current"`
	AWS        awslib.Storage `json:"aws"`
	Filesystem Filesystem     `json:"filesystem"`
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

	f.SMTP, err = f.SMTP.Encrypted(passphrase)
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

	f.SMTP, err = f.SMTP.Decrypted(passphrase)
	if err != nil {
		log.Fatalln(err)
	}

	return f
}
