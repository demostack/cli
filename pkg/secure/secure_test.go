package secure_test

import (
	"testing"

	"github.com/demostack/cli/pkg/secure"

	"github.com/stretchr/testify/assert"
)

func TestSuccess(t *testing.T) {
	text := "This is a test"
	pass := "password"

	enc, err := secure.Encrypt(text, pass)
	assert.Nil(t, err)

	dec, err := secure.Decrypt(enc, pass)
	assert.Nil(t, err)
	assert.Equal(t, text, dec)
}

func TestFail(t *testing.T) {
	text := "This is a test"
	pass := "password"
	badPass := "wrong"

	enc, err := secure.Encrypt(text, pass)
	assert.Nil(t, err)

	dec, err := secure.Decrypt(enc, badPass)
	assert.NotNil(t, err, err)
	assert.NotEqual(t, text, dec)
}
