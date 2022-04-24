package server

import (
	"log"
	"net/http"

	"github.com/alexey-mavrin/graduate-2/internal/store"
)

func verifyUser(user string, pass string) (bool, error) {
	s, err := store.NewStore()
	if err != nil {
		return false, err
	}

	ok, err := s.CheckUserAuth(user, pass)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func authUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, basicSet := r.BasicAuth()

		if basicSet {
			userOK, err := verifyUser(user, pass)
			if err != nil {
				log.Printf("verifyUser error: %v", err)
				return
			}
			if !userOK {
				http.Error(w,
					"Access Denied",
					http.StatusForbidden,
				)
				return
			}
			next.ServeHTTP(w, r)
			return
		}

		if !basicSet {
			if r.URL.Path == registerPath {
				// registering does not require user auth
				next.ServeHTTP(w, r)
				return
			}
			w.Header().Set("WWW-Authenticate", `Basic realm="storeapi"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

	})
}
