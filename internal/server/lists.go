package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/alexey-mavrin/graduate-2/internal/store"
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

	accs, err := s.ListAccounts(user)
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
