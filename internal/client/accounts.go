package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/alexey-mavrin/graduate-2/internal/common"
)

// ListAccounts lists account for the current user
func (c *Client) ListAccounts() (common.Accounts, error) {
	var accs common.Accounts
	req, err := c.prepaReq(http.MethodGet, "/accounts", nil)
	if err != nil {
		return accs, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return accs, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf(
			"getting account list: http status %d",
			resp.StatusCode,
		)
		return accs, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return accs, err
	}
	err = json.Unmarshal(respBody, &accs)
	if err != nil {
		return accs, err
	}
	return accs, nil
}

// DeleteAccount returns account record with the given id
func (c *Client) DeleteAccount(id int64) error {
	path := fmt.Sprintf("/accounts/%d", id)
	req, err := c.prepaReq(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf(
			"delete account: http status %d",
			resp.StatusCode,
		)
		return err
	}

	var status common.StoreAccountResponse
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(respBody, &status)
	if err != nil {
		return err
	}
	if status.Status != "OK" {
		return errors.New(status.Status)
	}
	return nil
}

// GetAccount returns account record with the given id
func (c *Client) GetAccount(id int64) (common.Account, error) {
	var acc common.Account

	path := fmt.Sprintf("/accounts/%d", id)
	req, err := c.prepaReq(http.MethodGet, path, nil)
	if err != nil {
		return acc, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return acc, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf(
			"get account: http status %d",
			resp.StatusCode,
		)
		return acc, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return acc, err
	}
	err = json.Unmarshal(respBody, &acc)
	if err != nil {
		return acc, err
	}
	return acc, nil
}

// UpdateAccount updates account record with the given id
func (c *Client) UpdateAccount(id int64, acc common.Account) error {
	body, err := json.Marshal(acc)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/accounts/%d", id)
	req, err := c.prepaReq(http.MethodPut, path, body)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf(
			"updating account: http status %d",
			resp.StatusCode,
		)
		return err
	}

	var status common.StoreAccountResponse
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(respBody, &status)
	if err != nil {
		return err
	}
	if status.Status != "OK" {
		return errors.New(status.Status)
	}
	return nil
}

// StoreAccount stores account record
func (c *Client) StoreAccount(acc common.Account) (int64, error) {
	body, err := json.Marshal(acc)
	if err != nil {
		return 0, err
	}

	req, err := c.prepaReq(http.MethodPost, "/accounts", body)
	if err != nil {
		return 0, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf(
			"storing account: http status %d",
			resp.StatusCode,
		)
		return 0, err
	}

	var status common.StoreAccountResponse
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	err = json.Unmarshal(respBody, &status)
	if err != nil {
		return 0, err
	}
	if status.Status != "OK" {
		return 0, errors.New(status.Status)
	}
	return status.ID, nil
}