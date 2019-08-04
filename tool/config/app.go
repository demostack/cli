package config

import "github.com/demostack/cli/pkg/awslib"

// File is the demostack config file.
type File struct {
	ID      string  `json:"id"`
	Storage Storage `json:"storage"`
}

// Storage is the storage of the config files.
type Storage struct {
	// Current supports the following values: filesystem, aws.
	Current    string         `json:"current"`
	AWS        awslib.Storage `json:"aws"`
	Filesystem Filesystem     `json:"filesystem"`
}

// Filesystem is for the local filesystem.
type Filesystem struct{}
