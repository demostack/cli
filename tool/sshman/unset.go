package sshman

import (
	"fmt"
	"log"

	"github.com/demostack/cli/pkg/validate"

	"github.com/manifoldco/promptui"
)

// Unset .
func (c Config) Unset(passphrase *validate.Passphrase) {
	fmt.Println("Delete an SSH entry.")

	// Load the entries.
	sshFile := new(SSHFile)
	err := c.store.Load(sshFile, c.Prefix)
	if err != nil {
		log.Fatalln(err)
	}

	if len(sshFile.Arr) == 0 {
		fmt.Println("No SSH entries found.")
		return
	}

	arr := make([]string, 0)
	for i := 0; i < len(sshFile.Arr); i++ {
		arr = append(arr, sshFile.Arr[i].Name)
	}

	pSelect := promptui.Select{
		Label: "Choose the SSH entry to delete (select)",
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

	newArr := make([]SSHEntry, 0)
	for _, v := range sshFile.Arr {
		if v.Name != name {
			newArr = append(newArr, v)
		}
	}

	sshFile.Arr = newArr

	err = c.store.Save(sshFile, c.Prefix)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Deleted:", name)
}
