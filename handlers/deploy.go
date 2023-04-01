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
func MakeDeploy(n string, client nomad.Job) http.HandlerFunc {
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

		// log.Printf("%+v", request)
		// log.Printf("%s", body)

		// Create job /v1/jobs

		jobName := "OpenFaaS-" + request.Service

		resourcesCPU := 500
		resourcesMemory := 256
		resourcesCount := 1

		config := map[string]interface{}{
			"image":        request.Image,
			"network_mode": "weave",
			"dns_servers":  []string{"10.99.255.100"},
		}

		if _, ok := request.EnvVars["DOCKER_USER"]; ok {
			config["auth"] = map[string]interface{}{
				"username": request.EnvVars["DOCKER_USER"],
				"password": request.EnvVars["DOCKER_PASSWORD"],
			}
		}

		job := &api.Job{
			ID:   &jobName,
			Name: &jobName,
			TaskGroups: []*api.TaskGroup{
				&api.TaskGroup{
					Name:  &jobName,
					Count: &resourcesCount,
					Tasks: []*api.Task{
						&api.Task{
							Name:   request.Service,
							Driver: "docker",
							Config: config,

							Services: []*api.Service{
								&api.Service{
									Name:        jobName,
									PortLabel:   "8080",
									Tags:        []string{"faas_function", "faas"},
									Provider:    "consul",
									AddressMode: "driver",
								},
							},

							Resources: &api.Resources{
								CPU:      &resourcesCPU,
								MemoryMB: &resourcesMemory,
							},
						},
					},
				},
			},
		}
		_, _, err = client.Register(job, nil)
		log.Printf("%v", err)
	}
}
