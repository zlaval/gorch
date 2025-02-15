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

func (w WorkerClient) InspectPod(workerAddress, containerID string) (string, error) {
	url := fmt.Sprintf("%s/pods/%s", workerAddress, containerID)
	return rest.Get(url, rest.ExtractBody[string])
}

func (w WorkerClient) PodLogs(workerAddress, containerID string) (pod.Logs, error) {
	url := fmt.Sprintf("%s/pods/%s/logs", workerAddress, containerID)
	return rest.Get(url, rest.ExtractBody[pod.Logs])
}
