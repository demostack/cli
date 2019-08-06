package appenv

import (
	"fmt"
	"log"

	"github.com/demostack/cli/pkg/validate"

	"github.com/manifoldco/promptui"
)

// Unset a secure environment variable.
func (c Config) Unset() {
	fmt.Println("Unset (delete) a secure environment variable.")

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
	for k := range vars {
		arr = append(arr, k)
	}

	pSelect = promptui.Select{
		Label: "Choose the env var to delete (select)",
		Items: arr,
	}
	name := validate.MustSelect(pSelect.Run())

	pSelect = promptui.Select{
		Label: "Confirm delete of: " + name + " (select)",
		Items: []string{
			"no",
			"yes",
		},
	}
	confirm := validate.MustSelect(pSelect.Run())

	if confirm != "yes" {
		fmt.Println("Delete cancelled.")
		return
	}

	delete(envFile.Apps[appName].Profiles[profileName].Vars, name)

	err = c.store.Save(envFile, c.Prefix)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Deleted:", name)
}
