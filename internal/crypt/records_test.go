package crypt

import (
	"testing"

	"github.com/alexey-mavrin/graduate-2/internal/common"
	"github.com/stretchr/testify/assert"
)

func Test_cryptRecord(t *testing.T) {
	key := MakeKey("qwerty")
	record := common.Record{
		Name:   "name",
		Type:   common.NoteRecord,
		Opaque: "1111",
		Meta:   "yo-ho-ho",
	}
	eRecord, err := EncryptRecord(key, record)
	assert.NoError(t, err)
	assert.NotEqual(t, eRecord, record)
	decr, err := DecryptRecord(key, eRecord)
	assert.NoError(t, err)
	assert.Equal(t, record, decr)
}
