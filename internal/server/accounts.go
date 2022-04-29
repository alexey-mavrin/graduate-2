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

func listAccounts(w http.ResponseWriter, r *http.Request) {
	log.Print("listAccounts")

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

	accs, err := s.GetAccounts(user)
	if err != nil {
		log.Print(err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}

	if err := json.NewEncoder(w).Encode(accs); err != nil {
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}
}

func getDeleteAccount(w http.ResponseWriter, r *http.Request) {
	log.Print("getDeleteAccount")

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

	var acc common.Account

	switch r.Method {
	case http.MethodGet:
		acc, err = s.GetAccount(user, int64(id))
	case http.MethodDelete:
		err = s.DeleteAccount(user, int64(id))
	}
	if err == store.ErrNotFound {
		log.Print("Account not found")
		writeStatus(w,
			http.StatusNotFound,
			"Account not found",
		)
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

	if r.Method == http.MethodDelete {
		writeStatus(w,
			http.StatusOK,
			"OK",
		)
		return
	}
	if err := json.NewEncoder(w).Encode(acc); err != nil {
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}
}

func storeUpdateAccount(w http.ResponseWriter, r *http.Request) {
	log.Print("storeUpdateAccount")

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

	var account common.Account
	err = json.Unmarshal(body, &account)
	if err != nil {
		writeStatus(w,
			http.StatusBadRequest,
			fmt.Sprintf("Cannot Parse Body: %v", err),
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

	var resp common.StoreAccountResponse
	resp.Name = account.Name
	resp.Status = "OK"
	switch r.Method {
	case http.MethodPost:
		resp.ID, err = s.StoreAccount(user, account)
	case http.MethodPut:
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			writeStatus(w, http.StatusBadRequest, "cannot parse 'id' param")
			return
		}
		err = s.UpdateAccount(user, int64(id), account)
	}

	if err != nil {
		log.Printf("storeUpdateAccount() error: %v", err)
		resp.Status = "error"
		writeStatus(w,
			http.StatusBadRequest,
			fmt.Sprintf("Cannot Store or Update Account: %v", err),
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
