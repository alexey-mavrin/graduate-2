package crypt

import (
	"testing"

	"github.com/alexey-mavrin/graduate-2/internal/common"
	"github.com/stretchr/testify/assert"
)

func Test_cryptAccount(t *testing.T) {
	key := MakeKey("qwerty")
	acc := common.Account{
		Name:     "name",
		UserName: "user_name",
		Password: "pass1",
		URL:      "http://localhost",
		Meta:     "yo-ho-ho",
	}
	eAcc, err := EncryptAccount(key, acc)
	assert.NoError(t, err)
	assert.NotEqual(t, eAcc, acc)
	decr, err := DecryptAccount(key, eAcc)
	assert.NoError(t, err)
	assert.Equal(t, acc, decr)
}

func Test_cryptNote(t *testing.T) {
	key := MakeKey("qwerty")
	acc := common.Note{
		Name: "name",
		Text: "some text",
		Meta: "yo-ho-ho",
	}
	eAcc, err := EncryptNote(key, acc)
	assert.NoError(t, err)
	assert.NotEqual(t, eAcc, acc)
	decr, err := DecryptNote(key, eAcc)
	assert.NoError(t, err)
	assert.Equal(t, acc, decr)
}
