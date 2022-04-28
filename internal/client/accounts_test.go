package client

import (
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

	clnt := NewClient(ts.URL, userName, userPass)

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
