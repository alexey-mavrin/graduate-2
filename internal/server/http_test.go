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
	router http.Handler,
	method string,
	path string,
	body string,
	user string,
	pass string,
) (*http.Response, string) {
	bodyReader := strings.NewReader(body)
	req := httptest.NewRequest(method, path, bodyReader)

	req.Header.Add("Content-Type", "application/json")
	if user != "" && pass != "" {
		req.SetBasicAuth(user, pass)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	resp := w.Result()
	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}
