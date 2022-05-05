package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexey-mavrin/graduate-2/internal/common"
	"github.com/alexey-mavrin/graduate-2/internal/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Account(t *testing.T) {
	require.NoError(t, store.DropStore())
	router := NewRouter()
	ts := httptest.NewServer(router)
	defer ts.Close()

	t.Run("Store account record", func(t *testing.T) {
		createUserBody := `{"name":"user1", "password":"pass"}`

		acc := common.Account{
			Name:     "localhost",
			URL:      "http://localhost",
			UserName: "u",
			Password: "p",
		}

		accBody, _ := json.Marshal(acc)

		resp, _ := testHTTPRequest(t,
			ts,
			http.MethodPost,
			"/users",
			createUserBody,
			// empty user and pass
			"",
			"",
		)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		createResp, _ := testHTTPRequest(t,
			ts,
			http.MethodPost,
			"/records/account",
			string(accBody),
			// correct user and pass
			"user1",
			"pass",
		)
		defer createResp.Body.Close()
		assert.Equal(t, http.StatusOK, createResp.StatusCode)

		getResp, getRespBody := testHTTPRequest(t,
			ts,
			http.MethodGet,
			"/records/account/1",
			"",
			// correct user and pass
			"user1",
			"pass",
		)
		defer getResp.Body.Close()
		assert.Equal(t, http.StatusOK, getResp.StatusCode)
		var gotAcc common.Account
		err := json.Unmarshal([]byte(getRespBody), &gotAcc)
		assert.NoError(t, err)
		assert.Equal(t, acc, gotAcc)

		listResp, listRespBody := testHTTPRequest(t,
			ts,
			http.MethodGet,
			"/records/account",
			"",
			// correct user and pass
			"user1",
			"pass",
		)
		assert.Equal(t, http.StatusOK, listResp.StatusCode)
		var listAccs common.Accounts
		err = json.Unmarshal([]byte(listRespBody), &listAccs)
		// list returns accounts w/o passwords
		wantAcc := acc
		wantAcc.Password = ""
		assert.NoError(t, err)
		assert.Equal(t, wantAcc, listAccs[1])

		delResp, _ := testHTTPRequest(t,
			ts,
			http.MethodDelete,
			"/records/account/1",
			"",
			// correct user and pass
			"user1",
			"pass",
		)
		assert.Equal(t, http.StatusOK, delResp.StatusCode)

		getResp2, _ := testHTTPRequest(t,
			ts,
			http.MethodGet,
			"/records/account/1",
			"",
			// correct user and pass
			"user1",
			"pass",
		)
		defer getResp.Body.Close()
		assert.Equal(t, http.StatusNotFound, getResp2.StatusCode)
	})
}
