package secureenv

import (
	"fmt"
	"log"

	"github.com/demostack/cli/pkg/validate"

	"github.com/manifoldco/promptui"
)

// View a secure environment variable.
func View() {
	fmt.Println("View a secure environment variable.")

	// App name.
	prompt := promptui.Prompt{
		Label:    "App name (string)",
		Default:  "",
		Validate: validate.RequireString,
	}
	app := validate.Must(prompt.Run())

	// Load the vars.
	envFile, err := LoadFile(app)
	if err != nil {
		log.Fatalln(err)
	}

	if len(envFile.Arr) == 0 {
		fmt.Println("No items found for app:", app)
		return
	}

	arr := make([]string, 0)
	arr = append(arr, "(All)")
	for i := 0; i < len(envFile.Arr); i++ {
		arr = append(arr, envFile.Arr[i].Name)
	}

	pSelect := promptui.Select{
		Label: "Choose the env var to view (select)",
		Items: arr,
	}
	name := validate.MustSelect(pSelect.Run())

	if name == "(All)" {
		var arr []string
		if ok, v := envFile.HasEncryptedValues(); ok {
			// If a password already exists, verify it.
			pass, _ := DecryptValue(v)
			arr = envFile.Strings(pass)
		} else {
			// Pass a blank password since it won't be used.
			arr = envFile.Strings("")
		}

		for _, v := range arr {
			fmt.Println(v)
		}
		return
	}

	// Find the item.
	for _, v := range envFile.Arr {
		if v.Name == name {
			if !v.Encrypted {
				fmt.Println(v.String(""))
				return
			}

			password, _ := DecryptValue(v.Value)
			fmt.Println("Password correct.")
			fmt.Println(v.String(password))
			return
		}
	}

	fmt.Println("Could not find value:", name)
}
