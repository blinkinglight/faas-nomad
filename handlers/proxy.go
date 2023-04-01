package handlers

import (
	"log"
	"net/http"
)

func MakeProxy(n string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %v", n, r.URL.String())
		w.WriteHeader(http.StatusNotFound)

	}
}
