package command

import (
	"fmt"
	"gorch/pkg/pod"
)

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

func (c *Command) String() string {
	return fmt.Sprintf(
		"Command: [Action : %d DeploymentID: %s, DeploymentName: %s, PodID: %s, Replicas: %d, Config: %s]",
		c.Action, c.DeploymentID, c.DeploymentName, c.PodID, c.Replicas, c.Config.String(),
	)
}

type Response struct {
	Data []string
}
