package client

import (
	"github.com/alexey-mavrin/graduate-2/internal/common"
	"github.com/alexey-mavrin/graduate-2/internal/store"
)

// ListAccounts lists account for the current user
func (c *Client) ListAccounts() (common.Accounts, error) {
	records, err := c.listRecords(common.AccountRecord)
	if err != nil {
		return common.Accounts{}, err
	}
	return *records.Accounts, nil
}

// DeleteAccount returns account record with the given id
func (c *Client) DeleteAccount(id int64) error {
	err := c.deleteRecord(id, common.AccountRecord)
	if err != nil {
		return err
	}
	err = c.cacheDeleteAccount(id)
	if err != nil && err != store.ErrNotFound {
		return err
	}
	return nil
}

// GetAccount returns account record with the given id
func (c *Client) GetAccount(id int64) (common.Account, error) {
	record, err := c.getRecord(id, common.AccountRecord)
	if err != nil {
		return common.Account{}, err
	}
	return *record.Account, nil
}

// UpdateAccount updates account record with the given id
func (c *Client) UpdateAccount(id int64, acc common.Account) error {
	record := common.Record{
		Type:    common.AccountRecord,
		Account: &acc,
	}
	return c.updateRecord(id, record)
}

// StoreAccount stores account record
func (c *Client) StoreAccount(acc common.Account) (int64, error) {
	record := common.Record{
		Type:    common.AccountRecord,
		Account: &acc,
	}
	return c.storeRecord(record)
}
