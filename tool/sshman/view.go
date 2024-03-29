package sshman

import (
	"fmt"
	"log"

	"github.com/demostack/cli/pkg/secure"
	"github.com/demostack/cli/pkg/validate"

	"github.com/manifoldco/promptui"
)

// View .
func (c Config) View(passphrase *validate.Passphrase) {
	fmt.Println("View SSH entry.")

	// Load the entries.
	sshFile := new(SSHFile)
	err := c.store.Load(sshFile, c.Prefix)
	if err != nil {
		log.Fatalln(err)
	}

	if len(sshFile.Arr) == 0 {
		fmt.Println("No SSH entries found.")
		return
	}

	arr := make([]string, 0)
	for i := 0; i < len(sshFile.Arr); i++ {
		arr = append(arr, sshFile.Arr[i].Name)
	}

	pSelect := promptui.Select{
		Label: "Choose the SSH entry (select)",
		Items: arr,
	}
	name := validate.MustSelect(pSelect.Run())

	var ent SSHEntry
	for _, v := range sshFile.Arr {
		if v.Name == name {
			ent = v
			break
		}
	}

	// Decrypt the private key.
	priKey, err := secure.Decrypt(ent.PrivateKey, passphrase.Password())
	if err != nil {
		log.Fatalln(err)
	}

	// Generate the private key that is used on all other steps.
	pri, err := secure.ParsePrivatePEM(priKey)
	if err != nil {
		log.Fatalln(err)
	}

	// Generate the public key.
	pubPem, err := secure.PublicKeyToPEM(&pri.PublicKey)
	if err != nil {
		log.Fatalln(err)
	}

	// Generate the authorized public key.
	pub, err := secure.PublicKeyToAuthorizedKey(&pri.PublicKey)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println()
	fmt.Println("Command:")
	fmt.Printf("ssh %v@%v\n", ent.User, ent.Hostname)
	fmt.Println()
	fmt.Println("SSH Config (~/.ssh/config):")
	fmt.Printf("Host %v\n", ent.Name)
	fmt.Printf("  Hostname %v\n", ent.Hostname)
	fmt.Printf("  User %v\n", ent.User)
	fmt.Println()
	fmt.Println("Private Key:")
	fmt.Println(secure.PrivateKeyToPEM(pri))
	fmt.Println("Public Key:")
	fmt.Println(pubPem)
	fmt.Println("Authorized Public Key:")
	fmt.Println(string(pub))
}
