package action

import (
	"encoding/json"
	"fmt"

	"github.com/alexey-mavrin/graduate-2/cmd/client/internal/config"
	"github.com/alexey-mavrin/graduate-2/internal/client"
	"github.com/alexey-mavrin/graduate-2/internal/common"
	"github.com/alexey-mavrin/graduate-2/internal/crypt"
)

func actCard(subop config.OpSubtype, record common.Card) error {
	clnt := client.NewClient(config.Cfg.ServerAddr,
		config.Cfg.UserName,
		config.Cfg.Password,
		config.Cfg.CacheFile,
		config.Cfg.HTTPSInsecure,
	)
	switch subop {
	case config.OpSubtypeRecordStore:
		record := common.Record{
			Name: config.Op.RecordName,
			Type: config.Op.RecordType,
			Meta: config.Op.RecordMeta,
		}
		opaque, err := json.Marshal(config.Op.Card)
		if err != nil {
			return err
		}
		record.Opaque = string(opaque)

		eRecord, err := crypt.EncryptRecord(*config.Key, record)
		if err != nil {
			return err
		}
		id, err := clnt.StoreRecord(eRecord)
		if err != nil {
			return err
		}
		fmt.Printf("record stored with id %d\n", id)
	case config.OpSubtypeRecordGetID:
		eRecord, err := clnt.GetRecordID(config.Op.RecordID)
		if err != nil {
			return err
		}
		record, err := crypt.DecryptRecord(*config.Key, eRecord)
		if err != nil {
			return err
		}
		fmt.Println(record)
	case config.OpSubtypeRecordListType:
		records, err := clnt.ListRecordsType(config.Op.RecordType)
		if err != nil {
			return err
		}
		fmt.Println(records)
	case config.OpSubtypeRecordUpdateID:
		record := common.Record{
			Name: config.Op.RecordName,
			Type: config.Op.RecordType,
			Meta: config.Op.RecordMeta,
		}
		opaque, err := json.Marshal(config.Op.Card)
		if err != nil {
			return err
		}
		record.Opaque = string(opaque)

		eRecord, err := crypt.EncryptRecord(*config.Key, record)
		if err != nil {
			return err
		}
		err = clnt.UpdateRecordID(config.Op.RecordID, eRecord)
		if err != nil {
			return err
		}
		fmt.Println("record updated")
	case config.OpSubtypeRecordDeleteID:
		err := clnt.DeleteRecordID(config.Op.RecordID)
		if err != nil {
			return err
		}
		fmt.Printf("Record %d deleted\n", config.Op.RecordID)
	}
	return nil
}
