package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/alexellis/faas/gateway/requests"
	"github.com/blinkinglight/faas-nomad/nomad"
	"github.com/hashicorp/nomad/api"
)

// MakeDeploy creates a handler for deploying functions
func MakeScale(n string, client nomad.Scale, job *api.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %v", n, r.URL.Path)
		defer r.Body.Close()

		body, _ := ioutil.ReadAll(r.Body)

		request := requests.CreateFunctionRequest{}
		err := json.Unmarshal(body, &request)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		options := api.QueryOptions{}
		// options.Prefix = "OpenFaas
		// options.Prefix = fmt.Sprintf("OpenFaaS-%s", request.Name)

		allocations, _, err := job.Allocations().List(&options)
		if err != nil {
			writeError("get allocs update", w, err)
			return
		}
		log.Printf("request name: %s", request.Service)
		functions, err := getFunctions(request.Service, job.Allocations(), allocations)
		if err != nil {
			writeError("get functions update", w, err)
			return
		}

		functionBytes, _ := json.Marshal(functions)
		w.Header().Set("Content-Type", "application/json")
		if len(functions) == 0 {
			w.WriteHeader(404)
			w.Write(functionBytes)
			return
		}

		// Create job /v1/jobs

		requestFunction := requests.Function{}
		err = json.Unmarshal(body, &requestFunction)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		jobName := "OpenFaaS-" + requestFunction.Name
		count := int(requestFunction.Replicas)
		client.Scale(jobName, jobName, &count, "autoscaling", false, map[string]interface{}{}, nil)
	}
}
