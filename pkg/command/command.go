package command

import "gorch/pkg/pod"

type Action int

const (
	Create Action = iota + 1
	ListPods
	Delete
	Scale
	Inspect
	Log
	History
	Workers
	ListDeployments
)

type Command struct {
	pod.Config

	Action Action

	DeploymentID string
	PodID        string

	DeploymentName string
	Replicas       int
}

type Response struct {
	Data []string
}
