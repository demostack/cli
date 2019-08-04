package sshman

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/demostack/cli/pkg"
	"github.com/demostack/cli/pkg/secure"
	"github.com/demostack/cli/pkg/validate"

	"github.com/manifoldco/promptui"
)

// Login .
func (c Config) Login() {
	fmt.Println("Login helper for SSH.")

	// Load the entries.
	sshFile := new(SSHFile)
	err := c.store.LoadFile(sshFile, c.Prefix)
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

	// Get the password.
	_, priKey := validate.DecryptValue(ent.PrivateKey)

	// Generate the private key that is used on all other steps.
	pri, err := secure.ParsePrivatePEM(priKey)
	if err != nil {
		log.Fatalln(err)
	}

	// Generate the public key.
	pub, err := secure.PublicKeyToAuthorizedKey(&pri.PublicKey)
	if err != nil {
		log.Fatalln(err)
	}

	// Create the private key on disk.
	priTemp, err := ioutil.TempFile("", "")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = priTemp.WriteString(secure.PrivateKeyToPEM(pri))
	if err != nil {
		log.Fatalln(err)
	}

	// Create the public key on disk.
	pubTemp, err := ioutil.TempFile("", "")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = pubTemp.WriteString(string(pub))
	if err != nil {
		log.Fatalln(err)
	}

	// Add the private key to the SSH agent.
	pkg.Run([]string{"ssh-add", priTemp.Name()}, os.Environ())
	// Remove the private key from disk.
	os.Remove(priTemp.Name())

	// Let the user know SSH is available for a limited amount of time.
	fmt.Println("Initial login available for 15 seconds:")
	fmt.Printf("ssh %v@%v\n", ent.User, ent.Hostname)

	time.Sleep(15 * time.Second)

	// Remove identity from the SSH agent using the public key.
	pkg.Run([]string{"ssh-add", "-d", pubTemp.Name()}, os.Environ())
	// Remove hte public key from disk.
	os.Remove(pubTemp.Name())
}
