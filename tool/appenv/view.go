package appenv

import (
	"fmt"
	"log"

	"github.com/demostack/cli/pkg/validate"

	"github.com/manifoldco/promptui"
)

// View a secure environment variable.
func (c Config) View(passphrase *validate.Passphrase) {
	fmt.Println("View a secure environment variable.")

	// Load the vars.
	envFile := new(EnvFile)
	err := c.store.Load(envFile, c.Prefix)
	if err != nil {
		log.Fatalln(err)
	}

	arr := make([]string, 0)
	for k := range envFile.Apps {
		arr = append(arr, k)
	}

	pSelect := promptui.Select{
		Label: "Choose the app to view (select)",
		Items: arr,
	}
	appName := validate.MustSelect(pSelect.Run())

	arr = make([]string, 0)
	for k := range envFile.Profiles(appName) {
		arr = append(arr, k)
	}

	pSelect = promptui.Select{
		Label: "Choose the profile to view (select)",
		Items: arr,
	}
	profileName := validate.MustSelect(pSelect.Run())

	vars := envFile.Vars(appName, profileName)
	if len(vars) == 0 {
		fmt.Printf("No items found for app (profile): %v (%v)\n", appName, profileName)
		return
	}

	arr = make([]string, 0)
	arr = append(arr, "(All)")
	for k := range vars {
		arr = append(arr, k)
	}

	pSelect = promptui.Select{
		Label: "Choose the env var to view (select)",
		Items: arr,
	}
	name := validate.MustSelect(pSelect.Run())

	if name == "(All)" {
		arr := envFile.Profile(appName, profileName).Strings(passphrase)
		for _, v := range arr {
			fmt.Println(v)
		}
		return
	}

	v, ok := envFile.Var(appName, profileName, name)
	if !ok {
		fmt.Println("Could not find value:", name)
		return
	}

	fmt.Println(v.String(name, passphrase))
}
