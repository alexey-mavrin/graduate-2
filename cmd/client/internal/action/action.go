package action

import (
	"errors"
	"fmt"
	"log"

	"github.com/alexey-mavrin/graduate-2/cmd/client/internal/config"
	"github.com/alexey-mavrin/graduate-2/internal/client"
	"github.com/alexey-mavrin/graduate-2/internal/common"
)

func actUser(subop config.OpSubtype, user common.User) error {
	clnt := client.NewClient(config.Cfg.ServerAddr,
		config.Cfg.UserName,
		config.Cfg.Password,
		config.Cfg.CacheFile,
		config.Cfg.HTTPSInsecure,
	)
	switch subop {
	case config.OpSubtypeUserRegister:
		id, err := clnt.RegisterUser(config.Cfg.FullName)
		if err != nil {
			return err
		}
		log.Printf("user is registered with id %d", id)
	case config.OpSubtypeUserVerify:
		err := clnt.VerifyUser()
		if err != nil {
			return err
		}
		log.Printf("user is verified")
	}
	return nil
}

func actAccount(subop config.OpSubtype, acc common.Account) error {
	clnt := client.NewClient(config.Cfg.ServerAddr,
		config.Cfg.UserName,
		config.Cfg.Password,
		config.Cfg.CacheFile,
		config.Cfg.HTTPSInsecure,
	)
	switch subop {
	case config.OpSubtypeRecordStore:
		id, err := clnt.StoreAccount(config.Op.Account)
		if err != nil {
			return err
		}
		fmt.Printf("accout record stored with id %d\n", id)
	case config.OpSubtypeRecordGet:
		acc, err := clnt.GetAccount(config.Op.AccountID)
		if err != nil {
			return err
		}
		fmt.Println(acc)
	case config.OpSubtypeRecordList:
		accs, err := clnt.ListAccounts()
		if err != nil {
			return err
		}
		fmt.Println(accs)
	case config.OpSubtypeRecordUpdate:
		err := clnt.UpdateAccount(config.Op.AccountID, config.Op.Account)
		if err != nil {
			return err
		}
		fmt.Println("account updated")
	case config.OpSubtypeRecordDelete:
		err := clnt.DeleteAccount(config.Op.AccountID)
		if err != nil {
			return err
		}
		fmt.Printf("Account %d deleted\n", config.Op.AccountID)
	}
	return nil
}

func actNote(subop config.OpSubtype, note common.Note) error {
	clnt := client.NewClient(config.Cfg.ServerAddr,
		config.Cfg.UserName,
		config.Cfg.Password,
		config.Cfg.CacheFile,
		config.Cfg.HTTPSInsecure,
	)
	switch subop {
	case config.OpSubtypeRecordStore:
		id, err := clnt.StoreNote(config.Op.Note)
		if err != nil {
			return err
		}
		fmt.Printf("note record stored with id %d\n", id)
	case config.OpSubtypeRecordGet:
		note, err := clnt.GetNote(config.Op.NoteID)
		if err != nil {
			return err
		}
		fmt.Println(note)
	case config.OpSubtypeRecordList:
		notes, err := clnt.ListNotes()
		if err != nil {
			return err
		}
		fmt.Println(notes)
	case config.OpSubtypeRecordUpdate:
		err := clnt.UpdateNote(config.Op.NoteID, config.Op.Note)
		if err != nil {
			return err
		}
		fmt.Println("note updated")
	case config.OpSubtypeRecordDelete:
		err := clnt.DeleteNote(config.Op.NoteID)
		if err != nil {
			return err
		}
		fmt.Printf("Note %d deleted\n", config.Op.NoteID)
	}
	return nil
}

// ChooseAct performs client actions
func ChooseAct() error {
	switch config.Op.Op {
	case config.OpTypeUser:
		return actUser(config.Op.Subop, config.Op.User)
	case config.OpTypeAccount:
		return actAccount(config.Op.Subop, config.Op.Account)
	case config.OpTypeNote:
		return actNote(config.Op.Subop, config.Op.Note)
	default:
		return errors.New("unknown operation type")
	}
}
