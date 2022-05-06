package client

import (
	"os"
	"testing"

	"github.com/alexey-mavrin/graduate-2/internal/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_notes(t *testing.T) {
	ts, err := newHTTPServer()
	require.NoError(t, err)
	defer ts.Close()

	userName := "user1"
	userPass := "pass"
	text := "note text"
	text1 := "note text 1"

	clnt := NewClient(ts.URL, userName, userPass, "", false)

	_, err = clnt.RegisterUser("")
	assert.NoError(t, err)

	note := common.Note{
		Name: "local host",
		Text: text,
	}

	id, err := clnt.StoreNote(note)
	assert.NoError(t, err)

	notes, err := clnt.ListNotes()
	respNote := note
	respNote.Text = ""

	assert.Equal(t, respNote, notes[id])

	updateNote := note
	updateNote.Text = text1

	err = clnt.UpdateNote(id, updateNote)
	assert.NoError(t, err)

	gotNote, err := clnt.GetNote(id)
	assert.NoError(t, err)

	assert.Equal(t, updateNote, gotNote)

	err = clnt.DeleteNote(id)
	assert.NoError(t, err)

	gotNote, err = clnt.GetNote(id)
	assert.Error(t, err)
}

func Test_notesCache(t *testing.T) {
	ts, err := newHTTPServer()
	require.NoError(t, err)
	defer ts.Close()

	userName := "user1"
	userPass := "pass"
	text := "note text"

	cacheName := "cache_storage.db"
	err = os.Remove(cacheName)
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}

	clnt := NewClient(ts.URL, userName, userPass, cacheName, false)

	_, err = clnt.RegisterUser("")
	assert.NoError(t, err)

	note := common.Note{
		Name: "note1",
		Text: text,
	}

	id, err := clnt.StoreNote(note)
	assert.NoError(t, err)

	ts.Close()

	// should get note# from cache
	gotNote, err := clnt.GetNote(id)
	assert.NoError(t, err)
	assert.Equal(t, gotNote, note)
}

func Test_notesUpdateCache(t *testing.T) {
	ts, err := newHTTPServer()
	require.NoError(t, err)
	defer ts.Close()

	userName := "user1"
	userPass := "pass"
	text := "note text"

	cacheName := "cache_storage.db"
	err = os.Remove(cacheName)
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}

	clnt := NewClient(ts.URL, userName, userPass, cacheName, false)

	_, err = clnt.RegisterUser("")
	assert.NoError(t, err)

	note := common.Note{
		Name: "note 1",
		Text: text,
	}

	noteUpd := common.Note{
		Name: "note 1 update",
		Text: text,
	}

	id, err := clnt.StoreNote(note)
	assert.NoError(t, err)

	err = clnt.UpdateNote(id, noteUpd)
	assert.NoError(t, err)

	ts.Close()

	// should get note# from cache
	gotNote, err := clnt.GetNote(id)
	assert.NoError(t, err)
	assert.Equal(t, noteUpd, gotNote)
}

func Test_notesDeleteCache(t *testing.T) {
	ts, err := newHTTPServer()
	require.NoError(t, err)
	defer ts.Close()

	userName := "user1"
	userPass := "pass"
	text := "note text"

	cacheName := "cache_storage.db"
	err = os.Remove(cacheName)
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}

	clnt := NewClient(ts.URL, userName, userPass, cacheName, false)

	_, err = clnt.RegisterUser("")
	assert.NoError(t, err)

	note := common.Note{
		Name: "local host",
		Text: text,
	}

	id, err := clnt.StoreNote(note)
	assert.NoError(t, err)

	err = clnt.DeleteNote(id)
	assert.NoError(t, err)

	ts.Close()

	// should NOT get note from cache
	_, err = clnt.GetNote(id)
	assert.Error(t, err)
}
