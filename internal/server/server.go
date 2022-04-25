package server

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/alexey-mavrin/graduate-2/internal/store"
	"github.com/go-chi/chi/v5"
)

const (
	// ListenAddress  is the address the server listen to
	ListenAddress = ":8080"
	// registerPath is the path to serve requests to register new users
	registerPath = "/users/"
)

func writeStatus(w http.ResponseWriter, code int, status string) {
	w.WriteHeader(code)
	w.Write([]byte(`{"Status":"` + status + `"}`))
}

func checkSetContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		cType := r.Header.Get("Content-Type")
		if cType != "application/json" {
			log.Print("checkContentType: bad content type " + cType)
			writeStatus(w, http.StatusBadRequest, "Bad Content Type")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// StartServer starts the server
func StartServer() error {
	r := NewRouter()
	c := make(chan error)
	go func() {
		log.Printf("Listening on %v...", ListenAddress)
		err := http.ListenAndServe(ListenAddress, r)
		c <- err
	}()

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
	)

	select {
	case sig := <-signalChannel:
		switch sig {
		case os.Interrupt:
			log.Print("sigint")
		case syscall.SIGTERM:
			log.Print("sigterm")
		case syscall.SIGINT:
			log.Print("sigint")
		case syscall.SIGQUIT:
			log.Print("sigquit")
		}
	case err := <-c:
		log.Print(err)
		return err
	}

	log.Print("Server finished")
	return store.CloseDB()
}

// NewRouter returns new Router
func NewRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(checkSetContentType)
	r.Use(authUser)

	r.Post("/users/", createUser)
	r.Get("/ping/", pingHandler)
	// r.Post("/accounts/", storeAccount)
	// r.Get("/accounts/", listAccounts)
	// r.Get("/accounts/{id}", listAccount)
	// r.Put("/accounts/{id}", updateAccount)
	// r.Delete("/accounts/{id}", deleteAccount)

	return r
}
