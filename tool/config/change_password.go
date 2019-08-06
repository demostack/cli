package config

import (
	"fmt"
	"log"

	"github.com/demostack/cli/pkg/validate"
	"github.com/demostack/cli/tool/appenv"
	"github.com/demostack/cli/tool/sshman"
)

// ChangePassword will change the config password.
func (c Config) ChangePassword(f File, passphrase *validate.Passphrase) {
	fmt.Println("Change the config password and re-encrypt all the data.")

	// Prompt for the password.
	passphrase.Password()

	newPassphrase := validate.GeneratePassphrase()

	// Re-encrypt the config.
	decFile := f.Decrypted(passphrase)
	encFile := decFile.Encrypted(newPassphrase)
	err := c.store.Save(encFile, c.Prefix)
	if err != nil {
		log.Fatalln(err)
	}

	// Load the vars. THIS NEEDS TO BE IN A LOOP.
	envFile := new(appenv.EnvFile)
	err = c.store.Load(envFile, "env")
	if err != nil {
		log.Fatalln(err)
	}
	decEnvFile := envFile.Decrypted(passphrase)
	encEnvFile := decEnvFile.Encrypted(newPassphrase)
	err = c.store.Save(encEnvFile, "env")
	if err != nil {
		log.Fatalln(err)
	}

	// Load the entries.
	sshFile := new(sshman.SSHFile)
	err = c.store.Load(sshFile, "ssh")
	if err != nil {
		log.Fatalln(err)
	}
	decSSHFile := sshFile.Decrypted(passphrase)
	encSSHFile := decSSHFile.Encrypted(newPassphrase)
	err = c.store.Save(encSSHFile, "ssh")
	if err != nil {
		log.Fatalln(err)
	}
}
