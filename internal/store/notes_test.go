package store

import (
	"reflect"
	"testing"

	"github.com/alexey-mavrin/graduate-2/internal/common"
	"github.com/stretchr/testify/assert"
)

func TestStore_GetNote(t *testing.T) {
	store := dropCreateStore(t)
	t.Run("Get single note#", func(t *testing.T) {
		user := "user1"

		note := common.Note{
			Name: "local host",
			Text: "text note 1",
		}

		_, err := store.AddUser(common.User{
			Name: user,
		})
		assert.NoError(t, err)

		id, err := store.StoreNote(user, note)
		assert.NoError(t, err)

		noteRet, err := store.GetNote(user, id)
		assert.NoError(t, err)

		assert.True(t, reflect.DeepEqual(note, noteRet))
	})
}

func TestStore_ListNotes(t *testing.T) {
	store := dropCreateStore(t)
	t.Run("Get multiple notes", func(t *testing.T) {
		user := "user1"
		name1 := "note 1"
		name2 := "note 2"
		text1 := "text note 1"
		text2 := "text note 2"

		_, err := store.AddUser(common.User{
			Name: user,
		})
		assert.NoError(t, err)

		id1, err := store.StoreNote(user, common.Note{
			Name: name1,
			Text: text1,
		})
		assert.NoError(t, err)

		id2, err := store.StoreNote(user, common.Note{
			Name: name2,
			Text: text2,
		})
		assert.NoError(t, err)

		notes, err := store.ListNotes(user)
		assert.NoError(t, err)

		wantNotes := make(common.Notes)
		wantNotes[id1] = common.Note{
			Name: name1,
		}
		wantNotes[id2] = common.Note{
			Name: name2,
		}
		assert.True(t, reflect.DeepEqual(notes, wantNotes))
	})
}

func TestStore_DeleteNote(t *testing.T) {
	store := dropCreateStore(t)
	t.Run("Delete note record", func(t *testing.T) {
		user := "user1"

		note := common.Note{
			Text: "text note",
			Name: "local host",
		}

		_, err := store.AddUser(common.User{
			Name: user,
		})
		assert.NoError(t, err)

		id, err := store.StoreNote(user, note)
		assert.NoError(t, err)

		err = store.DeleteNote(user, id)
		assert.NoError(t, err)

		// same attempt should result in error
		err = store.DeleteNote(user, id)
		assert.Error(t, err)

		// attempt to delete non-existing note# recourd
		// should result in error
		err = store.DeleteNote(user, 999999)
		assert.Error(t, err)
	})
}

func TestStore_UpdateNote(t *testing.T) {
	store := dropCreateStore(t)
	t.Run("Update note record", func(t *testing.T) {
		user := "user1"
		name := "local host"
		text1 := "text note 1"
		text2 := "text note 2"

		_, err := store.AddUser(common.User{
			Name: user,
		})
		assert.NoError(t, err)

		id, err := store.StoreNote(user, common.Note{
			Name: name,
			Text: text1,
		})
		assert.NoError(t, err)

		err = store.UpdateNote(user, id, common.Note{
			Name: name,
			Text: text2,
		})
		assert.NoError(t, err)

		note, err := store.GetNote(user, id)
		assert.NoError(t, err)
		assert.Equal(t, note.Text, text2, "Updated note text should change")
	})
}
