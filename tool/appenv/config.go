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
	Apps map[string]EnvApp `json:"apps"`
}

// Profiles .
func (f *EnvFile) Profiles(app string) map[string]EnvProfile {
	_, ok := f.Apps[app]
	if !ok {
		f.Apps = make(map[string]EnvApp)
		f.Apps[app] = EnvApp{
			Profiles: make(map[string]EnvProfile),
		}
	}

	return f.Apps[app].Profiles
}

// Profile .
func (f *EnvFile) Profile(app, profile string) EnvProfile {
	f.Profiles(app)

	_, ok := f.Apps[app].Profiles[profile]
	if !ok {
		f.Apps[app].Profiles[profile] = EnvProfile{
			Vars: make(map[string]EnvVar),
		}
	}

	return f.Apps[app].Profiles[profile]
}

// Vars .
func (f *EnvFile) Vars(app, profile string) map[string]EnvVar {
	f.Profile(app, profile)

	return f.Apps[app].Profiles[profile].Vars
}

// Var .
func (f *EnvFile) Var(app, profile, name string) (EnvVar, bool) {
	f.Vars(app, profile)

	v, ok := f.Apps[app].Profiles[profile].Vars[name]
	if !ok {
		return EnvVar{}, false
	}

	return v, true
}

// SetVar .
func (f *EnvFile) SetVar(app, profile, name string, value EnvVar) {
	f.Var(app, profile, name)
	f.Apps[app].Profiles[profile].Vars[name] = value
}

// Encrypted returns an encrypted object.
func (f EnvFile) Encrypted(passphrase *validate.Passphrase) EnvFile {
	for kApp, vApp := range f.Apps {
		for kProfile, vProfile := range vApp.Profiles {
			for k, v := range vProfile.Vars {
				if v.Encrypted {
					env, err := secure.Encrypt(v.Value, passphrase.Password())
					if err != nil {
						fmt.Println("Could not encrypt var:", k)
						v.Value = ""
					} else {
						v.Value = env
					}
				}
				f.Apps[kApp].Profiles[kProfile].Vars[k] = v
			}
		}
	}

	return f
}

// Decrypted returns a decrypted object.
func (f EnvFile) Decrypted(passphrase *validate.Passphrase) EnvFile {
	for kApp, vApp := range f.Apps {
		for kProfile, vProfile := range vApp.Profiles {
			for k, v := range vProfile.Vars {
				if v.Encrypted {
					env, err := secure.Decrypt(v.Value, passphrase.Password())
					if err != nil {
						fmt.Println("Could not decrypt var:", k)
						v.Value = ""
					} else {
						v.Value = env
					}
				}
				f.Apps[kApp].Profiles[kProfile].Vars[k] = v
			}
		}
	}

	return f
}

// EnvApp represents an application.
type EnvApp struct {
	Profiles map[string]EnvProfile `json:"profiles"`
}

// EnvProfile is a profile for an application.
type EnvProfile struct {
	Vars map[string]EnvVar `json:"vars"`
}

// Strings returns a string array of environment variables.
func (f EnvProfile) Strings(passphrase *validate.Passphrase) []string {
	arr := make([]string, 0)
	for i, v := range f.Vars {
		s := v.String(i, passphrase)
		if len(s) > 0 {
			arr = append(arr, s)
		}
	}
	return arr
}

// EnvVar represents an environment variable.
type EnvVar struct {
	Value     string `json:"value"`
	Encrypted bool   `json:"encrypted"`
}

// String returns the name and value in this format: name=value.
func (ev EnvVar) String(name string, passphrase *validate.Passphrase) string {
	if ev.Encrypted {
		v, err := secure.Decrypt(ev.Value, passphrase.Password())
		if err != nil {
			fmt.Println("Could not decrypt var:", name)
			return ""
		}
		return fmt.Sprintf("%v=%v", name, v)
	}
	return fmt.Sprintf("%v=%v", name, ev.Value)
}
