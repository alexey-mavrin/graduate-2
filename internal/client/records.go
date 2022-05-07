package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/alexey-mavrin/graduate-2/internal/common"
)

// ListRecords lists account for the current user
func (c *Client) listRecords(t common.RecordType) (common.Records, error) {
	var records common.Records
	path := fmt.Sprintf("/records/%s", t)
	req, err := c.prepaReq(http.MethodGet, path, nil)
	if err != nil {
		return records, err
	}

	client := c.httpClient()
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("cannot contact the server: %v, trying local cache", err)
		switch t {
		case common.AccountRecord:
			accounts, err := c.cacheListAccounts()
			records.Accounts = &accounts
			return records, err
		case common.NoteRecord:
			notes, err := c.cacheListNotes()
			records.Notes = &notes
			return records, err
		case common.CardRecord:
			cards, err := c.cacheListCards()
			records.Cards = &cards
			return records, err
		default:
			return records, fmt.Errorf("unknown record type %s", t)
		}
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

	var parseErr error
	records.Type = t
	switch t {
	case common.AccountRecord:
		accounts := make(common.Accounts)
		parseErr = json.Unmarshal(respBody, &accounts)
		records.Accounts = &accounts
	case common.NoteRecord:
		notes := make(common.Notes)
		parseErr = json.Unmarshal(respBody, &notes)
		records.Notes = &notes
	case common.CardRecord:
		cards := make(common.Cards)
		parseErr = json.Unmarshal(respBody, &cards)
		records.Cards = &cards
	default:
		return records, fmt.Errorf("unknown record type %s", t)
	}
	if parseErr != nil {
		return records, parseErr
	}
	return records, nil
}

// DeleteRecord returns account record with the given id
func (c *Client) deleteRecord(id int64, t common.RecordType) error {
	path := fmt.Sprintf("/records/%s/%d", t, id)
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
			"delete %s: http status %d: %s",
			t, resp.StatusCode, status.Status,
		)
		return err
	}

	return nil
}

// getRecord returns account record with the given id
func (c *Client) getRecord(id int64, t common.RecordType) (common.Record, error) {
	var record common.Record

	path := fmt.Sprintf("/records/%s/%d", t, id)
	req, err := c.prepaReq(http.MethodGet, path, nil)
	if err != nil {
		return record, err
	}

	client := c.httpClient()
	resp, err := client.Do(req)
	record.Type = t
	if err != nil {
		log.Printf("cannot contact the server: %v, trying local cache", err)
		switch t {
		case common.AccountRecord:
			account, err := c.cacheGetAccount(id)
			record.Account = &account
			return record, err
		case common.NoteRecord:
			note, err := c.cacheGetNote(id)
			record.Note = &note
			return record, err
		case common.CardRecord:
			card, err := c.cacheGetCard(id)
			record.Card = &card
			return record, err
		default:
			return record, fmt.Errorf("unknown record type %s", t)
		}
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

	var parseErr, cacheErr error
	switch t {
	case common.AccountRecord:
		var account common.Account
		parseErr = json.Unmarshal(respBody, &account)
		if parseErr != nil {
			break
		}
		record.Account = &account
		cacheErr = c.cacheAccount(id, account)
	case common.NoteRecord:
		var note common.Note
		parseErr = json.Unmarshal(respBody, &note)
		if parseErr != nil {
			break
		}
		record.Note = &note
		cacheErr = c.cacheNote(id, note)
	case common.CardRecord:
		var card common.Card
		parseErr = json.Unmarshal(respBody, &card)
		if parseErr != nil {
			break
		}
		record.Card = &card
		cacheErr = c.cacheCard(id, card)
	default:
		parseErr = fmt.Errorf("unknown type: %s", t)
		cacheErr = fmt.Errorf("unknown type: %s", t)
	}
	if parseErr != nil {
		return record, parseErr
	}

	if cacheErr != nil {
		log.Printf("cache record: %v", err)
	}

	return record, nil
}

// updateRecord updates account record with the given id
func (c *Client) updateRecord(id int64, record common.Record) error {
	body := make([]byte, 0)
	var encodeErr error
	switch record.Type {
	case common.AccountRecord:
		body, encodeErr = json.Marshal(*record.Account)
	case common.NoteRecord:
		body, encodeErr = json.Marshal(*record.Note)
	case common.CardRecord:
		body, encodeErr = json.Marshal(*record.Card)
	default:
		encodeErr = fmt.Errorf("unknown record type %s", record.Type)
	}
	if encodeErr != nil {
		return encodeErr
	}

	path := fmt.Sprintf("/records/%s/%d", record.Type, id)
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

	var cacheErr error
	switch record.Type {
	case common.AccountRecord:
		cacheErr = c.cacheAccount(id, *record.Account)
	case common.NoteRecord:
		cacheErr = c.cacheNote(id, *record.Note)
	case common.CardRecord:
		cacheErr = c.cacheCard(id, *record.Card)
	default:
		cacheErr = fmt.Errorf("unknown record type: %s", record.Type)
	}

	if cacheErr != nil {
		log.Printf("cache %s: %v", record.Type, cacheErr)
	}

	return nil
}

// storeRecord stores account record
func (c *Client) storeRecord(record common.Record) (int64, error) {
	body := make([]byte, 0)
	var encodeErr error
	switch record.Type {
	case common.AccountRecord:
		body, encodeErr = json.Marshal(*record.Account)
	case common.NoteRecord:
		body, encodeErr = json.Marshal(*record.Note)
	case common.CardRecord:
		body, encodeErr = json.Marshal(*record.Card)
	default:
		encodeErr = fmt.Errorf("unknown record type %s", record.Type)
	}
	if encodeErr != nil {
		return 0, encodeErr
	}

	path := fmt.Sprintf("/records/%s", record.Type)
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

	var cacheErr error
	switch record.Type {
	case common.AccountRecord:
		cacheErr = c.cacheAccount(status.ID, *record.Account)
	case common.NoteRecord:
		cacheErr = c.cacheNote(status.ID, *record.Note)
	case common.CardRecord:
		cacheErr = c.cacheCard(status.ID, *record.Card)
	default:
		cacheErr = fmt.Errorf("unknown record type: %s", record.Type)
	}

	if cacheErr != nil {
		log.Printf("cache %s: %v", record.Type, cacheErr)
	}

	return status.ID, nil
}
