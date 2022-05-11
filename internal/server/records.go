package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/alexey-mavrin/graduate-2/internal/common"
	"github.com/alexey-mavrin/graduate-2/internal/store"
	"github.com/go-chi/chi/v5"
)

func listRecords(w http.ResponseWriter, r *http.Request) {
	log.Print("listRecords")

	user, _, ok := r.BasicAuth()
	if !ok {
		writeStatus(w, http.StatusBadRequest, "no basic auth")
		return
	}

	s, err := store.NewStore()
	if err != nil {
		log.Print(err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}

	records, err := s.ListRecords(user)
	if err != nil {
		log.Print(err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}
	err = json.NewEncoder(w).Encode(records)
	if err != nil {
		log.Print(err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}
}

func listRecordsType(w http.ResponseWriter, r *http.Request) {
	recordType := common.RecordType(chi.URLParam(r, "record_type"))
	log.Print("listRecords " + recordType)

	user, _, ok := r.BasicAuth()
	if !ok {
		writeStatus(w, http.StatusBadRequest, "no basic auth")
		return
	}

	s, err := store.NewStore()
	if err != nil {
		log.Print(err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}

	records, err := s.ListRecordsType(user, recordType)
	if err != nil {
		log.Print(err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}
	err = json.NewEncoder(w).Encode(records)
	if err != nil {
		log.Print(err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}
}

func getRecordID(w http.ResponseWriter, r *http.Request) {
	log.Print("getRecordID")

	user, _, ok := r.BasicAuth()
	if !ok {
		writeStatus(w, http.StatusBadRequest, "no basic auth")
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeStatus(w, http.StatusBadRequest, "cannot parse 'id' param")
		return
	}

	s, err := store.NewStore()
	if err != nil {
		log.Print(err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}

	record, err := s.GetRecordID(user, int64(id))
	if err == store.ErrNotFound {
		msg := fmt.Sprintf("Record id %d not found", id)
		log.Print(msg)
		writeStatus(w, http.StatusNotFound, msg)
		return
	}
	if err != nil {
		log.Print(err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}
	err = json.NewEncoder(w).Encode(record)

	if err != nil {
		log.Print(err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}
}

func getRecordTypeName(w http.ResponseWriter, r *http.Request) {
	log.Print("getRecordTypeName")

	user, _, ok := r.BasicAuth()
	if !ok {
		writeStatus(w, http.StatusBadRequest, "no basic auth")
		return
	}

	recordType := common.RecordType(chi.URLParam(r, "record_type"))
	recordName := chi.URLParam(r, "record_name")

	s, err := store.NewStore()
	if err != nil {
		log.Print(err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}

	record, err := s.GetRecordTypeName(user, recordType, recordName)
	if err == store.ErrNotFound {
		msg := fmt.Sprintf("Record %s of type %s not found",
			recordName, recordType)
		log.Print(msg)
		writeStatus(w, http.StatusNotFound, msg)
		return
	}
	if err != nil {
		log.Print(err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}
	err = json.NewEncoder(w).Encode(record)

	if err != nil {
		log.Print(err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}
}

func deleteRecordID(w http.ResponseWriter, r *http.Request) {
	log.Print("deleteRecordID")

	user, _, ok := r.BasicAuth()
	if !ok {
		writeStatus(w, http.StatusBadRequest, "no basic auth")
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeStatus(w, http.StatusBadRequest, "cannot parse 'id' param")
		return
	}

	s, err := store.NewStore()
	if err != nil {
		log.Print(err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}

	err = s.DeleteRecordID(user, int64(id))
	if err == store.ErrNotFound {
		msg := fmt.Sprintf("Record id %d not found", id)
		log.Print(msg)
		writeStatus(w, http.StatusNotFound, msg)
		return
	}
	if err != nil {
		log.Print(err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}

	writeStatus(w, http.StatusOK, "OK")
}

func deleteRecordTypeName(w http.ResponseWriter, r *http.Request) {
	log.Print("deleteRecordTypeName")

	user, _, ok := r.BasicAuth()
	if !ok {
		writeStatus(w, http.StatusBadRequest, "no basic auth")
		return
	}

	recordType := common.RecordType(chi.URLParam(r, "record_type"))
	recordName := chi.URLParam(r, "record_name")

	s, err := store.NewStore()
	if err != nil {
		log.Print(err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}

	err = s.DeleteRecordTypeName(user, recordType, recordName)
	if err == store.ErrNotFound {
		msg := fmt.Sprintf("Record %s of type %s not found",
			recordName, recordType)
		log.Print(msg)
		writeStatus(w, http.StatusNotFound, msg)
		return
	}
	if err != nil {
		log.Print(err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}

	writeStatus(w, http.StatusOK, "OK")
}

func storeRecord(w http.ResponseWriter, r *http.Request) {
	log.Print("storeRecord")

	user, _, ok := r.BasicAuth()
	if !ok {
		writeStatus(w, http.StatusBadRequest, "no basic auth")
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}

	s, err := store.NewStore()
	if err != nil {
		log.Print(err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}

	var resp common.StoreRecordResponse
	resp.Status = "OK"

	var record common.Record
	err = json.Unmarshal(body, &record)
	if err != nil {
		writeStatus(w,
			http.StatusBadRequest,
			fmt.Sprintf("Cannot Parse Body: %v", err),
		)
		return
	}
	resp.Name = record.Name
	resp.ID, err = s.StoreRecord(user, record)

	if err != nil {
		log.Printf("storeRecord() error: %v", err)
		resp.Status = "error"
		writeStatus(w,
			http.StatusInternalServerError,
			fmt.Sprintf("Cannot Store Record: %v", err),
		)
		return
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}
}

func updateRecordID(w http.ResponseWriter, r *http.Request) {
	log.Print("updateRecordID")

	user, _, ok := r.BasicAuth()
	if !ok {
		writeStatus(w, http.StatusBadRequest, "no basic auth")
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeStatus(w, http.StatusBadRequest, "cannot parse 'id' param")
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}

	s, err := store.NewStore()
	if err != nil {
		log.Print(err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}

	var resp common.StoreRecordResponse
	resp.Status = "OK"
	var record common.Record
	err = json.Unmarshal(body, &record)
	if err != nil {
		writeStatus(w,
			http.StatusBadRequest,
			fmt.Sprintf("Cannot Parse Body: %v", err),
		)
		return
	}
	resp.Name = record.Name
	err = s.UpdateRecordID(user, int64(id), record)

	if err != nil {
		log.Printf("update record error: %v", err)
		resp.Status = "error"
		writeStatus(w,
			http.StatusInternalServerError,
			fmt.Sprintf("Cannot Update Record: %v", err),
		)
		return
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}
}

func updateRecordTypeName(w http.ResponseWriter, r *http.Request) {
	log.Print("updateRecordTypeName")

	user, _, ok := r.BasicAuth()
	if !ok {
		writeStatus(w, http.StatusBadRequest, "no basic auth")
		return
	}

	recordType := common.RecordType(chi.URLParam(r, "record_type"))
	recordName := chi.URLParam(r, "record_name")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}

	s, err := store.NewStore()
	if err != nil {
		log.Print(err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}

	var resp common.StoreRecordResponse
	resp.Status = "OK"
	var record common.Record
	err = json.Unmarshal(body, &record)
	if err != nil {
		writeStatus(w,
			http.StatusBadRequest,
			fmt.Sprintf("Cannot Parse Body: %v", err),
		)
		return
	}
	resp.Name = record.Name
	err = s.UpdateRecordTypeName(user, recordType, recordName, record)

	if err != nil {
		log.Printf("update record error: %v", err)
		resp.Status = "error"
		writeStatus(w,
			http.StatusInternalServerError,
			fmt.Sprintf("Cannot Update Record: %v", err),
		)
		return
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}
}
