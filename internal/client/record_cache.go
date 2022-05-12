package client

import (
	"errors"
	"log"

	"github.com/alexey-mavrin/graduate-2/internal/common"
	"github.com/alexey-mavrin/graduate-2/internal/store"
)

func (c *Client) cacheDeleteRecordID(id int64) error {
	if c.CacheFile == "" {
		return nil
	}
	cache, err := store.NewStore()
	if err != nil {
		log.Printf("cannot get local store: %v", err)
		return err
	}
	err = cache.DeleteRecordID(c.UserName, id)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) cacheRecordID(storeID int64, record common.Record) error {
	if c.CacheFile == "" {
		return nil
	}
	cache, err := store.NewStore()
	if err != nil {
		log.Printf("cannot get local store: %v", err)
		return err
	}

	_, err = cache.AddUser(common.User{
		Name: c.UserName,
	})
	if err != nil && err != store.ErrAlreadyExists {
		return err
	}

	records, err := cache.ListRecords(c.UserName)

	for id := range records {
		if records[id].Name != record.Name ||
			records[id].Type != record.Type {
			continue
		}
		if id != storeID {
			return errors.New("local cache is out of sync")
		}
		err = cache.DeleteRecordID(c.UserName, id)
		if err != nil {
			return err
		}
		break
	}

	err = cache.StoreRecordID(storeID, c.UserName, record)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) cacheGetRecordID(id int64) (common.Record, error) {
	if c.CacheFile == "" {
		return common.Record{}, nil
	}
	cache, err := store.NewStore()
	if err != nil {
		log.Printf("cannot get local store: %v", err)
		return common.Record{}, err
	}
	return cache.GetRecordID(c.UserName, id)
}

func (c *Client) cacheListRecords() (common.Records, error) {
	records := make(common.Records)
	if c.CacheFile == "" {
		return records, nil
	}
	cache, err := store.NewStore()
	if err != nil {
		log.Printf("cannot get local store: %v", err)
		return records, err
	}
	return cache.ListRecords(c.UserName)
}

func (c *Client) cacheListRecordsType(t common.RecordType) (common.Records, error) {
	records := make(common.Records)
	if c.CacheFile == "" {
		return records, nil
	}
	cache, err := store.NewStore()
	if err != nil {
		log.Printf("cannot get local store: %v", err)
		return records, err
	}
	return cache.ListRecordsType(c.UserName, t)
}
