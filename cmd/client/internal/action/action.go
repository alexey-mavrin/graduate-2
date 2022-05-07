package action

import (
	"errors"
	"log"

	"github.com/alexey-mavrin/graduate-2/cmd/client/internal/config"
	"github.com/alexey-mavrin/graduate-2/internal/client"
	"github.com/alexey-mavrin/graduate-2/internal/common"
)

//go:generate go run tmpl/generator.go Account
//go:generate go run tmpl/generator.go Note
//go:generate go run tmpl/generator.go Card

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

// ChooseAct performs client actions
func ChooseAct() error {
	switch config.Op.Op {
	case config.OpTypeUser:
		return actUser(config.Op.Subop, config.Op.User)
	case config.OpTypeAccount:
		return actAccount(config.Op.Subop, config.Op.Account)
	case config.OpTypeNote:
		return actNote(config.Op.Subop, config.Op.Note)
	case config.OpTypeCard:
		return actCard(config.Op.Subop, config.Op.Card)
	default:
		return errors.New("unknown operation type")
	}
}
