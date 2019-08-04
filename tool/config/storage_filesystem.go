package config

import (
	"fmt"
	"log"
)

// SetStorageFilesystem will set the storage to the filesystem.
func (c Config) SetStorageFilesystem(f File) {
	//fmt.Println("Set the storage provider to the local filesystem.")

	name := "filesystem"

	f.Storage.Current = name
	f.Storage.Filesystem = Filesystem{}

	err := c.store.Save(f, c.Prefix)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Set provider to:", name)
}
