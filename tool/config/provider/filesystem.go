package provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os/user"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/demostack/cli/tool"
)

// FilesystemProvider .
type FilesystemProvider struct {
	log tool.ILogger

	Base string
}

// NewFilesystemProvider .
func NewFilesystemProvider(l tool.ILogger) FilesystemProvider {
	return FilesystemProvider{
		log:  l,
		Base: ".demostack",
	}
}

// Filename returns the app configuration file. The strings are typically:
// prefix then app.
func (p FilesystemProvider) Filename(params ...string) string {
	f := ""

	if len(params) == 0 {
		f = fmt.Sprintf("%v.json", p.Base)
	} else {
		f = fmt.Sprintf("%v-%v.json", p.Base, strings.Join(params, "-"))
	}

	u, err := user.Current()
	if err != nil {
		fmt.Println("Cannot get the current user, using the current directory.")
		return f
	}

	return filepath.Join(u.HomeDir, f)
}

// LoadFile will load the configuration file for the app.
func (p FilesystemProvider) LoadFile(v interface{}, params ...string) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("value passed in must be a pointer")
	}

	filename := p.Filename(params...)

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		//fmt.Printf("Environment file for, %v, does not exist or cannot be read so a new one will be created.\n", f.App)
	} else {
		err = json.Unmarshal(b, v)
		if err != nil {
			return errors.New("unmarshal error: " + err.Error())
		}

		//fmt.Printf("Found %v secure environment variable(s).\n", len(f.Arr))
	}

	return nil
}

// Save .
func (p FilesystemProvider) Save(v interface{}, params ...string) error {
	filename := p.Filename(params...)

	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, b, 0644)
	if err != nil {

	} else {
		fmt.Printf("Saved to: %v\n", filename)
	}

	return err
}
