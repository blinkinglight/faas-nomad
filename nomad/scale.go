package nomad

import "github.com/hashicorp/nomad/api"

// Job defines the interface for creating a new new job
type Scale interface {
	// Register creates a new Nomad job
	Scale(string, string, *int, string, bool, map[string]interface{}, *api.WriteOptions) (*api.JobRegisterResponse, *api.WriteMeta, error)
}
