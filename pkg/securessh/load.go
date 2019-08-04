package securessh

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os/user"
	"path/filepath"
)

// Filename returns the app configuration file.
func Filename() string {
	f := fmt.Sprintf(".demostack-ssh.json")

	u, err := user.Current()
	if err != nil {
		fmt.Println("Cannot get the current user, using the current directory.")
		return f
	}

	return filepath.Join(u.HomeDir, f)
}

// LoadFile will load the configuration file for all SSH.
func LoadFile() (SSHFile, error) {
	f := SSHFile{}

	filename := Filename()

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		//fmt.Printf("SSH file does not exist or cannot be read so a new one will be created.\n")
	} else {
		err = json.Unmarshal(b, &f)
		if err != nil {
			return f, errors.New("unmarshal error: " + err.Error())
		}

		if len(f.Arr) > 1 {
			fmt.Printf("Found %v SSH entries.\n", len(f.Arr))
		} else {
			fmt.Printf("Found %v SSH entry.\n", len(f.Arr))
		}
	}

	return f, nil
}
