package sshman

import (
	"fmt"

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
		Prefix: "ssh",
	}
}

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

// Encrypted returns an encrypted object.
func (f SSHFile) Encrypted(passphrase *validate.Passphrase) SSHFile {
	for i := 0; i < len(f.Arr); i++ {
		v, err := secure.Encrypt(f.Arr[i].PrivateKey, passphrase.Password())
		if err != nil {
			fmt.Println("Could not encrypt private key:", f.Arr[i].Name)
			f.Arr[i].PrivateKey = ""
		} else {
			f.Arr[i].PrivateKey = v
		}
	}
	return f
}

// Decrypted returns an decrypted object.
func (f SSHFile) Decrypted(passphrase *validate.Passphrase) SSHFile {
	for i := 0; i < len(f.Arr); i++ {
		v, err := secure.Decrypt(f.Arr[i].PrivateKey, passphrase.Password())
		if err != nil {
			fmt.Println("Could not decrypt private key:", f.Arr[i].Name)
			f.Arr[i].PrivateKey = ""
		} else {
			f.Arr[i].PrivateKey = v
		}
	}
	return f
}
