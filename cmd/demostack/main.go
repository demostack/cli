package main

import (
	"fmt"
	"log"
	"os"

	"github.com/demostack/cli/pkg"
	"github.com/demostack/cli/pkg/secureenv"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	// Version of the app.
	Version = "1.0"
)

var (
	app = kingpin.New("demostack", "A command-line application to interact with your stacks.")

	cRun     = app.Command("run", "Run a command with environment variables from encrypted storage.")
	cRunArgs = cRun.Arg("arguments", "Command and optional arguments to run.").Required().Strings()

	cEnv      = app.Command("env", "Manage secure environment variables.")
	cEnvSet   = cEnv.Command("set", "Add or update a secure environment variable.")
	cEnvUnset = cEnv.Command("unset", "Remove a secure environment variable.")
	cEnvView  = cEnv.Command("view", "View a secure environment variable.")
)

func init() {
	// Verbose logging with file name and line number
	log.SetFlags(log.Lshortfile)
}

func main() {
	app.Version(Version)
	app.VersionFlag.Short('v')
	app.HelpFlag.Short('h')
	arg := kingpin.MustParse(app.Parse(os.Args[1:]))

	switch arg {
	case cRun.FullCommand():
		vars := os.Environ()
		f, err := secureenv.LoadFile((*cRunArgs)[0])
		if err == nil {
			var arr []string
			if ok, v := f.HasEncryptedValues(); ok {
				// If a password already exists, verify it.
				pass, _ := secureenv.DecryptValue(v)
				arr = f.Strings(pass)
			} else {
				// Pass a blank password since it won't be used.
				arr = f.Strings("")
			}
			vars = append(vars, arr...)

			fmt.Printf("Loaded %v secure environment variable(s).\n", len(arr))
		} else {
			fmt.Printf("Found 0 secure environment variables.\n")
		}

		pkg.Run(*cRunArgs, vars)
	case cEnvSet.FullCommand():
		secureenv.Set()
	case cEnvUnset.FullCommand():
		secureenv.Unset()
	case cEnvView.FullCommand():
		secureenv.View()
	}
}
