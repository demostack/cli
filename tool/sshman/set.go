package sshman

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/demostack/cli/pkg/secure"
	"github.com/demostack/cli/pkg/validate"

	"github.com/manifoldco/promptui"
)

// Set a new SSH item.
func (c Config) Set() {
	fmt.Println("Set or update an SSH entry.")

	// Load the entries.
	sshFile := new(SSHFile)
	err := c.store.LoadFile(sshFile, c.Prefix)
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

	password := ""

	if len(sshFile.Arr) > 0 {
		// If a password already exists, verify it by performing a decrypt.
		password, _ = DecryptValue(sshFile.Arr[0].PrivateKey)
	} else {
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

		fmt.Println("Passwords match.")
	}

	// Environment value.
	prompt = promptui.Prompt{
		Label:    "Private key path (string)",
		Default:  "",
		Validate: validate.RequirePEM,
	}
	priKeyFile := validate.Must(prompt.Run())

	// Read the file.
	b, err := ioutil.ReadFile(validate.ExpandPath(priKeyFile))
	if err != nil {
		log.Fatalln(err)
	}

	// Generate the private key that is used on all other steps.
	pri, err := secure.ParsePrivatePEM(string(b))
	if err != nil {
		log.Fatalln(err)
	}

	ent.PrivateKey, err = secure.Encrypt(secure.PrivateKeyToPEM(pri), password)
	if err != nil {
		log.Fatalln(err)
	}

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

	fmt.Printf("Added: %v@%v\n", ent.User, ent.Hostname)
}
