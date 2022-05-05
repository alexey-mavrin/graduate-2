package client

import (
	"os"
	"testing"

	"github.com/alexey-mavrin/graduate-2/internal/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_accounts(t *testing.T) {
	ts, err := newHTTPServer()
	require.NoError(t, err)
	defer ts.Close()

	userName := "user1"
	userPass := "pass"
	url := "http://localhost"
	url1 := "http://localhost:8080"

	clnt := NewClient(ts.URL, userName, userPass, "")

	_, err = clnt.RegisterUser("")
	assert.NoError(t, err)

	acc := common.Account{
		Name:     "local host",
		UserName: "u",
		Password: "p",
		URL:      url,
	}

	id, err := clnt.StoreAccount(acc)
	assert.NoError(t, err)

	accs, err := clnt.ListAccounts()
	respAcc := acc
	respAcc.Password = ""

	assert.Equal(t, respAcc, accs[id])

	updateAcc := acc
	updateAcc.URL = url1

	err = clnt.UpdateAccount(id, updateAcc)
	assert.NoError(t, err)

	gotAcc, err := clnt.GetAccount(id)
	assert.NoError(t, err)

	assert.Equal(t, updateAcc, gotAcc)

	err = clnt.DeleteAccount(id)
	assert.NoError(t, err)

	gotAcc, err = clnt.GetAccount(id)
	assert.Error(t, err)
}

func Test_accountsCache(t *testing.T) {
	ts, err := newHTTPServer()
	require.NoError(t, err)
	defer ts.Close()

	userName := "user1"
	userPass := "pass"
	url := "http://localhost"

	cacheName := "cache_storage.db"
	err = os.Remove(cacheName)
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}

	clnt := NewClient(ts.URL, userName, userPass, cacheName)

	_, err = clnt.RegisterUser("")
	assert.NoError(t, err)

	acc1 := common.Account{
		Name:     "local host",
		UserName: "u",
		Password: "p",
		URL:      url,
	}

	acc2 := common.Account{
		Name:     "local host",
		UserName: "u",
		Password: "p",
		URL:      url,
	}

	id1, err := clnt.StoreAccount(acc1)
	assert.NoError(t, err)
	id2, err := clnt.StoreAccount(acc2)
	assert.NoError(t, err)
	err = clnt.DeleteAccount(id2)
	assert.NoError(t, err)

	ts.Close()

	// should get account from cache
	gotAcc, err := clnt.GetAccount(id1)
	assert.NoError(t, err)
	assert.Equal(t, gotAcc, acc1)

	// make sure cache has no stall records
	ts.Close()
	gotAcc, err = clnt.GetAccount(id2)
	assert.Error(t, err)
}

func Test_accountsUpdateCache(t *testing.T) {
	ts, err := newHTTPServer()
	require.NoError(t, err)
	defer ts.Close()

	userName := "user1"
	userPass := "pass"
	url := "http://localhost"

	cacheName := "cache_storage.db"
	err = os.Remove(cacheName)
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}

	clnt := NewClient(ts.URL, userName, userPass, cacheName)

	_, err = clnt.RegisterUser("")
	assert.NoError(t, err)

	acc := common.Account{
		Name:     "local host",
		UserName: "u",
		Password: "p",
		URL:      url,
	}

	accUpd := common.Account{
		Name:     "local host",
		UserName: "u",
		Password: "pNew",
		URL:      url,
	}

	id, err := clnt.StoreAccount(acc)
	assert.NoError(t, err)

	err = clnt.UpdateAccount(id, accUpd)
	assert.NoError(t, err)

	ts.Close()

	// should get account from cache
	gotAcc, err := clnt.GetAccount(id)
	assert.NoError(t, err)
	assert.Equal(t, accUpd, gotAcc)
}

func Test_accountsDeleteCache(t *testing.T) {
	ts, err := newHTTPServer()
	require.NoError(t, err)
	defer ts.Close()

	userName := "user1"
	userPass := "pass"
	url := "http://localhost"

	cacheName := "cache_storage.db"
	err = os.Remove(cacheName)
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}

	clnt := NewClient(ts.URL, userName, userPass, cacheName)

	_, err = clnt.RegisterUser("")
	assert.NoError(t, err)

	acc := common.Account{
		Name:     "local host",
		UserName: "u",
		Password: "p",
		URL:      url,
	}

	id, err := clnt.StoreAccount(acc)
	assert.NoError(t, err)

	err = clnt.DeleteAccount(id)
	assert.NoError(t, err)

	ts.Close()

	// should NOT get account from cache
	_, err = clnt.GetAccount(id)
	assert.Error(t, err)
}
