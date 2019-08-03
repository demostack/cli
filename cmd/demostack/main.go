package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"

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
)

// init sets runtime settings.
func init() {
	// Verbose logging with file name and line number
	log.SetFlags(log.Lshortfile)
}

func main() {
	app.Version(Version)
	app.VersionFlag.Short('v')
	app.HelpFlag.Short('h')

	argList := os.Args[1:]
	arg := kingpin.MustParse(app.Parse(argList))

	switch arg {
	case cRun.FullCommand():
		Run(*cRunArgs)
	}
}

// Run a command and pass optional arguments. Supports passing in environment
// variables.
func Run(osArgs []string) {
	app := ""
	args := make([]string, 0)

	if len(osArgs) >= 1 {
		app = osArgs[0]
	}

	if len(osArgs) >= 2 {
		args = osArgs[1:]
	}

	cmd := exec.Command(app, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = os.Environ()

	err := cmd.Start()
	if err != nil {
		log.Fatalf("cmd.Start() failed with '%s'\n", err)
	}

	err = cmd.Wait()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				os.Exit(status.ExitStatus())
			}
			log.Fatalf("cmd.Wait error 1: %v\n", err)
		} else {
			log.Fatalf("cmd.Wait error 2: %v\n", err)
		}
	}
}
