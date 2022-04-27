package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/alexey-mavrin/graduate-2/internal/store"
)

// StoreAccountResponse is the responce for store account
type StoreAccountResponse struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	ID     int64  `json:"id"`
}

func storeAccount(w http.ResponseWriter, r *http.Request) {
	log.Print("storeAccount")
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

	var account store.Account
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

	var resp StoreAccountResponse
	resp.Name = account.Name
	resp.Status = "OK"
	resp.ID, err = s.StoreAccount(user, account)
	if err != nil {
		log.Printf("StoreAccount() error: %v", err)
		resp.Status = "error"
		writeStatus(w,
			http.StatusBadRequest,
			fmt.Sprintf("Cannot Store Account: %v", err),
		)
		return
	}
}
