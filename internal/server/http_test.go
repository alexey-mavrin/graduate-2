package server

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func testHTTPRequest(t *testing.T,
	ts *httptest.Server,
	method string,
	path string,
	body string,
	user string,
	pass string,
) (*http.Response, string) {
	bodyReader := strings.NewReader(body)
	req, err := http.NewRequest(method, ts.URL+path, bodyReader)
	require.NoError(t, err)

	req.Header.Add("Content-Type", "application/json")
	if user != "" && pass != "" {
		req.SetBasicAuth(user, pass)
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}
