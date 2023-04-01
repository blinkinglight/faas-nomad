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

		realFunction := fmt.Sprintf("http://%s-function.service.consul:8080", r.URL.Path[10:])
		log.Printf("%s", realFunction)
		req, _ := http.NewRequest("GET", realFunction, nil)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("%v", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		defer resp.Body.Close()
		io.Copy(w, resp.Body)

	}
}
