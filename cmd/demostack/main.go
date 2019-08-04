package main

import (
	"fmt"
	"log"
	"os"

	"github.com/demostack/cli/pkg"
	"github.com/demostack/cli/pkg/logger"
	"github.com/demostack/cli/pkg/validate"
	"github.com/demostack/cli/tool"
	"github.com/demostack/cli/tool/appenv"
	"github.com/demostack/cli/tool/config"
	"github.com/demostack/cli/tool/config/provider"
	"github.com/demostack/cli/tool/sshman"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	// Version of the app.
	Version = "1.0"
)

var (
	app = kingpin.New("demostack", "A command-line application to interact with your stacks.")

	cRun     = app.Command("run", "Run a command with secure environment variables.")
	cRunArgs = cRun.Arg("arguments", "Command and optional arguments to run.").Required().Strings()

	cEnv      = app.Command("env", "Manage secure environment variables.")
	cEnvSet   = cEnv.Command("set", "Add or update a secure environment variable.")
	cEnvUnset = cEnv.Command("unset", "Remove a secure environment variable.")
	cEnvView  = cEnv.Command("view", "View a secure environment variable.")

	cSSH      = app.Command("ssh", "Manage SSH session.")
	cSSHNew   = cSSH.Command("new", "Set an SSH session and generate a new private key.")
	cSSHSet   = cSSH.Command("set", "Set an SSH session with an existing private key.")
	cSSHLogin = cSSH.Command("login", "Set up a SSH session helper.")
	cSSHView  = cSSH.Command("view", "View a SSH entry.")

	cConfig                  = app.Command("config", "Manage settings for demostack application.")
	cConfigStorage           = cConfig.Command("storage", "Manage storage for the application.")
	cConfigStorageFilesystem = cConfigStorage.Command("filesystem", "Set the storage to the local filesystem.")
	cConfigStorageAWS        = cConfigStorage.Command("aws", "Set the storage to AWS.")
)

func init() {
	// Verbose logging with file name and line number
	log.SetFlags(log.Lshortfile)
}

func main() {
	// Create the logger.
	l := logger.New(log.New(os.Stderr, "", log.LstdFlags))

	app.Version(Version)
	app.VersionFlag.Short('v')
	app.HelpFlag.Short('h')
	arg := kingpin.MustParse(app.Parse(os.Args[1:]))

	// Load the filesystem storage provider since it's always used.
	fs := provider.NewFilesystemProvider(l)

	// Load the configuration file.
	appConfig := config.NewConfig(l, fs)
	c, err := appConfig.Load()
	if err != nil {
		log.Fatalln(err)
	}

	// Load the current storage provider for the tools.
	var sp tool.IStorage
	switch c.Storage.Current {
	case "filesystem":
		sp = fs
	case "aws":
		pass, _ := validate.DecryptValue(c.ID)
		dec, err := c.Storage.AWS.Decrypted(pass)
		if err != nil {
			log.Fatalln(err)
		}
		sp = provider.NewAWSProvider(l, dec)
	}

	// Load the tools.
	appenvConfig := appenv.NewConfig(l, sp)
	sshmanConfig := sshman.NewConfig(l, sp)

	switch arg {
	case cRun.FullCommand():
		vars := os.Environ()
		f := new(appenv.EnvFile)
		err := sp.LoadFile(f, "env", (*cRunArgs)[0])
		if err == nil {
			var arr []string
			if ok, v := f.HasEncryptedValues(); ok {
				// If a password already exists, verify it.
				pass, _ := validate.DecryptValue(v)
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
		appenvConfig.Set()
	case cEnvUnset.FullCommand():
		appenvConfig.Unset()
	case cEnvView.FullCommand():
		appenvConfig.View()
	case cSSHNew.FullCommand():
		sshmanConfig.New()
	case cSSHSet.FullCommand():
		sshmanConfig.Set()
	case cSSHLogin.FullCommand():
		sshmanConfig.Login()
	case cSSHView.FullCommand():
		sshmanConfig.View()
	case cConfigStorageAWS.FullCommand():
		pass, _ := validate.DecryptValue(c.ID)
		appConfig.SetStorageAWS(c, pass)
	case cConfigStorageFilesystem.FullCommand():
		appConfig.SetStorageFilesystem(c)
	}
}
