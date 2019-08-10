package config

import (
	"fmt"
	"log"
)

// SetStorageAWS will set the AWS storage.
func (c Config) SetStorageAWS(f File) {
	//fmt.Println("Set the storage provider to AWS.")

	if len(f.Storage.AWS.AccessKeyID) == 0 {
		log.Fatalln("AWS credentials don't exist, please login first.")
	}

	name := "aws"

	f.Storage.Current = name

	err := c.store.Save(f, c.Prefix)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Set provider to:", name)
}
