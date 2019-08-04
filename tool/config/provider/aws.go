package provider

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/demostack/cli/pkg/awslib"

	"github.com/demostack/cli/tool"
)

// AWSProvider .
type AWSProvider struct {
	log   tool.ILogger
	creds awslib.Storage

	Base string
}

// NewAWSProvider .
func NewAWSProvider(l tool.ILogger, creds awslib.Storage) AWSProvider {
	return AWSProvider{
		log:   l,
		Base:  "demostack",
		creds: creds,
	}
}

// Filename returns the app configuration file. The strings are typically:
// prefix then app.
func (p AWSProvider) Filename(params ...string) string {
	f := ""

	if len(params) == 0 {
		f = fmt.Sprintf("%v.json", p.Base)
	} else {
		f = fmt.Sprintf("%v-%v.json", p.Base, strings.Join(params, "-"))
	}

	return f
}

// LoadFile will load the configuration file for the app.
func (p AWSProvider) LoadFile(v interface{}, params ...string) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("value passed in must be a pointer")
	}

	filename := p.Filename(params...)

	b, err := awslib.Download(p.creds, p.creds.Bucket, filename)
	if err != nil {

		//return err
	} else {
		err = json.Unmarshal(b, v)
		if err != nil {
			return errors.New("unmarshal error: " + err.Error())
		}
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

	filename := p.Filename(params...)

	err = awslib.Upload(p.creds, p.creds.Bucket, filename, bytes.NewBuffer(b))
	if err != nil {

	} else {
		fmt.Printf("Saved to: %v\n", filename)
	}
	return err
}
