package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os/user"
	"path/filepath"

	"github.com/demostack/cli/pkg/awslib"
	"github.com/demostack/cli/pkg/secure"
	"github.com/demostack/cli/pkg/validate"

	"github.com/manifoldco/promptui"
)

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

// Filename returns the demostack configuration file.
func Filename() string {
	f := fmt.Sprint(".demostack-config.json")

	u, err := user.Current()
	if err != nil {
		fmt.Println("Cannot get the current user, using the current directory.")
		return f
	}

	return filepath.Join(u.HomeDir, f)
}

// LoadFile will load the configuration file for the app.
func LoadFile() (File, error) {
	f := File{
		Storage: Storage{
			Current: "filesystem",
		},
	}

	filename := Filename()

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Initialization - Please set a password.")
		password := ""

		for true {
			// If no password is set, create a new one.
			prompt := promptui.Prompt{
				Label:    "New password (secure)",
				Default:  "",
				Mask:     '*',
				Validate: validate.EncryptionKey,
			}
			password = validate.Must(prompt.Run())

			// If no password is set, create a new one.
			prompt = promptui.Prompt{
				Label:    "Verify new password (password)",
				Default:  "",
				Mask:     '*',
				Validate: validate.EncryptionKey,
			}
			verifyPassword := validate.Must(prompt.Run())

			if password == verifyPassword {
				break
			}

			fmt.Println("Password don't match, please try again.")
		}

		id, err := secure.UUID()
		if err != nil {
			return f, errors.New("cannot generate UUID: " + err.Error())
		}

		f.ID, err = secure.Encrypt(id, password)
		if err != nil {
			return f, errors.New("cannot encrypt UUID: " + err.Error())
		}

		err = SaveFile(f)
		if err != nil {
			return f, errors.New("config save error: " + err.Error())
		}
		fmt.Printf("New config created: %v.\n", filename)
	} else {
		err = json.Unmarshal(b, &f)
		if err != nil {
			return f, errors.New("config load error: " + err.Error())
		}

		fmt.Printf("Current storage provider: %v.\n", f.Storage.Current)
	}

	return f, nil
}

// SaveFile .
func SaveFile(f File) error {
	b, err := json.Marshal(f)
	if err != nil {
		return err
	}

	filename := Filename()
	err = ioutil.WriteFile(filename, b, 0644)
	if err != nil {
		return err
	}

	return nil
}
