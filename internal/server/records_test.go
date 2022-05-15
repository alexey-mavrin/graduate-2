package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/alexey-mavrin/graduate-2/internal/common"
	"github.com/alexey-mavrin/graduate-2/internal/store"
	"github.com/stretchr/testify/assert"
)

const (
	testUser = "user1"
	testPass = "pass1"
)

func prepareTest(t *testing.T) http.Handler {
	store.DropStore(defaultStoreFile)
	err := InitStore(defaultStoreFile)
	assert.NoError(t, err)
	router := NewRouter()
	user := common.User{
		Name:     testUser,
		Password: testPass,
	}
	createUserBody, _ := json.Marshal(user)
	resp, _ := testHTTPRequest(t,
		router,
		http.MethodPost,
		"/users",
		string(createUserBody),
		// empty user and pass
		"",
		"",
	)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	return router
}

func storeTestRecord(t *testing.T,
	router http.Handler,
	record common.Record,
) int64 {
	recordBody, _ := json.Marshal(record)

	createResp, createRespBody := testHTTPRequest(t,
		router,
		http.MethodPost,
		"/records",
		string(recordBody),
		testUser,
		testPass,
	)
	defer createResp.Body.Close()
	assert.Equal(t, http.StatusOK, createResp.StatusCode)
	var resp common.StoreRecordResponse
	err := json.Unmarshal([]byte(createRespBody), &resp)
	assert.NoError(t, err)
	return resp.ID
}

func Test_Record(t *testing.T) {
	t.Run("Store record", func(t *testing.T) {
		router := prepareTest(t)
		record := common.Record{
			Name:   "rec1",
			Opaque: "0000",
			Type:   common.NoteRecord,
		}
		_ = storeTestRecord(t, router, record)
	})

	t.Run("Update record by ID", func(t *testing.T) {
		router := prepareTest(t)
		record := common.Record{
			Name:   "rec1",
			Opaque: "0000",
			Type:   common.NoteRecord,
		}
		id := storeTestRecord(t, router, record)

		updateRecord := common.Record{
			Name:   "rec1",
			Opaque: "1111",
			Type:   common.NoteRecord,
		}
		updateRecordBody, err := json.Marshal(updateRecord)
		assert.NoError(t, err)

		updateResp, _ := testHTTPRequest(t,
			router,
			http.MethodPut,
			fmt.Sprintf("/records/%d", id),
			string(updateRecordBody),
			testUser,
			testPass,
		)
		defer updateResp.Body.Close()
		assert.Equal(t, http.StatusOK, updateResp.StatusCode)

		getResp, getRespBody := testHTTPRequest(t,
			router,
			http.MethodGet,
			fmt.Sprintf("/records/%d", id),
			"",
			testUser,
			testPass,
		)
		defer getResp.Body.Close()
		assert.Equal(t, http.StatusOK, getResp.StatusCode)

		var gotRecord common.Record
		err = json.Unmarshal([]byte(getRespBody), &gotRecord)
		assert.NoError(t, err)
		assert.Equal(t, updateRecord, gotRecord)
	})

	t.Run("Get record by ID", func(t *testing.T) {
		router := prepareTest(t)
		record := common.Record{
			Name:   "rec1",
			Opaque: "0000",
			Type:   common.NoteRecord,
		}
		id := storeTestRecord(t, router, record)
		getResp, getRespBody := testHTTPRequest(t,
			router,
			http.MethodGet,
			fmt.Sprintf("/records/%d", id),
			"",
			testUser,
			testPass,
		)
		defer getResp.Body.Close()
		assert.Equal(t, http.StatusOK, getResp.StatusCode)

		var gotRecord common.Record
		err := json.Unmarshal([]byte(getRespBody), &gotRecord)
		assert.NoError(t, err)
		assert.Equal(t, record, gotRecord)
	})

	t.Run("Get record ID by type and name", func(t *testing.T) {
		router := prepareTest(t)
		recName := "rec1"
		recType := common.NoteRecord
		record := common.Record{
			Name:   recName,
			Opaque: "0000",
			Type:   recType,
		}
		id := storeTestRecord(t, router, record)
		getResp, getRespBody := testHTTPRequest(t,
			router,
			http.MethodGet,
			fmt.Sprintf("/records/%s/%s", recType, recName),
			"",
			testUser,
			testPass,
		)
		defer getResp.Body.Close()
		assert.Equal(t, http.StatusOK, getResp.StatusCode)

		var gotResp common.StoreRecordResponse
		err := json.Unmarshal([]byte(getRespBody), &gotResp)
		assert.NoError(t, err)
		assert.Equal(t, id, gotResp.ID)
	})

	t.Run("List records", func(t *testing.T) {
		router := prepareTest(t)
		record1 := common.Record{
			Name:   "rec1",
			Opaque: "1111",
			Type:   common.NoteRecord,
		}
		id1 := storeTestRecord(t, router, record1)
		record2 := common.Record{
			Name:   "rec2",
			Opaque: "2222",
			Type:   common.NoteRecord,
		}
		id2 := storeTestRecord(t, router, record2)

		listResp, listRespBody := testHTTPRequest(t,
			router,
			http.MethodGet,
			"/records",
			"",
			testUser,
			testPass,
		)
		assert.Equal(t, http.StatusOK, listResp.StatusCode)
		var listRecords common.Records
		err := json.Unmarshal([]byte(listRespBody), &listRecords)
		assert.NoError(t, err)

		expected1 := record1
		expected1.Opaque = ""
		expected2 := record2
		expected2.Opaque = ""
		assert.Equal(t, expected1, listRecords[id1])
		assert.Equal(t, expected2, listRecords[id2])
	})

	t.Run("Delete record by ID", func(t *testing.T) {
		router := prepareTest(t)
		record := common.Record{
			Name:   "rec1",
			Opaque: "0000",
			Type:   common.NoteRecord,
		}
		id := storeTestRecord(t, router, record)

		delResp, _ := testHTTPRequest(t,
			router,
			http.MethodDelete,
			fmt.Sprintf("/records/%d", id),
			"",
			testUser,
			testPass,
		)
		defer delResp.Body.Close()
		assert.Equal(t, http.StatusOK, delResp.StatusCode)

		delResp2, _ := testHTTPRequest(t,
			router,
			http.MethodGet,
			fmt.Sprintf("/records/%d", id),
			"",
			testUser,
			testPass,
		)
		defer delResp2.Body.Close()
		assert.Equal(t, http.StatusNotFound, delResp2.StatusCode)
	})
}
