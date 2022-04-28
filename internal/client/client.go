package client

import (
	"bytes"
	"net/http"
)

// Client describes general client configuration
type Client struct {
	ServerAddr string
	UserName   string
	UserPass   string
}

// NewClient returns new client
func NewClient(serverAddr string, userName string, userPass string) *Client {
	return &Client{
		ServerAddr: serverAddr,
		UserName:   userName,
		UserPass:   userPass,
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
