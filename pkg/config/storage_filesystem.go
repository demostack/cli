package config

import (
	"fmt"
	"log"
)

// SetStorageFilesystem will set the storage to the filesystem.
func SetStorageFilesystem(c File) {
	//fmt.Println("Set the storage provider to the local filesystem.")

	name := "filesystem"

	c.Storage.Current = name
	c.Storage.Filesystem = Filesystem{}

	err := SaveFile(c)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Set provider to:", name)
	fmt.Printf("Saved to: %v\n", Filename())
}
