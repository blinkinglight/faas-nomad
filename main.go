package main

import (
	"log"
	"os"

	"github.com/blinkinglight/faas-nomad/handlers"
	"github.com/hashicorp/nomad/api"
	bootstrap "github.com/openfaas/faas-provider"
	"github.com/openfaas/faas-provider/types"
)

func main() {
	// region := os.Getenv("NOMAD_REGION")
	address := os.Getenv("NOMAD_ADDR")

	// c := api.DefaultConfig()
	// c.SecretID = os.Getenv("NOMAD_TOKEN")
	c := &api.Config{
		Address:  address,
		SecretID: os.Getenv("NOMAD_TOKEN"),
	}

	client, err := api.NewClient(c)
	if err != nil {
		log.Fatal(err)
	}

	handlers := &types.FaaSHandlers{
		FunctionLister: handlers.MakeReader(client.Allocations()),
		DeployFunction: handlers.MakeDeploy(client.Jobs()),
		DeleteFunction: handlers.MakeNull("DeleteFunction"),
		FunctionProxy:  handlers.MakeNull("FunctionProxy"),
		ListNamespaces: handlers.MakeNull("ListNamespaces"),
		UpdateFunction: handlers.MakeReader(client.Allocations()),
		FunctionStatus: handlers.MakeReader(client.Allocations()),
		ScaleFunction:  handlers.MakeNull("ScaleFunction"),
		Secrets:        handlers.MakeNull("Secrets"),
		Logs:           handlers.MakeNull("Logs"),
		Health:         handlers.MakeNull("Health"),
		Info:           handlers.MakeNull("Info"),
	}
	config := &types.FaaSConfig{}
	port := 9999
	config.TCPPort = &port
	// http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("Hello from Nomad!"))
	// }))
	// http.ListenAndServe(":9999", nil)
	// bootstrap.Serve(handlers, config)
	log.Printf("Starting...")
	bootstrap.Serve(handlers, config)
}
