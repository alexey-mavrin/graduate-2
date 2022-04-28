package client

import (
	"testing"

	"github.com/alexey-mavrin/graduate-2/internal/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_createUser(t *testing.T) {
	type want struct {
		respData common.AddUserResponse
	}

	tests := []struct {
		name     string
		userName string
		userPass string
		want     want
		wantID   int64
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name:     "Register a user",
			userName: "user1",
			userPass: "pass",
			wantID:   1,
			wantErr:  assert.NoError,
			want: want{
				respData: common.AddUserResponse{
					Name:   "user1",
					Status: "OK",
					ID:     1,
				},
			},
		},
		{
			name:     "Register the same user twice",
			userName: "user1",
			wantErr:  assert.Error,
			want: want{
				respData: common.AddUserResponse{
					Name:   "",
					Status: "User Already Exists",
					ID:     0,
				},
			},
		},
	}

	ts, err := newHTTPServer()
	require.NoError(t, err)
	defer ts.Close()

	for _, tt := range tests {
		clnt := NewClient(ts.URL, tt.userName, tt.userPass)
		t.Run(tt.name, func(t *testing.T) {
			id, err := clnt.RegisterUser("")
			tt.wantErr(t, err)
			assert.Equal(t, id, tt.wantID)
		})
	}
}
