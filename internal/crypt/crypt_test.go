package crypt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_EncryptDecrypt(t *testing.T) {
	clearText1 := "o la la"
	clearText2 := "blah-vlah-blah-vlah-blah-vlah-blah-vlah-blah-vlah"

	var key1 [32]byte
	copy(key1[:], "01234567890123456789012345678901")
	key2 := MakeKey("this is a key phrase, could be long or short")

	cyperText1, err := Encrypt(key1, []byte(clearText1))
	assert.NoError(t, err)
	cyperText2, err := Encrypt(key2, []byte(clearText2))
	assert.NoError(t, err)

	assert.NotEqual(t, cyperText1, clearText1)
	assert.NotEqual(t, cyperText2, clearText2)

	newClearText1, err := Decrypt(key1, cyperText1)
	assert.NoError(t, err)
	newClearText2, err := Decrypt(key2, cyperText2)
	assert.NoError(t, err)

	assert.Equal(t, string(newClearText1), clearText1)
	assert.Equal(t, string(newClearText2), clearText2)
}
