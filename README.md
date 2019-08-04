# demostack - CLI

This command-line application helps you manage sensitive environment variables and SSH information. Data is stored in JSON configuration files with sensitive information encrypted using a password.

Features:

- store environment variables so they are encrypted and then pass them into applications at runtime. For example, instead of storing your AWS credentials in plaintext, you can store them securely in demostack and then execute your AWS commands like this: `demostack run aws s3 ls` or `demostack run aws sts get-caller-identity`.
- store your SSH keys with a password and then add them to your ssh-agent via temp files for 15 seconds while you login and then removes from the ssh-agent and the temp files. It also makes it very easy to generate your public key or authorized_key public key from your private key.

Keep in mind that that environment variables are still visible for all running applications on a Mac using `ps eww <PID>` or on Linux using `ps faux | grep 'PROCESSHERE'` and then `cat /proc/PIDHERE/environ`.

Here is the syntax for the application:

```
usage: demostack [<flags>] <command> [<args> ...]

A command-line application to manage sensitive environment variables and SSH securely.

Flags:
  -h, --help     Show context-sensitive help (also try --help-long and --help-man).
  -v, --version  Show application version.

Commands:
  help [<command>...]
    Show help.

  run <arguments>...
    Run a command with secure environment variables.

  env set
    Add or update a secure environment variable.

  env unset
    Remove a secure environment variable.

  env view
    View a secure environment variable.

  ssh new
    Set an SSH session and generate a new private key.

  ssh set
    Set an SSH session with an existing private key.

  ssh login
    Set up a SSH session helper.

  ssh view
    View a SSH entry.

  config storage filesystem
    Set the storage to the local filesystem.

  config storage aws
    Set the storage to AWS.
```