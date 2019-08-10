package main

import (
	"fmt"
	"log"
	"os"

	"github.com/demostack/cli/pkg"
	"github.com/demostack/cli/pkg/logger"
	"github.com/demostack/cli/pkg/secure"
	"github.com/demostack/cli/pkg/validate"
	"github.com/demostack/cli/tool"
	"github.com/demostack/cli/tool/appenv"
	"github.com/demostack/cli/tool/config"
	"github.com/demostack/cli/tool/config/provider"
	"github.com/demostack/cli/tool/email"
	"github.com/demostack/cli/tool/sshman"

	"github.com/manifoldco/promptui"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	// Version of the app.
	Version = "1.0"
)

var (
	app = kingpin.New("demostack", "A command-line application utility for managing encrypted environment variables and SSH connection information. Also contains SMTP functionality.")

	cRun     = app.Command("run", "Run a command with secure environment variables.")
	cRunArgs = cRun.Arg("arguments", "Command and optional arguments to run.").Required().Strings()

	cLogin     = app.Command("login", "Login to the authentication service.")
	cLoginTest = app.Command("login-test", "Test the AWS credentials from authentication service.")
	cLogout    = app.Command("logout", "Logout and delete local files.")

	cEncrypt = app.Command("encrypt", "Encrypt a variable with a password.")
	cDecrypt = app.Command("decrypt", "Decrypt a variable with a password.")

	cEnv      = app.Command("env", "Manage secure environment variables.")
	cEnvSet   = cEnv.Command("set", "Add or update a secure environment variable.")
	cEnvUnset = cEnv.Command("unset", "Remove a secure environment variable.")
	cEnvView  = cEnv.Command("view", "View a secure environment variable.")

	cSSH      = app.Command("ssh", "Manage SSH session.")
	cSSHNew   = cSSH.Command("new", "Set an SSH session and generate a new private key.")
	cSSHSet   = cSSH.Command("set", "Set an SSH session with an existing private key.")
	cSSHUnset = cSSH.Command("unset", "Remove an SSH session and private key.")
	cSSHLogin = cSSH.Command("login", "Set up a SSH session helper.")
	cSSHView  = cSSH.Command("view", "View a SSH entry.")

	cConfig               = app.Command("config", "Manage settings for demostack application.")
	cConfigView           = cConfig.Command("view", "View the settings for the application.")
	cConfigChangePassword = cConfig.Command("change-password", "Change the config password.")

	cEmail     = app.Command("email", "Manage email settings an email via SMTP.")
	cEmailSet  = cEmail.Command("set", "Set settings for the SMTP server for the application.")
	cEmailSend = cEmail.Command("send", "Set settings for the SMTP server for the application.")

	cSMS = app.Command("sms", "Send an SMS text message via AWS SNS.")

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

	// Set up the application.
	app.Version(Version)
	app.VersionFlag.Short('v')
	app.HelpFlag.Short('h')
	arg := kingpin.MustParse(app.Parse(os.Args[1:]))

	// Load the filesystem storage provider since it's always used.
	fs := provider.NewFilesystemProvider(l)

	// Load the configuration file.
	appConfig := config.NewConfig(l, fs)
	c, passphrase, err := appConfig.Load()
	if err != nil {
		log.Fatalln(err)
	}

	// Set the passphrase if not already set.
	if passphrase == nil {
		// Create the passphrase object. It will only prompt the user for the
		// password when it's required.
		passphrase = validate.NewPassphrase(c.ID)
	}

	// Load the current storage provider for the tools.
	var sp tool.IStorage
	switch c.Storage.Current {
	case "filesystem":
		sp = fs
	case "aws":
		sp = provider.NewAWSProvider(l, c.Storage.AWS, passphrase)
	}

	// Load the tools.
	appenvConfig := appenv.NewConfig(l, sp)
	sshmanConfig := sshman.NewConfig(l, sp)
	emailConfig := email.NewConfig(l, sp)

	switch arg {
	case cLogin.FullCommand():
		appConfig.Login(c, passphrase)
	case cLoginTest.FullCommand():
		appConfig.LoginTest(c, passphrase)
	case cLogout.FullCommand():
		appConfig.Logout(c, passphrase)

	case cRun.FullCommand():
		if len(*cRunArgs) < 2 {
			log.Fatalln("Command requires a profile and then the app command.")
		}
		vars := os.Environ()
		f := new(appenv.EnvFile)
		err := sp.Load(f, "env")
		if err == nil {
			profile, ok := f.Apps[(*cRunArgs)[1]].Profiles[(*cRunArgs)[0]]
			if ok {
				arr := profile.Strings(passphrase)
				vars = append(vars, arr...)
				fmt.Printf("Loaded %v secure environment variable(s).\n", len(arr))
			} else {
				fmt.Printf("Found 0 secure environment variables.\n")
			}
		} else {
			fmt.Printf("Found 0 secure environment variables.\n")
		}

		pkg.Run((*cRunArgs)[1:], vars)
	case cEncrypt.FullCommand():
		fmt.Println("Encrypt a value with a password.")

		prompt := promptui.Prompt{
			Label:    "Value (secure)",
			Default:  "",
			Mask:     '*',
			Validate: validate.RequireString,
		}
		value := validate.Must(prompt.Run())

		prompt = promptui.Prompt{
			Label:    "Password (secure)",
			Default:  "",
			Mask:     '*',
			Validate: validate.RequireString,
		}
		password := validate.Must(prompt.Run())

		enc, err := secure.Encrypt(value, password)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println("Encrypted Value:", enc)
	case cDecrypt.FullCommand():
		fmt.Println("Decrypt a value with a password.")

		prompt := promptui.Prompt{
			Label:    "Value (string)",
			Default:  "",
			Validate: validate.RequireString,
		}
		value := validate.Must(prompt.Run())

		enc := ""
		for {
			prompt = promptui.Prompt{
				Label:    "Password (secure)",
				Default:  "",
				Mask:     '*',
				Validate: validate.RequireString,
			}
			password := validate.Must(prompt.Run())

			enc, err = secure.Decrypt(value, password)
			if err != nil {
				fmt.Println("wrong password")
			} else {
				break
			}
		}

		fmt.Println("Decrypted Value:", enc)

	case cEnvSet.FullCommand():
		appenvConfig.Set(passphrase)
	case cEnvUnset.FullCommand():
		appenvConfig.Unset()
	case cEnvView.FullCommand():
		appenvConfig.View(passphrase)

	case cSSHNew.FullCommand():
		sshmanConfig.New(passphrase)
	case cSSHSet.FullCommand():
		sshmanConfig.Set(passphrase)
	case cSSHUnset.FullCommand():
		sshmanConfig.Unset(passphrase)
	case cSSHLogin.FullCommand():
		sshmanConfig.Login(passphrase)
	case cSSHView.FullCommand():
		sshmanConfig.View(passphrase)

	case cConfigView.FullCommand():
		appConfig.View(c, passphrase)

	case cConfigChangePassword.FullCommand():
		appConfig.ChangePassword(c, passphrase)

	case cEmailSet.FullCommand():
		emailConfig.SetSMTP(passphrase)

	case cEmailSend.FullCommand():
		emailConfig.SendSMTP(passphrase)

	case cSMS.FullCommand():
		appConfig.SendSMS(c, passphrase)

	case cConfigStorageAWS.FullCommand():
		appConfig.SetStorageAWS(c, passphrase.Password())
	case cConfigStorageFilesystem.FullCommand():
		appConfig.SetStorageFilesystem(c)
	}
}
