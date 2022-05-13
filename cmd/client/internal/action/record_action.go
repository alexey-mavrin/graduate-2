package action

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/alexey-mavrin/graduate-2/cmd/client/internal/config"
	"github.com/alexey-mavrin/graduate-2/internal/client"
	"github.com/alexey-mavrin/graduate-2/internal/common"
	"github.com/alexey-mavrin/graduate-2/internal/crypt"
)

const defaultFileMode = 0600

func actRecord(subop config.OpSubtype, subrecord common.Opaque) error {
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
		opaque, err := subrecord.Pack()
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
	case config.OpSubtypeRecordGet:
		eRecord, err := clnt.GetRecordID(config.Op.RecordID)
		if err != nil {
			return err
		}
		record, err := crypt.DecryptRecord(*config.Key, eRecord)
		if err != nil {
			return err
		}
		fmt.Println(record)
		if config.Op.RecordType == common.BinaryRecord {
			err = writeDecodeFile(config.Op.FileName, record.Opaque)
			if err != nil {
				return err
			}
			fmt.Printf("  File %s is written\n", config.Op.FileName)
		}
	case config.OpSubtypeRecordList:
		records, err := clnt.ListRecordsType(config.Op.RecordType)
		if err != nil {
			return err
		}
		fmt.Println(records)
	case config.OpSubtypeRecordUpdate:
		record := common.Record{
			Name: config.Op.RecordName,
			Type: config.Op.RecordType,
			Meta: config.Op.RecordMeta,
		}
		opaque, err := subrecord.Pack()
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
	case config.OpSubtypeRecordDelete:
		err := clnt.DeleteRecordID(config.Op.RecordID)
		if err != nil {
			return err
		}
		fmt.Printf("Record %d deleted\n", config.Op.RecordID)
	}
	return nil
}

func writeDecodeFile(file, str string) error {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return err
	}
	err = os.WriteFile(file, data, defaultFileMode)
	if err != nil {
		return err
	}
	return nil
}
