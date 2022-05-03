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
