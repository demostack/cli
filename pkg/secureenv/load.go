package secureenv

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os/user"
	"path/filepath"
)

// Filename returns the app configuration file.
func Filename(app string) string {
	f := fmt.Sprintf(".demostack-env-%v.json", app)

	u, err := user.Current()
	if err != nil {
		fmt.Println("Cannot get the current user, using the current directory.")
		return f
	}

	return filepath.Join(u.HomeDir, f)
}

// LoadFile will load the configuration file for the app.
func LoadFile(app string) (EnvFile, error) {
	f := EnvFile{
		App: app,
	}

	filename := Filename(f.App)

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		//fmt.Printf("Environment file for, %v, does not exist or cannot be read so a new one will be created.\n", f.App)
	} else {
		err = json.Unmarshal(b, &f)
		if err != nil {
			return f, errors.New("unmarshal error: " + err.Error())
		}

		fmt.Printf("Found %v secure environment variable(s).\n", len(f.Arr))
	}

	return f, nil
}
