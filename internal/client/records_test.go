package client

import (
	"testing"

	"github.com/alexey-mavrin/graduate-2/internal/common"
	"github.com/alexey-mavrin/graduate-2/internal/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	userName = "user1"
	userPass = "pass"
)

func Test_records(t *testing.T) {
	ts, err := newHTTPServer()
	require.NoError(t, err)
	defer ts.Close()

	clnt := NewClient(ts.URL, userName, userPass, "", false)

	_, err = clnt.RegisterUser("")
	assert.NoError(t, err)

	recType := common.NoteRecord
	record := common.Record{
		Name:   "record1",
		Type:   recType,
		Opaque: "1111",
	}

	id, err := clnt.StoreRecord(record)
	assert.NoError(t, err)

	records, err := clnt.ListRecordsType(recType)
	assert.NoError(t, err)
	expRecord := record
	expRecord.Opaque = ""

	assert.Equal(t, expRecord, records[id])

	updateRecord := record
	updateRecord.Opaque = "2222"

	err = clnt.UpdateRecordID(id, updateRecord)
	assert.NoError(t, err)

	gotRecord, err := clnt.GetRecordID(id)
	assert.NoError(t, err)

	assert.Equal(t, updateRecord, gotRecord)

	err = clnt.DeleteRecordID(id)
	assert.NoError(t, err)

	gotRecord, err = clnt.GetRecordID(id)
	assert.Error(t, err)
}

func Test_recordsCache(t *testing.T) {
	ts, err := newHTTPServer()
	require.NoError(t, err)
	defer ts.Close()

	cacheName := "cache_storage.db"
	store.DropStore(cacheName)
	clnt := NewClient(ts.URL, userName, userPass, cacheName, false)

	_, err = clnt.RegisterUser("")
	assert.NoError(t, err)

	record := common.Record{
		Name:   "record1",
		Type:   common.NoteRecord,
		Opaque: "1111",
	}

	id, err := clnt.StoreRecord(record)
	assert.NoError(t, err)

	ts.Close()

	// should get record# from cache
	gotRecord, err := clnt.GetRecordID(id)
	assert.NoError(t, err)
	assert.Equal(t, gotRecord, record)
}

func Test_recordsUpdateCache(t *testing.T) {
	ts, err := newHTTPServer()
	require.NoError(t, err)
	defer ts.Close()

	cacheName := "cache_storage.db"
	store.DropStore(cacheName)
	clnt := NewClient(ts.URL, userName, userPass, cacheName, false)

	_, err = clnt.RegisterUser("")
	assert.NoError(t, err)

	recType := common.NoteRecord
	record := common.Record{
		Name:   "record 1",
		Type:   recType,
		Opaque: "1111",
	}

	recordUpd := common.Record{
		Name:   "record 1 update",
		Type:   recType,
		Opaque: "2222",
	}

	id, err := clnt.StoreRecord(record)
	assert.NoError(t, err)

	err = clnt.UpdateRecordID(id, recordUpd)
	assert.NoError(t, err)

	ts.Close()

	// should get record# from cache
	gotRecord, err := clnt.GetRecordID(id)
	assert.NoError(t, err)
	assert.Equal(t, recordUpd, gotRecord)
}

func Test_recordsDeleteCache(t *testing.T) {
	ts, err := newHTTPServer()
	require.NoError(t, err)
	defer ts.Close()

	cacheName := "cache_storage.db"
	store.DropStore(cacheName)

	clnt := NewClient(ts.URL, userName, userPass, cacheName, false)

	_, err = clnt.RegisterUser("")
	assert.NoError(t, err)

	record := common.Record{
		Name:   "record1",
		Type:   common.NoteRecord,
		Opaque: "1111",
	}

	id, err := clnt.StoreRecord(record)
	assert.NoError(t, err)

	err = clnt.DeleteRecordID(id)
	assert.NoError(t, err)

	ts.Close()

	// should NOT get record from cache
	_, err = clnt.GetRecordID(id)
	assert.Error(t, err)
}
