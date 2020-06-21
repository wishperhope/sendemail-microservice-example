package main

import (
	"log"
	"net/http"
)

// Just static token setting, implement another auth system if nescessary
func (s *Server) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if err := r.ParseForm(); err != nil {
			log.Println("Cannot Extract Form Data")
			http.Error(w, "Bad Request", 400)
			return
		}

		if err := r.Header.Get("Authorization"); err == "" {
			log.Println("Empty Token")
			http.Error(w, "Bad Request", 400)
			return
		}

		if r.Header.Get("Authorization") != s.token {
			log.Println("Wrong Token : ", r.Header.Get("Authorization"))
			http.Error(w, "Bad Request", 400)
			return
		}

		next(w, r)
	}
}
