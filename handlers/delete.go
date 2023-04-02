package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/alexellis/faas/gateway/requests"
	"github.com/hashicorp/nomad/api"
)

// MakeReader implements the OpenFaaS reader handler
func MakeDelete(n string, client *api.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %v", n, r.URL.Path)
		defer r.Body.Close()

		body, _ := ioutil.ReadAll(r.Body)

		request := requests.DeleteFunctionRequest{}
		err := json.Unmarshal(body, &request)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		name := fmt.Sprintf("OpenFaaS-%s", request.FunctionName)
		client.Jobs().Deregister(name, true, nil)

	}
}
