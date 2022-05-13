package client

import (
	"github.com/alexey-mavrin/graduate-2/internal/common"
	"github.com/alexey-mavrin/graduate-2/internal/store"
)

func (c *Client) cacheDeleteRecordID(id int64) error {
	if c.CacheFile == "" {
		return nil
	}
	err := c.Store.DeleteRecordID(c.UserName, id)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) cacheRecordID(storeID int64, record common.Record) error {
	if c.CacheFile == "" {
		return nil
	}

	_, err := c.Store.AddUser(common.User{
		Name: c.UserName,
	})
	if err != nil && err != store.ErrAlreadyExists {
		return err
	}

	err = c.Store.DeleteRecordID(c.UserName, storeID)
	if err != nil && err != store.ErrNotFound {
		return err
	}

	err = c.Store.StoreRecordID(storeID, c.UserName, record)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) cacheGetRecordID(id int64) (common.Record, error) {
	if c.CacheFile == "" {
		return common.Record{}, nil
	}
	return c.Store.GetRecordID(c.UserName, id)
}

func (c *Client) cacheListRecords() (common.Records, error) {
	records := make(common.Records)
	if c.CacheFile == "" {
		return records, nil
	}
	return c.Store.ListRecords(c.UserName)
}

func (c *Client) cacheListRecordsType(t common.RecordType) (common.Records, error) {
	records := make(common.Records)
	if c.CacheFile == "" {
		return records, nil
	}
	return c.Store.ListRecordsType(c.UserName, t)
}
