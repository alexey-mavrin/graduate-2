package client

import (
	"bytes"
	"net/http"
	"time"

	"github.com/alexey-mavrin/graduate-2/internal/store"
)

const (
	defaultClientTimeout = time.Second * 1
)

// Client describes general client configuration
type Client struct {
	ServerAddr string
	UserName   string
	UserPass   string
	CacheFile  string
	Timeout    time.Duration
}

// NewClient returns new client
func NewClient(serverAddr string,
	userName string,
	userPass string,
	cacheFile string,
) *Client {
	store.DBFile = cacheFile
	return &Client{
		ServerAddr: serverAddr,
		UserName:   userName,
		UserPass:   userPass,
		CacheFile:  cacheFile,
		Timeout:    defaultClientTimeout,
	}
}

func (c *Client) prepaReq(method, path string, body []byte) (*http.Request, error) {
	b := bytes.NewReader(body)
	req, err := http.NewRequest(method, c.ServerAddr+path, b)
	if err != nil {
		return req, err
	}

	req.SetBasicAuth(c.UserName, c.UserPass)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}
