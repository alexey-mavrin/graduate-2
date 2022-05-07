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

	var listErr, encodeErr error
	switch recordType {
	case common.AccountRecord:
		var accs common.Accounts
		accs, listErr = s.ListAccounts(user)
		if listErr != nil {
			break
		}
		encodeErr = json.NewEncoder(w).Encode(accs)
	case common.NoteRecord:
		var notes common.Notes
		notes, listErr = s.ListNotes(user)
		if listErr != nil {
			break
		}
		encodeErr = json.NewEncoder(w).Encode(notes)
	case common.CardRecord:
		var cards common.Cards
		cards, listErr = s.ListCards(user)
		if listErr != nil {
			break
		}
		encodeErr = json.NewEncoder(w).Encode(cards)
	default:
		msg := "unknown record type requested"
		log.Print(msg)
		writeStatus(w, http.StatusBadRequest, msg)
		return
	}
	if listErr != nil {
		log.Print(listErr)
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
	case common.NoteRecord:
		var note common.Note
		note, getErr = s.GetNote(user, int64(id))
		if getErr != nil {
			break
		}
		encodeErr = json.NewEncoder(w).Encode(note)
	case common.CardRecord:
		var card common.Card
		card, getErr = s.GetCard(user, int64(id))
		if getErr != nil {
			break
		}
		encodeErr = json.NewEncoder(w).Encode(card)
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
	recordType := common.RecordType(chi.URLParam(r, "record_type"))
	log.Print("deleteRecords " + recordType)

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

	var deleteErr error
	switch recordType {
	case common.AccountRecord:
		deleteErr = s.DeleteAccount(user, int64(id))
	case common.NoteRecord:
		deleteErr = s.DeleteNote(user, int64(id))
	case common.CardRecord:
		deleteErr = s.DeleteCard(user, int64(id))
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
	recordType := common.RecordType(chi.URLParam(r, "record_type"))
	log.Print("storeRecord " + recordType)

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
	var storeErr error
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
		resp.ID, storeErr = s.StoreAccount(user, account)
	case common.NoteRecord:
		var note common.Note
		err = json.Unmarshal(body, &note)
		if err != nil {
			writeStatus(w,
				http.StatusBadRequest,
				fmt.Sprintf("Cannot Parse Body: %v", err),
			)
			return
		}
		resp.Name = note.Name
		resp.ID, storeErr = s.StoreNote(user, note)
	case common.CardRecord:
		var card common.Card
		err = json.Unmarshal(body, &card)
		if err != nil {
			writeStatus(w,
				http.StatusBadRequest,
				fmt.Sprintf("Cannot Parse Body: %v", err),
			)
			return
		}
		resp.Name = card.Name
		resp.ID, storeErr = s.StoreCard(user, card)
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
	recordType := common.RecordType(chi.URLParam(r, "record_type"))
	log.Print("updateRecord " + recordType)

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
	case common.NoteRecord:
		var note common.Note
		err = json.Unmarshal(body, &note)
		if err != nil {
			writeStatus(w,
				http.StatusBadRequest,
				fmt.Sprintf("Cannot Parse Body: %v", err),
			)
			return
		}
		resp.Name = note.Name
		updateErr = s.UpdateNote(user, int64(id), note)
	case common.CardRecord:
		var card common.Card
		err = json.Unmarshal(body, &card)
		if err != nil {
			writeStatus(w,
				http.StatusBadRequest,
				fmt.Sprintf("Cannot Parse Body: %v", err),
			)
			return
		}
		resp.Name = card.Name
		updateErr = s.UpdateCard(user, int64(id), card)
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
