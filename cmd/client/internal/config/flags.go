package config

import (
	"encoding/base64"
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
	// OpTypeBinary is for binary operations
	OpTypeBinary
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
	// OpSubtypeOther is unknown operation
	OpSubtypeOther
)

// Operation describes the current operation type
type Operation struct {
	Op         OpType
	Subop      OpSubtype
	User       common.User
	Account    common.Account
	Note       common.Note
	Card       common.Card
	Binary     common.Binary
	RecordID   int64
	RecordName string
	RecordMeta string
	RecordType common.RecordType
	FileName   string
}

// Op describes the current operation
var Op Operation

func actionType(a *string) OpSubtype {
	switch *a {
	case "store":
		return OpSubtypeRecordStore
	case "get":
		return OpSubtypeRecordGet
	case "list":
		return OpSubtypeRecordList
	case "update":
		return OpSubtypeRecordUpdate
	case "delete":
		return OpSubtypeRecordDelete
	}
	return OpSubtypeOther
}

// ParseFlags parses cmd line arguments
func ParseFlags() error {
	userFlags := flag.NewFlagSet("user", flag.ExitOnError)
	accFlags := flag.NewFlagSet("acc", flag.ExitOnError)
	noteFlags := flag.NewFlagSet("note", flag.ExitOnError)
	cardFlags := flag.NewFlagSet("card", flag.ExitOnError)
	binFlags := flag.NewFlagSet("bin", flag.ExitOnError)

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

	binAction := binFlags.String("a",
		"list",
		"action: list|store|get|update|delete",
	)
	binName := binFlags.String("n", "", "binary record name")
	binFile := binFlags.String("f", "", "file name")
	binID := binFlags.Int64("i", 0, "binary record ID")

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
	case "bin":
		binFlags.Parse(os.Args[2:])
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
		Op.RecordType = common.AccountRecord
		Op.Subop = actionType(accAction)

		Op.RecordName = *accName
		Op.Account.UserName = *accUserName
		Op.Account.Password = *accPassword
		Op.Account.URL = *accURL
		Op.RecordMeta = *accMeta
		Op.RecordID = *accID
	} else if noteFlags.Parsed() {
		Op.Op = OpTypeNote
		Op.RecordType = common.NoteRecord
		Op.Subop = actionType(noteAction)

		Op.RecordName = *noteName
		Op.Note.Text = *noteText
		Op.RecordMeta = *noteMeta
		Op.RecordID = *noteID
	} else if cardFlags.Parsed() {
		Op.Op = OpTypeCard
		Op.RecordType = common.CardRecord
		Op.Subop = actionType(cardAction)

		Op.RecordName = *cardName
		Op.Card.Holder = *cardHolder
		Op.Card.Number = *cardNumber
		Op.Card.ExpMonth = *cardExpMonth
		Op.Card.ExpYear = *cardExpYear
		Op.Card.CVC = *cardCVC
		Op.RecordMeta = *cardMeta
		Op.RecordID = *cardID
	} else if binFlags.Parsed() {
		Op.Op = OpTypeBinary
		Op.RecordType = common.BinaryRecord
		Op.Subop = actionType(binAction)

		Op.RecordName = *binName
		var err error
		if Op.Subop == OpSubtypeRecordStore || Op.Subop == OpSubtypeRecordUpdate {
			Op.Binary.Data, err = readEncodeFile(*binFile)
		}
		Op.FileName = *binFile
		if err != nil {
			return err
		}
		Op.RecordID = *binID
	}

	return nil
}

func readEncodeFile(file string) (string, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	str := base64.StdEncoding.EncodeToString(data)
	return str, nil
}
