package store

import (
	"reflect"
	"testing"

	"github.com/alexey-mavrin/graduate-2/internal/common"
	"github.com/stretchr/testify/assert"
)

func TestStore_GetRecord(t *testing.T) {
	store := dropCreateStore(t)
	t.Run("Get single record#", func(t *testing.T) {
		user := "user1"

		recType := common.NoteRecord
		recName := "rec1"
		record := common.Record{
			Name:   recName,
			Opaque: "1111",
			Type:   recType,
			Meta:   "abc",
		}

		_, err := store.AddUser(common.User{
			Name: user,
		})
		assert.NoError(t, err)

		id, err := store.StoreRecord(user, record)
		assert.NoError(t, err)

		recordRet, err := store.GetRecordByID(user, id)
		assert.NoError(t, err)

		assert.Equal(t, record, recordRet)

		recordRet, err = store.GetRecordByTypeName(user, recType, recName)
		assert.NoError(t, err)

		assert.Equal(t, record, recordRet)
	})
}

func TestStore_ListRecords(t *testing.T) {
	store := dropCreateStore(t)
	t.Run("Get multiple records", func(t *testing.T) {
		user := "user1"
		name1 := "record 1"
		name2 := "record 2"
		recType := common.NoteRecord
		opaque1 := "1111"
		opaque2 := "2222"

		_, err := store.AddUser(common.User{
			Name: user,
		})
		assert.NoError(t, err)

		id1, err := store.StoreRecord(user, common.Record{
			Name:   name1,
			Type:   recType,
			Opaque: opaque1,
		})
		assert.NoError(t, err)

		id2, err := store.StoreRecord(user, common.Record{
			Name:   name2,
			Type:   recType,
			Opaque: opaque2,
		})
		assert.NoError(t, err)

		records, err := store.ListRecords(user)
		assert.NoError(t, err)

		wantRecords := make(common.Records)
		wantRecords[id1] = common.Record{
			Name: name1,
			Type: recType,
		}
		wantRecords[id2] = common.Record{
			Name: name2,
			Type: recType,
		}
		assert.True(t, reflect.DeepEqual(records, wantRecords))
	})
}

func TestStore_DeleteRecord(t *testing.T) {
	store := dropCreateStore(t)
	t.Run("Delete record record", func(t *testing.T) {
		user := "user1"

		recName := "rec1"
		recType := common.NoteRecord
		record := common.Record{
			Name:   recName,
			Type:   recType,
			Opaque: "1111",
		}

		_, err := store.AddUser(common.User{
			Name: user,
		})
		assert.NoError(t, err)

		id, err := store.StoreRecord(user, record)
		assert.NoError(t, err)

		err = store.DeleteRecordByID(user, id)
		assert.NoError(t, err)

		// same attempt should result in error
		err = store.DeleteRecordByID(user, id)
		assert.Error(t, err)

		// attempt to delete non-existing record# recourd
		// should result in error
		err = store.DeleteRecordByID(user, 999999)
		assert.Error(t, err)

		_, err = store.StoreRecord(user, record)
		assert.NoError(t, err)

		err = store.DeleteRecordByTypeName(user, recType, recName)
		assert.NoError(t, err)

		// same attempt should result in error
		err = store.DeleteRecordByTypeName(user, recType, recName)
		assert.Error(t, err)

	})
}

func TestStore_UpdateRecord(t *testing.T) {
	store := dropCreateStore(t)
	t.Run("Update record record", func(t *testing.T) {
		user := "user1"
		name := "record"
		recType := common.NoteRecord
		opaque1 := "1111"
		opaque2 := "2222"
		opaque3 := "3333"

		_, err := store.AddUser(common.User{
			Name: user,
		})
		assert.NoError(t, err)

		id, err := store.StoreRecord(user, common.Record{
			Name:   name,
			Type:   recType,
			Opaque: opaque1,
		})
		assert.NoError(t, err)

		err = store.UpdateRecordByID(user, id, common.Record{
			Name:   name,
			Type:   recType,
			Opaque: opaque2,
		})
		assert.NoError(t, err)

		record, err := store.GetRecordByID(user, id)
		assert.NoError(t, err)
		assert.Equal(t, opaque2, record.Opaque,
			"Updated record number should change")

		err = store.UpdateRecordByTypeName(user, recType, name,
			common.Record{
				Name:   name,
				Type:   recType,
				Opaque: opaque3,
			})
		assert.NoError(t, err)

		record, err = store.GetRecordByID(user, id)
		assert.NoError(t, err)
		assert.Equal(t, opaque3, record.Opaque,
			"Updated record number should change (again)")
	})
}
