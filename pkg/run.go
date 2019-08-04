package pkg

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

// Run a command and pass optional arguments. Supports passing in environment
// variables.
func Run(osArgs []string, envs []string) {
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
	cmd.Env = envs

	err := cmd.Start()
	if err != nil {
		log.Fatalln(err)
	}

	err = cmd.Wait()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				os.Exit(status.ExitStatus())
			}
		}
		log.Fatalln(err)
	}
}
