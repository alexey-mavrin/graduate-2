package config

import (
	"errors"
	"flag"
	"os"

	"github.com/alexey-mavrin/graduate-2/internal/common"
)

type (
	// OpType is the operation type
	OpType int
	// OpSubtype is the operation subtype
	OpSubtype int
)

const (
	// OpTypeUser is for user operations
	OpTypeUser OpType = iota
	// OpTypeAccount is for account operations
	OpTypeAccount
	// OpTypeNote is for note operations
	OpTypeNote
	// OpTypeCard is for note operations
	OpTypeCard
)

const (
	// OpSubtypeUserRegister is the user registration
	OpSubtypeUserRegister OpSubtype = iota
	// OpSubtypeUserVerify is the user auth verification
	OpSubtypeUserVerify

	// OpSubtypeRecordStore is the account record creation
	OpSubtypeRecordStore
	// OpSubtypeRecordGet is the account regord retrieval
	OpSubtypeRecordGet
	// OpSubtypeRecordList is the listing of account records
	OpSubtypeRecordList
	// OpSubtypeRecordUpdate is the account record update
	OpSubtypeRecordUpdate
	// OpSubtypeRecordDelete is the removal of the account record
	OpSubtypeRecordDelete
)

// Operation describes the current operation type
type Operation struct {
	Op        OpType
	Subop     OpSubtype
	User      common.User
	Account   common.Account
	AccountID int64
	Note      common.Note
	NoteID    int64
	Card      common.Card
	CardID    int64
}

// Op describes the current operation
var Op Operation

// ParseFlags parses cmd line arguments
func ParseFlags() error {
	userFlags := flag.NewFlagSet("user", flag.ExitOnError)
	accFlags := flag.NewFlagSet("acc", flag.ExitOnError)
	noteFlags := flag.NewFlagSet("note", flag.ExitOnError)
	cardFlags := flag.NewFlagSet("card", flag.ExitOnError)

	userAction := userFlags.String("a", "verify", "action: verify|register")

	accAction := accFlags.String("a",
		"list",
		"action: list|store|get|update|delete",
	)
	accName := accFlags.String("n", "", "account name")
	accUserName := accFlags.String("u", "", "account user name")
	accPassword := accFlags.String("p", "", "account password")
	accURL := accFlags.String("l", "", "account URL")
	accMeta := accFlags.String("m", "", "account metainfo")
	accID := accFlags.Int64("i", 0, "account ID")

	noteAction := noteFlags.String("a",
		"list",
		"action: list|store|get|update|delete",
	)
	noteName := noteFlags.String("n", "", "note name")
	noteText := noteFlags.String("t", "", "note text")
	noteMeta := noteFlags.String("m", "", "note metainfo")
	noteID := noteFlags.Int64("i", 0, "note ID")

	cardAction := cardFlags.String("a",
		"list",
		"action: list|store|get|update|delete",
	)
	cardName := cardFlags.String("n", "", "card name")
	cardHolder := cardFlags.String("ch", "", "card holder")
	cardNumber := cardFlags.String("num", "", "card number")
	cardExpMonth := cardFlags.Int("em", 0, "card expiry month")
	cardExpYear := cardFlags.Int("ey", 0, "card expiry year")
	cardCVC := cardFlags.String("c", "", "card CVC code")
	cardMeta := cardFlags.String("m", "", "card metainfo")
	cardID := cardFlags.Int64("i", 0, "card ID")

	if len(os.Args) < 2 {
		return errors.New("mode is not set")
	}

	switch os.Args[1] {
	case "user":
		userFlags.Parse(os.Args[2:])
	case "acc":
		accFlags.Parse(os.Args[2:])
	case "note":
		noteFlags.Parse(os.Args[2:])
	case "card":
		cardFlags.Parse(os.Args[2:])
	default:
		return errors.New("unknown mode")
	}

	if userFlags.Parsed() {
		Op.Op = OpTypeUser
		switch *userAction {
		case "verify":
			Op.Subop = OpSubtypeUserVerify
		case "register":
			Op.Subop = OpSubtypeUserRegister
		default:
			return errors.New("unknown user action")
		}
	} else if accFlags.Parsed() {
		Op.Op = OpTypeAccount
		switch *accAction {
		case "store":
			Op.Subop = OpSubtypeRecordStore
		case "get":
			Op.Subop = OpSubtypeRecordGet
		case "list":
			Op.Subop = OpSubtypeRecordList
		case "update":
			Op.Subop = OpSubtypeRecordUpdate
		case "delete":
			Op.Subop = OpSubtypeRecordDelete
		}

		Op.Account.Name = *accName
		Op.Account.UserName = *accUserName
		Op.Account.Password = *accPassword
		Op.Account.URL = *accURL
		Op.Account.Meta = *accMeta
		Op.AccountID = *accID
	} else if noteFlags.Parsed() {
		Op.Op = OpTypeNote
		switch *noteAction {
		case "store":
			Op.Subop = OpSubtypeRecordStore
		case "get":
			Op.Subop = OpSubtypeRecordGet
		case "list":
			Op.Subop = OpSubtypeRecordList
		case "update":
			Op.Subop = OpSubtypeRecordUpdate
		case "delete":
			Op.Subop = OpSubtypeRecordDelete
		}

		Op.Note.Name = *noteName
		Op.Note.Text = *noteText
		Op.Note.Meta = *noteMeta
		Op.NoteID = *noteID
	} else if cardFlags.Parsed() {
		Op.Op = OpTypeCard
		switch *cardAction {
		case "store":
			Op.Subop = OpSubtypeRecordStore
		case "get":
			Op.Subop = OpSubtypeRecordGet
		case "list":
			Op.Subop = OpSubtypeRecordList
		case "update":
			Op.Subop = OpSubtypeRecordUpdate
		case "delete":
			Op.Subop = OpSubtypeRecordDelete
		}

		Op.Card.Name = *cardName
		Op.Card.Holder = *cardHolder
		Op.Card.Number = *cardNumber
		Op.Card.ExpMonth = *cardExpMonth
		Op.Card.ExpYear = *cardExpYear
		Op.Card.CVC = *cardCVC
		Op.Card.Meta = *cardMeta
		Op.CardID = *cardID
	}

	return nil
}
