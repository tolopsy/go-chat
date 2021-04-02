package main

import (
	"net/http"
)

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
