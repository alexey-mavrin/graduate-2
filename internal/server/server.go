package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/alexey-mavrin/graduate-2/internal/store"
	"github.com/go-chi/chi/v5"
)

const (
	defaultListenAddress = ":8080"
	defaultStoreFile     = "server_store.db"
	// registerPath is the path to serve requests to register new users
	registerPath = "/users"
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
func StartServer(listenPort int, storeFile string) error {
	store.DBFile = storeFile
	if storeFile == "" {
		store.DBFile = defaultStoreFile
	}

	listenAddress := fmt.Sprintf(":%d", listenPort)
	if listenPort == 0 {
		listenAddress = defaultListenAddress
	}

	r := NewRouter()
	c := make(chan error)
	go func() {
		log.Printf("Listening on %v...", listenAddress)
		err := http.ListenAndServe(listenAddress, r)
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

	r.Post("/users", createUser)
	r.Get("/ping", pingHandler)
	r.Post("/records/{record_type}", storeRecord)
	r.Get("/records/{record_type}", listRecords)
	r.Get("/records/{record_type}/{id}", getRecord)
	r.Put("/records/{record_type}/{id}", updateRecord)
	r.Delete("/records/{record_type}/{id}", deleteRecord)

	return r
}
