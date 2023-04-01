package main

import (
	"log"
	"os"

	bootstrap "github.com/alexellis/faas-provider"
	"github.com/alexellis/faas-provider/types"
	"github.com/blinkinglight/faas-nomad/handlers"
	"github.com/hashicorp/nomad/api"
)

func main() {
	region := os.Getenv("NOMAD_REGION")
	address := os.Getenv("NOMAD_ADDR")

	c := api.DefaultConfig()
	c.SecretID = os.Getenv("NOMAD_TOKEN")

	client, err := api.NewClient(c.ClientConfig(region, address, false))
	if err != nil {
		log.Fatal(err)
	}

	handlers := &types.FaaSHandlers{
		FunctionReader: handlers.MakeReader(client.Allocations()),
		DeployHandler:  handlers.MakeNull(),
		DeleteHandler:  handlers.MakeNull(),
		ReplicaReader:  handlers.MakeNull(),
		FunctionProxy:  handlers.MakeNull(),
		ReplicaUpdater: handlers.MakeNull(),
	}
	config := &types.FaaSConfig{}

	bootstrap.Serve(handlers, config)
}
