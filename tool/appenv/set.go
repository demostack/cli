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

	// App name.
	prompt := promptui.Prompt{
		Label:    "App name (string)",
		Default:  "",
		Validate: validate.RequireString,
	}
	app := validate.Must(prompt.Run())

	// Load the vars.
	envFile := new(EnvFile)
	envFile.App = app
	err := c.store.LoadFile(envFile, c.Prefix, app)
	if err != nil {
		log.Fatalln(err)
	}

	env := EnvVar{}

	for true {
		// Environment name.
		prompt = promptui.Prompt{
			Label:    "Var name (string)",
			Default:  "",
			Validate: validate.RequireString,
		}
		env.Name = validate.Must(prompt.Run())

		found := false
		for _, v := range envFile.Arr {
			if v.Name == env.Name {
				found = true

				pSelect := promptui.Select{
					Label: "Existing var name found, overwrite it (select)",
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

	// Loop through to see if there are any matches.
	found := false
	for i := 0; i < len(envFile.Arr); i++ {
		// If the item is found, replace it.
		if envFile.Arr[i].Name == env.Name {
			found = true
			envFile.Arr[i] = env
			break
		}
	}

	// If the item is not found, add it to the end.
	if !found {
		envFile.Arr = append(envFile.Arr, env)
	}

	err = c.store.Save(envFile, c.Prefix, envFile.App)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Added:", env.Name)
}
