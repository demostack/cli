package appenv

import (
	"fmt"
	"log"

	"github.com/demostack/cli/pkg/secure"
	"github.com/demostack/cli/pkg/validate"

	"github.com/manifoldco/promptui"
)

// Set a new secure environment variable.
func (c Config) Set(passphrase *validate.Passphrase) {
	fmt.Println("Set or update a secure environment variable.")

	// Load the vars.
	envFile := new(EnvFile)
	err := c.store.Load(envFile, c.Prefix)
	if err != nil {
		log.Fatalln(err)
	}

	// App name.
	prompt := promptui.Prompt{
		Label:    "App name (string)",
		Default:  "",
		Validate: validate.RequireString,
	}
	appName := validate.Must(prompt.Run())

	// Profile name.
	prompt = promptui.Prompt{
		Label:    "Profile name (string)",
		Default:  "",
		Validate: validate.RequireString,
	}
	profileName := validate.Must(prompt.Run())

	varName := ""
	for true {
		// Variable name.
		prompt = promptui.Prompt{
			Label:    "Var name (string)",
			Default:  "",
			Validate: validate.RequireString,
		}
		varName = validate.Must(prompt.Run())

		_, ok := envFile.Var(appName, profileName, varName)
		if !ok {
			break
		} else {
			pSelect := promptui.Select{
				Label: "Existing var name found, overwrite it (select)",
				Items: []string{
					"yes",
					"no",
				},
			}
			replace := validate.MustSelect(pSelect.Run())
			if replace == "yes" {
				break
			}
		}
	}

	env := EnvVar{}

	pSelect := promptui.Select{
		Label: "Encrypt var value (select)",
		Items: []string{
			"yes",
			"no",
		},
	}
	encryptValue := validate.MustSelect(pSelect.Run())

	if encryptValue == "yes" {
		env.Encrypted = true

		// Environment value.
		prompt = promptui.Prompt{
			Label:    "Var value (string)",
			Default:  "",
			Mask:     '*',
			Validate: validate.RequireString,
		}
		newValue := validate.Must(prompt.Run())

		env.Value, err = secure.Encrypt(newValue, passphrase.Password())
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		env.Encrypted = false

		// Environment value.
		prompt = promptui.Prompt{
			Label:    "Var value (string)",
			Default:  "",
			Validate: validate.RequireString,
		}
		env.Value = validate.Must(prompt.Run())
	}

	// Store the value.
	envFile.SetVar(appName, profileName, varName, env)

	err = c.store.Save(envFile, c.Prefix)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Added:", varName)
}
