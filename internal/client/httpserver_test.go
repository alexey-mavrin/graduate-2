package client

import (
	"net/http/httptest"

	"github.com/alexey-mavrin/graduate-2/internal/server"
	"github.com/alexey-mavrin/graduate-2/internal/store"
)

func newHTTPServer() (*httptest.Server, error) {
	err := store.DropStore()
	if err != nil {
		return nil, err
	}
	router := server.NewRouter()
	ts := httptest.NewServer(router)

	return ts, nil
}
