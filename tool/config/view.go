package config

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/demostack/cli/pkg/validate"

	"github.com/manifoldco/promptui"
)

// View will output the application settings..
func (c Config) View(f File, passphrase *validate.Passphrase) {
	fmt.Println("Current application settings.")

	pSelect := promptui.Select{
		Label: "Would you like to view with variables decrypted (select)",
		Items: []string{
			"no",
			"yes",
		},
	}
	answer := validate.MustSelect(pSelect.Run())

	fmt.Println("Location:", c.store.Filename(c.Prefix))

	if answer == "yes" {
		b, err := json.MarshalIndent(f.Decrypted(passphrase), "", "    ")
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(string(b))
		return
	}

	b, err := json.MarshalIndent(f, "", "    ")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(b))
}
