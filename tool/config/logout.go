package config

import (
	"fmt"
	"log"

	"github.com/demostack/cli/pkg/validate"

	"github.com/manifoldco/promptui"
)

// Logout .
func (c Config) Logout(f File, prefixList []string, passphrase *validate.Passphrase) {
	fmt.Println("Logout and remove local files.")

	pSelect := promptui.Select{
		Label: "Logout (select)",
		Items: []string{
			"yes",
			"no",
		},
	}
	answer := validate.MustSelect(pSelect.Run())
	if answer != "yes" {
		log.Fatalln("Logout cancelled.")
	}

	c.store.Delete(c.Prefix)

	for _, v := range prefixList {
		c.store.Delete(v)
	}

	fmt.Println("Logout successful.")
}
