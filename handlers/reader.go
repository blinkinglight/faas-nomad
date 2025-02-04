package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/alexellis/faas/gateway/requests"
	"github.com/blinkinglight/faas-nomad/nomad"
	"github.com/hashicorp/nomad/api"
)

func getFunctions(service string,
	client nomad.Allocations,
	allocs []*api.AllocationListStub) ([]requests.Function, error) {

	functions := make([]requests.Function, 0)
	for _, a := range allocs {

		if a.ClientStatus == "running" {
			allocation, _, err := client.Info(a.ID, nil)
			if err != nil {
				return functions, err
			}
			name := *allocation.Job.Name
			if !strings.HasPrefix(name, "OpenFaaS-"+service) {
				continue
			}

			image := allocation.Job.TaskGroups[0].Tasks[0].Config["image"]

			if image != nil {
				functions = append(functions, requests.Function{
					Name:            allocation.Job.TaskGroups[0].Tasks[0].Name,
					Image:           image.(string),
					Replicas:        uint64(*allocation.Job.TaskGroups[0].Count),
					InvocationCount: 0,
				})
			}
		}

	}

	return functions, nil

}

func writeError(s string, w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
	log.Println(s + " " + err.Error())
	return
}

// MakeReader implements the OpenFaaS reader handler
func MakeReader(n string, client nomad.Allocations) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %v", n, r.URL.Path)
		defer r.Body.Close()

		body, _ := ioutil.ReadAll(r.Body)

		request := requests.Function{}
		err := json.Unmarshal(body, &request)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Printf("%+v", request)

		// Not sure if prefix is the right option
		options := api.QueryOptions{}
		// options.Prefix = "faas_function"

		allocations, _, err := client.List(&options)
		if err != nil {
			writeError("alloc 1", w, err)
			return
		}

		functions, err := getFunctions(request.Name, client, allocations)
		if err != nil {
			writeError("get functions", w, err)
			return
		}

		functionBytes, _ := json.Marshal(functions)
		w.Header().Set("Content-Type", "application/json")
		if len(functions) == 0 {
			w.WriteHeader(404)
			w.Write(functionBytes)
			return
		}
		w.WriteHeader(200)
		w.Write(functionBytes)
	}
}
