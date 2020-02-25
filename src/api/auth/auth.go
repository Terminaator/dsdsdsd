package auth

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Auth struct {
	Token string
}

func (a *Auth) authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Session-Token")
		if token == a.Token || "/readiness" == r.URL.Path {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}

func (a *Auth) Middleware(r *mux.Router) {
	log.Println("adding token check")
	r.Use(a.authentication)
}
