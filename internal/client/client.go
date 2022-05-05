package client

import (
	"bytes"
	"crypto/tls"
	"net/http"
	"time"

	"github.com/alexey-mavrin/graduate-2/internal/store"
)

const (
	defaultClientTimeout = time.Second * 1
)

//go:generate go run tmpl/generator.go Account
//go:generate go run tmpl/generator.go Note

// Client describes general client configuration
type Client struct {
	ServerAddr    string
	UserName      string
	UserPass      string
	CacheFile     string
	Timeout       time.Duration
	HTTPSInsecure bool
}

// NewClient returns new client
func NewClient(serverAddr string,
	userName string,
	userPass string,
	cacheFile string,
	httpsInsecure bool,
) *Client {
	store.DBFile = cacheFile
	return &Client{
		ServerAddr:    serverAddr,
		UserName:      userName,
		UserPass:      userPass,
		CacheFile:     cacheFile,
		Timeout:       defaultClientTimeout,
		HTTPSInsecure: httpsInsecure,
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

func (c *Client) httpClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: c.HTTPSInsecure,
		},
	}
	client := &http.Client{
		Timeout:   c.Timeout,
		Transport: tr,
	}
	return client
}
