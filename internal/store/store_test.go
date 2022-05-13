package store

import (
	"testing"

	"github.com/alexey-mavrin/graduate-2/internal/common"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func dropCreateStore(t *testing.T) *Store {
	err := DropStore(defaultDBFile)
	assert.NoError(t, err)

	store, err := NewStore(defaultDBFile)
	assert.NoError(t, err)

	return store
}

func TestStore_AddUser(t *testing.T) {
	store := dropCreateStore(t)
	t.Run("Create user", func(t *testing.T) {
		user := "user1"
		_, err := store.AddUser(common.User{
			Name: user,
		})
		assert.NoError(t, err)

		// attempt to create the same user twice
		_, err = store.AddUser(common.User{
			Name: user,
		})
		assert.Error(t, err)
	})
}

func TestStore_CheckUserAuth(t *testing.T) {
	store := dropCreateStore(t)
	t.Run("Check user auth", func(t *testing.T) {
		user := "user1"
		wrongUser := "user2"
		pass := "pass1"
		wrongPass := "pass2"
		_, err := store.AddUser(common.User{
			Name:     user,
			Password: pass,
		})
		assert.NoError(t, err)

		ok, err := store.CheckUserAuth(user, pass)
		assert.NoError(t, err)
		assert.True(t, ok)

		ok, err = store.CheckUserAuth(user, wrongPass)
		assert.NoError(t, err)
		assert.False(t, ok)

		ok, err = store.CheckUserAuth(wrongUser, pass)
		assert.Error(t, err)
	})
}
