package controller

import "gorch/pkg/metrics"

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
