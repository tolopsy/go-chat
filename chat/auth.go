package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// authHandler checks if a user is authenticated or not
type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("auth")
	if err == http.ErrNoCookie {
		//redirect to /login if "auth" cookie is not present

		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	if err != nil {
		// in case of other error that's not http.ErrNoCookie

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// "auth" Cookie present
	h.next.ServeHTTP(w, r)
}

func MustAuth(h http.Handler) http.Handler {
	return &authHandler{next: h}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	urlSegments := strings.Split(r.URL.Path, "/")
	action := urlSegments[2]
	provider := urlSegments[3]
	switch action {
	case "login":
		log.Println("Login with", provider)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Auth action %s is not supporterd", action)
	}
}
