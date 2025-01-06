package pod

import (
	"github.com/docker/go-connections/nat"
	"time"
)

type CreateRequest struct {
	Config

	ID   string
	Name string
}

type ClientResponse struct {
	ID             string
	ContainerID    string
	IP             string
	EphemeralPorts []string
	ErrorMsg       string
	StartedAt      time.Time
	FinishedAt     time.Time
}

type Config struct {
	Image         string
	Cmd           []string
	Env           []string
	ExposedPorts  nat.PortSet
	CPURequest    int64
	MemoryRequest int64
}
