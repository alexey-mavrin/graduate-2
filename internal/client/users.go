package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/alexey-mavrin/graduate-2/internal/common"
)

// RegisterUser attempts to register current user into the system
// returns new user id and error
func (c Client) RegisterUser(fullName string) (int64, error) {
	user := common.User{
		Name:     c.UserName,
		FullName: fullName,
		Password: c.UserPass,
	}

	var body bytes.Buffer

	err := json.NewEncoder(&body).Encode(user)
	if err != nil {
		return 0, err
	}
	url := c.ServerAddr + "/users"

	req, err := http.NewRequest(http.MethodPost, url, &body)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("registing %s, http status %d",
			c.UserName,
			resp.StatusCode,
		)
		return 0, fmt.Errorf("http error %d", resp.StatusCode)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var addUserResp common.StoreAccountResponse
	err = json.Unmarshal(respBody, &addUserResp)
	if err != nil {
		return 0, err
	}

	return addUserResp.ID, nil
}
