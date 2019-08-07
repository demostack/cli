# demostack - CLI

[![Go Report Card](https://goreportcard.com/badge/github.com/demostack/cli)](https://goreportcard.com/report/github.com/demostack/cli)

This command-line application helps you manage sensitive environment variables and SSH information. Data is stored in JSON configuration files with sensitive information encrypted using a password. Also has other useful utilities like built in SMTP functionality.

Features:

- store environment variables so they are encrypted and then pass them into applications at runtime. For example, instead of storing your AWS credentials in plaintext, you can store them securely in demostack and then execute your AWS commands like this: `demostack run default aws s3 ls` or `demostack run aws default sts get-caller-identity`.
- store your SSH keys with a password and then add them to your ssh-agent via temp files for 15 seconds while you login and then removes from the ssh-agent and the temp files. It also makes it very easy to generate your public key or authorized_key public key from your private key.
- set SMTP credentials and send emails with an optional file attachment.
- encrypt and decrypt single values with a password.
- change the password used to encrypt all sensitive items.

Keep in mind that that environment variables are still visible for all running applications on a Mac using `ps eww <PID>` or on Linux using `ps faux | grep 'PROCESSHERE'` and then `cat /proc/PIDHERE/environ`.

Here is the syntax for the application:

```
usage: demostack [<flags>] <command> [<args> ...]

A command-line application utility for managing encrypted environment
variables and SSH connection information. Also contains SMTP functionality.

Flags:
  -h, --help     Show context-sensitive help (also try --help-long and --help-man).
  -v, --version  Show application version.

Commands:
  help [<command>...]
    Show help.

  run <arguments>...
    Run a command with secure environment variables.

  encrypt
    Encrypt a variable with a password.

  decrypt
    Decrypt a variable with a password.

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

  ssh unset
    Remove an SSH session and private key.

  ssh login
    Set up a SSH session helper.

  ssh view
    View a SSH entry.

  config view
    View the settings for the application.

  config change-password
    Change the config password.

  config smtp
    Set settings for the SMTP server for the application.

  config storage filesystem
    Set the storage to the local filesystem.

  config storage aws
    Set the storage to AWS.

  email
    Send an email via SMTP.
```