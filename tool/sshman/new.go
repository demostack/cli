package sshman

import (
	"fmt"
	"log"

	"github.com/demostack/cli/pkg/secure"
	"github.com/demostack/cli/pkg/validate"

	"github.com/manifoldco/promptui"
)

// New SSH item and generate a private key.
func (c Config) New(passphrase *validate.Passphrase) {
	fmt.Println("Create a new SSH entry and generate a private key.")

	// Load the entries.
	sshFile := new(SSHFile)
	err := c.store.Load(sshFile, c.Prefix)
	if err != nil {
		log.Fatalln(err)
	}

	ent := SSHEntry{}

	// Name.
	prompt := promptui.Prompt{
		Label:    "Name (string)",
		Default:  "",
		Validate: validate.RequireString,
	}
	ent.Name = validate.Must(prompt.Run())

	for true {
		// Hostname.
		prompt := promptui.Prompt{
			Label:    "Hostname (string)",
			Default:  "",
			Validate: validate.RequireString,
		}
		ent.Hostname = validate.Must(prompt.Run())

		// Environment name.
		prompt = promptui.Prompt{
			Label:    "User (string)",
			Default:  "",
			Validate: validate.RequireString,
		}
		ent.User = validate.Must(prompt.Run())

		found := false
		for _, v := range sshFile.Arr {
			if v.Hostname == ent.Hostname && v.User == ent.User {
				found = true

				pSelect := promptui.Select{
					Label: "Existing hostname and user found, overwrite it (select)",
					Items: []string{
						"yes",
						"no",
					},
				}
				replace := validate.MustSelect(pSelect.Run())
				if replace == "yes" {
					found = false
				}
				break
			}
		}
		if !found {
			break
		}
	}

	// Generate a new private key.
	pri, _ := secure.GenerateRSAKeyPair()

	ent.PrivateKey, err = secure.Encrypt(secure.PrivateKeyToPEM(pri), passphrase.Password())
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("SSH key generated.")

	// Loop through to see if there are any matches.
	found := false
	for i := 0; i < len(sshFile.Arr); i++ {
		// If the item is found, replace it.
		if sshFile.Arr[i].Hostname == ent.Hostname && sshFile.Arr[i].User == ent.User {
			found = true
			sshFile.Arr[i] = ent
			break
		}
	}

	// If the item is not found, add it to the end.
	if !found {
		sshFile.Arr = append(sshFile.Arr, ent)
	}

	err = c.store.Save(sshFile, c.Prefix)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Created: %v@%v\n", ent.User, ent.Hostname)
}
