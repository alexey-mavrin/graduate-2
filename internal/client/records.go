package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/alexey-mavrin/graduate-2/internal/common"
)

// ListRecordsByType lists account for the current user
func (c *Client) ListRecordsByType(t common.RecordType) (common.Records, error) {
	var records common.Records
	path := fmt.Sprintf("/records/by_type/%s", t)
	req, err := c.prepaReq(http.MethodGet, path, nil)
	if err != nil {
		return records, err
	}

	client := c.httpClient()
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("cannot contact the server: %v, trying local cache", err)
		records, err := c.cacheListRecordsByType(t)
		return records, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf(
			"getting %s list: http status %d",
			t, resp.StatusCode,
		)
		return records, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return records, err
	}

	err = json.Unmarshal(respBody, &records)
	if err != nil {
		return records, err
	}
	return records, nil
}

// DeleteRecordByID returns account record with the given id
func (c *Client) DeleteRecordByID(id int64) error {
	path := fmt.Sprintf("/records/%d", id)
	req, err := c.prepaReq(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	client := c.httpClient()
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var status common.StoreRecordResponse
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
			"delete record %d: http status %d: %s",
			id, resp.StatusCode, status.Status,
		)
		return err
	}

	err = c.cacheDeleteRecordByID(id)
	if err != nil {
		log.Print("cannot delete record from cache: ", err)
	}

	return nil
}

// GetRecordByID returns account record with the given id
func (c *Client) GetRecordByID(id int64) (common.Record, error) {
	var record common.Record

	path := fmt.Sprintf("/records/%d", id)
	req, err := c.prepaReq(http.MethodGet, path, nil)
	if err != nil {
		return record, err
	}

	client := c.httpClient()
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("cannot contact the server: %v, trying local cache", err)
		record, err := c.cacheGetRecordByID(id)
		return record, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf(
			"get account: http status %d",
			resp.StatusCode,
		)
		return record, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return record, err
	}

	err = json.Unmarshal(respBody, &record)
	if err != nil {
		return record, err
	}
	err = c.cacheRecordWithID(id, record)

	if err != nil {
		log.Printf("cache record: %v", err)
	}

	return record, nil
}

// UpdateRecordByID updates account record with the given id
func (c *Client) UpdateRecordByID(id int64, record common.Record) error {
	body, err := json.Marshal(record)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/records/%d", id)
	req, err := c.prepaReq(http.MethodPut, path, body)
	if err != nil {
		return err
	}

	client := c.httpClient()
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var status common.StoreRecordResponse
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

	err = c.cacheRecordWithID(id, record)
	if err != nil {
		log.Printf("cache %s: %v", record.Type, err)
	}

	return nil
}

// StoreRecord stores account record
func (c *Client) StoreRecord(record common.Record) (int64, error) {
	body, err := json.Marshal(record)
	if err != nil {
		return 0, err
	}

	path := fmt.Sprintf("/records")
	req, err := c.prepaReq(http.MethodPost, path, body)
	if err != nil {
		return 0, err
	}

	client := c.httpClient()
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var status common.StoreRecordResponse
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

	err = c.cacheRecordWithID(status.ID, record)
	if err != nil {
		log.Printf("cache %s: %v", record.Type, err)
	}

	return status.ID, nil
}
