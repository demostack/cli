package sendmail

import (
	"crypto/tls"
	"os"

	"github.com/demostack/cli/pkg/secure"
	"github.com/demostack/cli/pkg/validate"

	gomail "gopkg.in/gomail.v2"
)

// SMTP is the set of SMTP credentials from the environment variables
type SMTP struct {
	Host       string `json:"host"`
	Port       int    `json:"port"`
	Password   string `json:"password"`
	Username   string `json:"username"`
	From       string `json:"from"`
	SkipVerify bool   `json:"skip_verify"`
}

// Encrypted returns an encrypted object.
func (f SMTP) Encrypted(passphrase *validate.Passphrase) (SMTP, error) {
	var err error
	f.Password, err = secure.Encrypt(f.Password, passphrase.Password())
	if err != nil {
		return f, err
	}
	return f, nil
}

// Decrypted returns a decrypted object.
func (f SMTP) Decrypted(passphrase *validate.Passphrase) (SMTP, error) {
	if f.Host == "" {
		return f, nil
	}

	var err error
	f.Password, err = secure.Decrypt(f.Password, passphrase.Password())
	if err != nil {
		return f, err
	}
	return f, nil
}

// SendMail will send an email.
func (f SMTP) SendMail(toAddresses []string, emailSubject, emailBody string,
	skipVerify bool, attachments ...string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", f.From)
	m.SetHeader("To", toAddresses...)
	m.SetHeader("Subject", emailSubject)
	m.SetBody("text/html", emailBody)

	for _, mail := range attachments {
		if len(mail) == 0 {
			continue
		}

		if _, err := os.Stat(mail); err == nil {
			m.Attach(mail)
		}
	}

	d := gomail.NewPlainDialer(f.Host, f.Port, f.Username, f.Password)
	if skipVerify {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	return d.DialAndSend(m)
}
