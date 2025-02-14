package deployment

import (
	"gorch/pkg/pod"
	"time"
)

type State int

const (
	Created State = iota
	Deleted
)

func (s State) String() string {
	switch s {
	case Created:
		return "Created"
	case Deleted:
		return "Deleted"
	}
	return ""
}

type Deployment struct {
	ID        string
	Name      string
	Replicas  int
	State     State
	CreatedAt time.Time
	Config    pod.Config
	History   []string
}
