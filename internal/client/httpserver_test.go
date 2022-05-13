package client

import (
	"net/http/httptest"

	"github.com/alexey-mavrin/graduate-2/internal/server"
)

func newHTTPServer() (*httptest.Server, error) {
	err := server.DropServerStore("")
	if err != nil {
		return nil, err
	}
	err = server.InitStore("")
	if err != nil {
		return nil, err
	}
	router := server.NewRouter()
	ts := httptest.NewServer(router)

	return ts, nil
}
