package appenv

import (
	"fmt"

	"github.com/demostack/cli/pkg/secure"
	"github.com/demostack/cli/pkg/validate"
	"github.com/demostack/cli/tool"
)

// Config .
type Config struct {
	log   tool.ILogger
	store tool.IStorage

	Prefix string
}

// NewConfig .
func NewConfig(l tool.ILogger, store tool.IStorage) Config {
	return Config{
		log:    l,
		store:  store,
		Prefix: "env",
	}
}

// EnvFile is an environment config file.
type EnvFile struct {
	App string   `json:"app"`
	Arr []EnvVar `json:"vars"`
}

// Strings returns a string array of environment variables.
func (ef EnvFile) Strings(passphrase *validate.Passphrase) []string {
	arr := make([]string, 0)
	for _, v := range ef.Arr {
		s := v.String(passphrase)
		if len(s) > 0 {
			arr = append(arr, s)
		}
	}
	return arr
}

// EnvVar represents an environment variable.
type EnvVar struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	Encrypted bool   `json:"encrypted"`
}

// String returns the name and value in this format: name=value.
func (ev EnvVar) String(passphrase *validate.Passphrase) string {
	if ev.Encrypted {
		v, err := secure.Decrypt(ev.Value, passphrase.Password())
		if err != nil {
			fmt.Println("Could not decrypt var:", ev.Name)
			return ""
		}
		return fmt.Sprintf("%v=%v", ev.Name, v)
	}
	return fmt.Sprintf("%v=%v", ev.Name, ev.Value)
}
