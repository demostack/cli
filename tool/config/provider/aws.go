package provider

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/demostack/cli/pkg/awslib"
	"github.com/demostack/cli/pkg/validate"
	"github.com/demostack/cli/tool"
)

// AWSProvider .
type AWSProvider struct {
	log        tool.ILogger
	creds      awslib.Storage
	passphrase *validate.Passphrase

	Base string
}

// NewAWSProvider .
func NewAWSProvider(l tool.ILogger, creds awslib.Storage, passphrase *validate.Passphrase) AWSProvider {
	return AWSProvider{
		log:        l,
		creds:      creds,
		passphrase: passphrase,

		Base: "demostack",
	}
}

// Key returns the app configuration file path. The strings are typically:
// prefix then app.
func (p AWSProvider) Key(params ...string) string {
	f := ""

	if len(params) == 0 {
		f = fmt.Sprintf("%v.json", p.Base)
	} else {
		f = fmt.Sprintf("%v-%v.json", p.Base, strings.Join(params, "-"))
	}

	return f
}

// Filename returns the path to the file in the S3 bucket.
func (p AWSProvider) Filename(params ...string) string {
	return fmt.Sprintf("s3://%v/%v", p.creds.Bucket, p.Key(params...))
}

// Load will load the configuration file for the app.
func (p AWSProvider) Load(v interface{}, params ...string) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("value passed in must be a pointer")
	}

	filename := p.Key(params...)

	dec, err := p.creds.Decrypted(p.passphrase.Password())
	if err != nil {
		return err
	}

	b, err := awslib.Download(dec, dec.Bucket, filename)
	if err != nil {
		return ErrFileNotFound
	}

	err = json.Unmarshal(b, v)
	if err != nil {
		fmt.Println("unmarshal error: ", err)
		return ErrFileUnmarshalError
	}

	//fmt.Printf("Found %v secure environment variable(s).\n", len(f.Arr))

	return nil
}

// Save .
func (p AWSProvider) Save(v interface{}, params ...string) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	filename := p.Key(params...)

	dec, err := p.creds.Decrypted(p.passphrase.Password())
	if err != nil {
		return err
	}

	err = awslib.Upload(dec, dec.Bucket, filename, bytes.NewBuffer(b))
	if err != nil {

	} else {
		fmt.Printf("Saved to: %v\n", p.Filename(params...))
	}
	return err
}

// Delete .
func (p AWSProvider) Delete(params ...string) error {
	filename := p.Key(params...)

	dec, err := p.creds.Decrypted(p.passphrase.Password())
	if err != nil {
		return err
	}

	err = awslib.DeleteObject(dec, dec.Bucket, filename)
	if err != nil {

	} else {
		fmt.Printf("Deleted: %v\n", p.Filename(params...))
	}
	return err
}
