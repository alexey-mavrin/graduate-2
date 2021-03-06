package server

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/alexey-mavrin/graduate-2/internal/common"
	"github.com/alexey-mavrin/graduate-2/internal/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_createUser(t *testing.T) {
	type want struct {
		respData common.AddUserResponse
		code     int
	}

	tests := []struct {
		name   string
		method string
		body   string
		user   string
		pass   string
		want   want
	}{
		{
			name:   "Create a user",
			method: http.MethodPost,
			want: want{
				code: 200,
				respData: common.AddUserResponse{
					Name:   "user1",
					Status: "OK",
					ID:     1,
				},
			},
			body: `{"name":"user1"}`,
		},
		{
			name:   "Create the same user twice",
			method: http.MethodPost,
			want: want{
				code: http.StatusBadRequest,
				respData: common.AddUserResponse{
					Name:   "",
					Status: "User Already Exists",
					ID:     0,
				},
			},
			body: `{"name":"user1"}`,
		},
	}

	require.NoError(t, store.DropStore(defaultStoreFile))
	err := InitStore(defaultStoreFile)
	assert.NoError(t, err)
	router := NewRouter()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, body := testHTTPRequest(t,
				router,
				tt.method,
				"/users",
				tt.body,
				tt.user,
				tt.pass,
			)
			defer resp.Body.Close()
			assert.Equal(t, tt.want.code, resp.StatusCode)
			var respData common.AddUserResponse
			err := json.Unmarshal([]byte(body), &respData)
			assert.NoError(t, err)
			assert.Equal(t, tt.want.respData, respData)
		})
	}
}

func Test_pingHandler(t *testing.T) {
	body := `{"name":"user1", "password":"pass"}`

	require.NoError(t, store.DropStore(defaultStoreFile))
	err := InitStore(defaultStoreFile)
	assert.NoError(t, err)
	router := NewRouter()

	resp, _ := testHTTPRequest(t,
		router,
		http.MethodPost,
		"/users",
		body,
		// empty user and pass
		"",
		"",
	)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp1, _ := testHTTPRequest(t,
		router,
		http.MethodGet,
		"/ping",
		"",
		// correct user and pass
		"user1",
		"pass",
	)
	defer resp1.Body.Close()
	assert.Equal(t, http.StatusOK, resp1.StatusCode)

	resp2, _ := testHTTPRequest(t,
		router,
		http.MethodGet,
		"/ping",
		"",
		// incorrect user and pass
		"user1",
		"passX",
	)
	defer resp2.Body.Close()
	assert.Equal(t, http.StatusForbidden, resp2.StatusCode)
}
