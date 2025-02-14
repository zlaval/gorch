package controller

import (
	"fmt"
	"gorch/pkg/command"
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

func (w WorkerClient) CreatePod(workerAddress, id string, c command.Command) (pod.ClientResponse, error) {
	return rest.Post(
		fmt.Sprintf("%s/pods", workerAddress),
		pod.CreateRequest{
			Config: c.Config,
			ID:     id,
			Name:   c.DeploymentName,
		},
		rest.ExtractBody[pod.ClientResponse],
	)
}
