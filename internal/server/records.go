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
	recordType := common.RecordType(chi.URLParam(r, "record_type"))

	switch recordType {
	case common.AccountRecord:
		listAccounts(w, r)
		return
	default:
		msg := "unknown record type requested"
		log.Print(msg)
		writeStatus(w,
			http.StatusBadRequest,
			msg,
		)
		return
	}

}

func getRecord(w http.ResponseWriter, r *http.Request) {
	log.Print("getRecord")

	user, _, ok := r.BasicAuth()
	if !ok {
		writeStatus(w, http.StatusBadRequest, "no basic auth")
		return
	}

	recordType := common.RecordType(chi.URLParam(r, "record_type"))
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

	var getErr, encodeErr error
	switch recordType {
	case common.AccountRecord:
		var acc common.Account
		acc, getErr = s.GetAccount(user, int64(id))
		if getErr != nil {
			break
		}
		encodeErr = json.NewEncoder(w).Encode(acc)
	default:
		msg := "unknown record type requested"
		log.Print(msg)
		writeStatus(w,
			http.StatusBadRequest,
			msg,
		)
		return
	}
	if getErr == store.ErrNotFound {
		msg := "Record not found"
		log.Print(msg)
		writeStatus(w, http.StatusNotFound, msg)
		return
	}
	if getErr != nil {
		log.Print(err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}
	if encodeErr != nil {
		log.Print(encodeErr)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}
}

func deleteRecord(w http.ResponseWriter, r *http.Request) {
	log.Print("deleteRecords")

	user, _, ok := r.BasicAuth()
	if !ok {
		writeStatus(w, http.StatusBadRequest, "no basic auth")
		return
	}

	recordType := common.RecordType(chi.URLParam(r, "record_type"))
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

	var deleteErr error
	switch recordType {
	case common.AccountRecord:
		deleteErr = s.DeleteAccount(user, int64(id))
	default:
		msg := "unknown record type requested"
		log.Print(msg)
		writeStatus(w,
			http.StatusBadRequest,
			msg,
		)
		return
	}
	if deleteErr == store.ErrNotFound {
		msg := "Record not found"
		log.Print(msg)
		writeStatus(w, http.StatusNotFound, msg)
		return
	}
	if deleteErr != nil {
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

	recordType := common.RecordType(chi.URLParam(r, "record_type"))

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
	var storeErr error
	switch recordType {
	case common.AccountRecord:
		var account common.Account
		resp.Name = account.Name
		err = json.Unmarshal(body, &account)
		if err != nil {
			writeStatus(w,
				http.StatusBadRequest,
				fmt.Sprintf("Cannot Parse Body: %v", err),
			)
			return
		}
		resp.ID, storeErr = s.StoreAccount(user, account)
	default:
		msg := "unknown record type requested"
		log.Print(msg)
		writeStatus(w,
			http.StatusBadRequest,
			msg,
		)
		return
	}

	if storeErr != nil {
		log.Printf("storeAccount() error: %v", err)
		resp.Status = "error"
		writeStatus(w,
			http.StatusInternalServerError,
			fmt.Sprintf("Cannot Store Account: %v", err),
		)
		return
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}
}

func updateRecord(w http.ResponseWriter, r *http.Request) {
	log.Print("updateRecord")

	user, _, ok := r.BasicAuth()
	if !ok {
		writeStatus(w, http.StatusBadRequest, "no basic auth")
		return
	}

	recordType := common.RecordType(chi.URLParam(r, "record_type"))

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
	var updateErr error
	switch recordType {
	case common.AccountRecord:
		var account common.Account
		err = json.Unmarshal(body, &account)
		if err != nil {
			writeStatus(w,
				http.StatusBadRequest,
				fmt.Sprintf("Cannot Parse Body: %v", err),
			)
			return
		}
		resp.Name = account.Name
		updateErr = s.UpdateAccount(user, int64(id), account)
	default:
		msg := "unknown record type requested"
		log.Print(msg)
		writeStatus(w,
			http.StatusBadRequest,
			msg,
		)
		return
	}

	if updateErr != nil {
		log.Printf("updateAccount() error: %v", err)
		resp.Status = "error"
		writeStatus(w,
			http.StatusInternalServerError,
			fmt.Sprintf("Cannot Update Account: %v", err),
		)
		return
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}
}
