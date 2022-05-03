package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

	client := &http.Client{
		Timeout: c.Timeout,
	}
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

	client := &http.Client{
		Timeout: c.Timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var status common.StoreAccountResponse
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(respBody, &status)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf(
			"delete account: http status %d: %s",
			resp.StatusCode,
			status.Status,
		)
		return err
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

	client := &http.Client{
		Timeout: c.Timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("cannot contact the server: %v, trying local cache", err)
		return c.cacheGetAccount(id)
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

	err = c.cacheAccount(id, acc)
	if err != nil {
		log.Printf("cache account: %v", err)
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

	client := &http.Client{
		Timeout: c.Timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var status common.StoreAccountResponse
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(respBody, &status)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf(
			"updating account: http status %d: %s",
			resp.StatusCode,
			status.Status,
		)
		return err
	}

	err = c.cacheAccount(id, acc)
	if err != nil {
		log.Printf("cache account: %v", err)
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

	client := &http.Client{
		Timeout: c.Timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var status common.StoreAccountResponse
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	err = json.Unmarshal(respBody, &status)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf(
			"storing account: http status %d: %s",
			resp.StatusCode,
			status.Status,
		)
		return 0, err
	}

	err = c.cacheAccount(status.ID, acc)
	if err != nil {
		log.Printf("cache account: %v", err)
	}

	return status.ID, nil
}
