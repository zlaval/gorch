package controller

import (
	"fmt"
	"gorch/pkg/deployment"
	"gorch/pkg/metrics"
	"gorch/pkg/pod"
	"gorch/pkg/rest"
)

type WorkerStatus int

const (
	WorkerUp WorkerStatus = iota
	WorkerDown
)

type WorkerEntity struct {
	Name    string
	Address string
	IP      string
	Port    int

	Status  WorkerStatus
	Metrics metrics.Metrics
}

type WorkerClient struct {
}

func (w WorkerClient) CreatePod(workerAddress, id string, d deployment.Deployment) (pod.ClientResponse, error) {
	return rest.Post(
		fmt.Sprintf("%s/pods", workerAddress),
		pod.CreateRequest{
			Config: d.Config,
			ID:     id,
			Name:   d.Name,
		},
		rest.ExtractBody[pod.ClientResponse],
	)
}
