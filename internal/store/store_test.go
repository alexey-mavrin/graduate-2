package store

import (
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func dropCreateStore(t *testing.T) *Store {
	err := DropStore()
	assert.NoError(t, err)

	store, err := NewStore()
	assert.NoError(t, err)

	return store
}

func TestStore_StoreAccount(t *testing.T) {
	type args struct {
		user    string
		account Account
	}
	tests := []struct {
		name       string
		createUser bool
		args       args
		wantP      *int64
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name:       "store sample account",
			createUser: false,
			args: args{
				user: "user1",
				account: Account{
					URL:      "http://localhost",
					UserName: "user1",
					Password: "thePass",
				},
			},
			// try to store data w/o creating user first
			wantErr: assert.Error,
		},
		{
			name:       "store sample account",
			createUser: true,
			args: args{
				user: "user1",
				account: Account{
					URL:      "http://localhost",
					UserName: "user1",
					Password: "thePass",
				},
			},
			wantErr: assert.NoError,
		},
	}

	store := dropCreateStore(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.createUser {
				_, err := store.AddUser(User{
					Name: tt.args.user,
				})
				assert.NoError(t, err)
			}
			got, err := store.StoreAccount(tt.args.user, tt.args.account)
			tt.wantErr(t, err)
			if tt.wantP != nil && got != *tt.wantP {
				t.Errorf("Store.StoreAccount() = %v, want %v",
					got,
					*tt.wantP,
				)
			}
		})
	}
}

func TestStore_AddUser(t *testing.T) {
	store := dropCreateStore(t)
	t.Run("Create user", func(t *testing.T) {
		user := "user1"
		_, err := store.AddUser(User{
			Name: user,
		})
		assert.NoError(t, err)

		// attempt to create the same user twice
		_, err = store.AddUser(User{
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
		_, err := store.AddUser(User{
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

func TestStore_GetAccount(t *testing.T) {
	store := dropCreateStore(t)
	t.Run("Get single account", func(t *testing.T) {
		user := "user1"

		acc := Account{
			Name:     "local host",
			URL:      "http://localhost",
			UserName: user,
			Password: "secret1000",
		}

		_, err := store.AddUser(User{
			Name: user,
		})
		assert.NoError(t, err)

		id, err := store.StoreAccount(user, acc)
		assert.NoError(t, err)

		accRet, err := store.GetAccount(user, id)
		assert.NoError(t, err)

		assert.True(t, reflect.DeepEqual(acc, accRet))
	})
}

func TestStore_GetAccounts(t *testing.T) {
	store := dropCreateStore(t)
	t.Run("Get multiple accounts", func(t *testing.T) {
		user := "user1"
		name1 := "local host"
		name2 := "local host alt"
		url1 := "http://localhost"
		url2 := "http://localhost:8080"

		_, err := store.AddUser(User{
			Name: user,
		})
		assert.NoError(t, err)

		id1, err := store.StoreAccount(user, Account{
			UserName: user,
			Name:     name1,
			URL:      url1,
		})
		assert.NoError(t, err)

		id2, err := store.StoreAccount(user, Account{
			UserName: user,
			Name:     name2,
			URL:      url2,
		})
		assert.NoError(t, err)

		accs, err := store.GetAccounts(user)
		assert.NoError(t, err)

		wantAccs := make(Accounts)
		wantAccs[id1] = Account{
			UserName: user,
			Name:     name1,
			URL:      url1,
		}
		wantAccs[id2] = Account{
			UserName: user,
			Name:     name2,
			URL:      url2,
		}

		assert.True(t, reflect.DeepEqual(accs, wantAccs))
	})
}

func TestStore_DeleteAccount(t *testing.T) {
	store := dropCreateStore(t)
	t.Run("Delete account record", func(t *testing.T) {
		user := "user1"

		acc := Account{
			URL:      "http://localhost",
			Name:     "local host",
			UserName: user,
			Password: "secret1000",
		}

		_, err := store.AddUser(User{
			Name: user,
		})
		assert.NoError(t, err)

		id, err := store.StoreAccount(user, acc)
		assert.NoError(t, err)

		err = store.DeleteAccount(user, id)
		assert.NoError(t, err)

		// same attempt should result in error
		err = store.DeleteAccount(user, id)
		assert.Error(t, err)

		// attempt to delete non-existing account recourd
		// should result in error
		err = store.DeleteAccount(user, 999999)
		assert.Error(t, err)
	})
}
func TestStore_UpdateAccount(t *testing.T) {
	store := dropCreateStore(t)
	t.Run("Update account record", func(t *testing.T) {
		user := "user1"
		name := "local host"
		url1 := "http://localhost"
		url2 := "http://localhost:8080"

		_, err := store.AddUser(User{
			Name: user,
		})
		assert.NoError(t, err)

		id, err := store.StoreAccount(user, Account{
			URL:      url1,
			Name:     name,
			UserName: user,
		})
		assert.NoError(t, err)

		err = store.UpdateAccount(user, id, Account{
			URL:      url2,
			UserName: user,
		})
		assert.NoError(t, err)

		acc, err := store.GetAccount(user, id)
		assert.NoError(t, err)
		assert.Equal(t, acc.URL, url2, "Updated account url should change")
	})
}
