package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/alexey-mavrin/graduate-2/internal/store"
)

// AddUserResponse is the responce for AddUser request
type AddUserResponse struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	ID     int64  `json:"id"`
}

func createUser(w http.ResponseWriter, r *http.Request) {
	log.Print("createUser")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}

	var user store.User
	err = json.Unmarshal(body, &user)
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
	var resp AddUserResponse
	resp.Name = user.Name
	resp.Status = "OK"
	resp.ID, err = s.AddUser(user)
	if err != nil {
		log.Printf("AddUser() error: %v", err)
		if errors.Is(err, store.ErrAlreadyExists) {
			resp.Status = "already exists"
			writeStatus(w,
				http.StatusFound,
				"User Already Exists",
			)
		} else {
			resp.Status = "error"
			writeStatus(w,
				http.StatusBadRequest,
				"Cannot Add User",
			)
		}
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Printf("cannot encode AddUserResponse: %v", err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}
}
