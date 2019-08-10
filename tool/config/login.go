package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/demostack/cli/pkg/awslib"
	"github.com/demostack/cli/pkg/validate"

	"github.com/manifoldco/promptui"
)

// Login .
func (c Config) Login(f File, passphrase *validate.Passphrase) {
	fmt.Println("Login to authentication service")

	URLDefault := ""
	UsernameDefault := ""
	if len(f.Login.Host) > 0 {
		URLDefault = f.Login.Host
		UsernameDefault = f.Login.Username
	}

	prompt := promptui.Prompt{
		Label:     "URL (string)",
		Default:   URLDefault,
		AllowEdit: true,
		Validate:  validate.RequireString,
	}
	URL := validate.Must(prompt.Run())

	prompt = promptui.Prompt{
		Label:    "Username (string)",
		Default:  UsernameDefault,
		Validate: validate.RequireString,
	}
	username := validate.Must(prompt.Run())

	prompt = promptui.Prompt{
		Label:    "Password (secure)",
		Default:  "",
		Mask:     '*',
		Validate: validate.RequireString,
	}
	password := validate.Must(prompt.Run())

	prompt = promptui.Prompt{
		Label:    "Duration in minutes (int, 15-2160)",
		Default:  "15",
		Validate: validate.RequireAWSSessionInt,
	}
	duration := validate.Must(prompt.Run())

	dur, _ := strconv.Atoi(duration)

	type Request struct {
		Username        string `json:"username"`
		Password        string `json:"password"`
		DurationSeconds int64  `json:"duration_seconds"`
	}
	r := Request{
		Username:        username,
		Password:        password,
		DurationSeconds: int64(dur * 60),
	}
	b, err := json.Marshal(r)
	if err != nil {
		log.Fatalln(err)
	}

	req, err := http.NewRequest("POST", URL+"/auth", bytes.NewReader(b))
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Add("Content-Type", "application/json")

	fmt.Println("Sending MFA request...")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	if resp.StatusCode != 200 {
		b, _ := ioutil.ReadAll(resp.Body)
		log.Fatalf("error: %v", string(b))
	}

	type Response struct {
		Status          int    `json:"status"`
		AccessKeyID     string `json:"access_key_id"`
		SecretAccessKey string `json:"secret_access_key"`
		SessionToken    string `json:"session_token"`
	}

	result := Response{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Login successful.")

	expiration := time.Now()
	key := awslib.Storage{
		AccessKeyID:     result.AccessKeyID,
		SecretAccessKey: result.SecretAccessKey,
		SessionToken:    result.SessionToken,
		Expiration:      expiration.Add(time.Second * time.Duration(r.DurationSeconds)),
	}

	prompt = promptui.Prompt{
		Label:    "AWS Region (string)",
		Default:  "us-east-1",
		Validate: validate.RequireString,
	}
	key.Region = validate.Must(prompt.Run())

	account, err := awslib.AccountNumber(key)
	if err != nil {
		fmt.Println(err)
	}

	key.Bucket = fmt.Sprintf("%v-demostack-config", account)

	prompt = promptui.Prompt{
		Label:    "S3 Bucket to store config (string)",
		Default:  key.Bucket,
		Validate: validate.RequireString,
	}
	key.Bucket = validate.Must(prompt.Run())

	// Create the S3 bucket. Will not fail if bucket already exists.
	err = awslib.CreateBucket(key)
	if err != nil {
		log.Fatalln(err)
	}

	// Encrypt the sensitive information.
	key, err = key.Encrypted(passphrase.Password())
	if err != nil {
		log.Fatalln(err)
	}

	f.Storage.Current = "aws"
	f.Storage.AWS = key
	f.Login = Login{
		Host:     URL,
		Username: username,
	}

	err = c.store.Save(f, c.Prefix)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Login done.")
}
