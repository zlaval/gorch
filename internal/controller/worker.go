package controller

import (
	"fmt"
	"gorch/pkg/deployment"
	"gorch/pkg/metrics"
	"gorch/pkg/pod"
	"gorch/pkg/rest"
	"net/http"
)

type WorkerStatus int

const (
	WorkerUp WorkerStatus = iota
	WorkerDown
)

func (w WorkerStatus) String() string {
	switch w {
	case WorkerUp:
		return "Up"
	case WorkerDown:
		return "Down"
	}
	return ""
}

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

func (w WorkerClient) DeletePod(workerAddress, containerID string) error {
	url := fmt.Sprintf("%s/pods/%s", workerAddress, containerID)
	return rest.Delete(url)
}

func (w WorkerClient) NodeMetrics(workerAddress string) (metrics.Metrics, error) {
	return rest.Get(
		fmt.Sprintf("%s/metrics", workerAddress),
		rest.ExtractBody[metrics.Metrics],
	)
}

func (w WorkerClient) PodMetrics(workerAddress string) ([]metrics.PodStats, error) {
	return rest.Get(
		fmt.Sprintf("%s/pods/metrics", workerAddress),
		rest.ExtractBody[[]metrics.PodStats],
	)
}

func (w WorkerClient) Health(workerAddress string) error {
	st, err := rest.Get(
		fmt.Sprintf("%s/health", workerAddress),
		rest.ExtractStatus,
	)
	if err != nil || st != http.StatusOK {
		return fmt.Errorf("worker %s is unhealthy", workerAddress)
	}
	return nil
}
