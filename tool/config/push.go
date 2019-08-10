package config

import (
	"fmt"
	"log"

	"github.com/demostack/cli/pkg/validate"
	"github.com/demostack/cli/tool"
)

// Push .
func (c Config) Push(f File, from string, to string, src tool.IStorage, dst tool.IStorage, passphrase *validate.Passphrase) {
	fmt.Printf("Copy data %v filesystem to %v", from, to)

	// Check AWS credentials.
	if (from == "aws" || to == "aws") && !f.Storage.AWS.Valid() {
		log.Fatalln("AWS credentials are expired. Please login again to renew them.")
	}

	for _, prefix := range []string{
		"ssh",
		"env",
		"email",
	} {
		var data map[string]interface{}

		err := src.Load(&data, prefix)
		if err != nil {
			fmt.Println("Skipping:", prefix)
			continue
		}

		err = dst.Save(data, prefix)
		if err != nil {
			//log.Fatalln(err)
			fmt.Printf("Could not upload: %v. %v\n", prefix, err)
			continue
		}
	}

	fmt.Println("Push complete.")
}
