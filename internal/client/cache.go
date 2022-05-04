package client

import (
	"errors"
	"log"

	"github.com/alexey-mavrin/graduate-2/internal/common"
	"github.com/alexey-mavrin/graduate-2/internal/store"
)

func (c *Client) cacheDeleteAccount(id int64) error {
	if c.CacheFile == "" {
		return nil
	}
	cache, err := store.NewStore()
	if err != nil {
		log.Printf("cannot get local store: %v", err)
		return err
	}
	err = cache.DeleteAccount(c.UserName, id)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) cacheDeleteNote(id int64) error {
	if c.CacheFile == "" {
		return nil
	}
	cache, err := store.NewStore()
	if err != nil {
		log.Printf("cannot get local store: %v", err)
		return err
	}
	err = cache.DeleteNote(c.UserName, id)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) cacheAccount(storeID int64, acc common.Account) error {
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

	accs, err := cache.ListAccounts(c.UserName)

	for id := range accs {
		if accs[id].Name != acc.Name {
			continue
		}
		if id != storeID {
			return errors.New("local cache is out of sync")
		}
		err = cache.DeleteAccount(c.UserName, id)
		if err != nil {
			return err
		}
		break
	}

	err = cache.StoreAccountID(storeID, c.UserName, acc)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) cacheNote(storeID int64, note common.Note) error {
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

	notes, err := cache.ListNotes(c.UserName)

	for id := range notes {
		if notes[id].Name != note.Name {
			continue
		}
		if id != storeID {
			return errors.New("local cache is out of sync")
		}
		err = cache.DeleteNote(c.UserName, id)
		if err != nil {
			return err
		}
		break
	}

	err = cache.StoreNoteID(storeID, c.UserName, note)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) cacheGetAccount(id int64) (common.Account, error) {
	if c.CacheFile == "" {
		return common.Account{}, nil
	}
	cache, err := store.NewStore()
	if err != nil {
		log.Printf("cannot get local store: %v", err)
		return common.Account{}, err
	}
	return cache.GetAccount(c.UserName, id)
}

func (c *Client) cacheGetNote(id int64) (common.Note, error) {
	if c.CacheFile == "" {
		return common.Note{}, nil
	}
	cache, err := store.NewStore()
	if err != nil {
		log.Printf("cannot get local store: %v", err)
		return common.Note{}, err
	}
	return cache.GetNote(c.UserName, id)
}
