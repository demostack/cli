package secureenv

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/demostack/cli/pkg/validate"

	"github.com/manifoldco/promptui"
)

// Unset a secure environment variable.
func Unset() {
	fmt.Println("Unset (delete) a secure environment variable.")

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

	arr := make([]string, 0)
	for i := 0; i < len(envFile.Arr); i++ {
		arr = append(arr, envFile.Arr[i].Name)
	}

	pSelect := promptui.Select{
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

	newArr := make([]EnvVar, 0)

	// Loop through to see if there are any matches.
	for _, v := range envFile.Arr {
		// Only copy the item to the new array if it's not marked for deletion.
		if v.Name != name {
			newArr = append(newArr, v)
		}
	}

	// Set the new array.
	envFile.Arr = newArr

	b, err := json.Marshal(envFile)
	if err != nil {
		log.Fatalln(err)
	}

	filename := Filename(envFile.App)
	err = ioutil.WriteFile(filename, b, 0644)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Deleted:", name)
	fmt.Printf("Saved to: %v\n", filename)
}
