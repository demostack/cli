package secure

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"

	"golang.org/x/crypto/ssh"
)

// GenerateRSAKeyPair will generate a private and public key.
func GenerateRSAKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {
	pri, _ := rsa.GenerateKey(rand.Reader, 4096)
	return pri, &pri.PublicKey
}

// PrivateKeyToPEM converts private key to a PEM string.
func PrivateKeyToPEM(privkey *rsa.PrivateKey) string {
	priBytes := x509.MarshalPKCS1PrivateKey(privkey)
	return string(pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: priBytes,
		},
	))
}

// ParsePrivatePEM converts a string PEM into a private key.
func ParsePrivatePEM(privPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}

// PublicKeyToPEM converts a public key to a string PEM.
func PublicKeyToPEM(pubkey *rsa.PublicKey) (string, error) {
	pubBytes, err := x509.MarshalPKIXPublicKey(pubkey)
	if err != nil {
		return "", err
	}

	return string(pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubBytes,
		},
	)), nil
}

// PublicKeyToAuthorizedKey converts a public key to an authorized key for use
// with the authorized_keys file.
func PublicKeyToAuthorizedKey(pub *rsa.PublicKey) ([]byte, error) {
	sshPubKey, err := ssh.NewPublicKey(pub)
	if err != nil {
		return nil, err
	}

	return ssh.MarshalAuthorizedKey(sshPubKey), nil
}
