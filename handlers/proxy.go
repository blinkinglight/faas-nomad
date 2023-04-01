package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func MakeProxy(n string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %v", n, r.URL.String())
		// w.WriteHeader(http.StatusNotFound)

		realFunction := fmt.Sprintf("http://%s-function.service.consul", r.URL.Path[10:])

		req, _ := http.NewRequest("GET", realFunction, nil)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		defer resp.Body.Close()
		io.Copy(w, resp.Body)

	}
}
