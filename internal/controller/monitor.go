package controller

import (
	"context"
	"fmt"
	"github.com/renstrom/shortuuid"
	"gorch/pkg/deployment"
	"gorch/pkg/pod"
	"log/slog"
	"slices"
	"time"
)

type Monitor struct {
	db        *Database
	scheduler *Scheduler
}

func NewMonitor(db *Database, scheduler *Scheduler) *Monitor {
	return &Monitor{db, scheduler}
}

func (m Monitor) Run(ctx context.Context) {
	t := time.NewTicker(5 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			m.monitor()
		case <-ctx.Done():
			return
		}
	}
}

func (m Monitor) monitor() {
	deployments, err := m.db.LoadDeployments()
	if err != nil {
		slog.Error("load deployments", slog.Any("error", err))
	}

	//PROBLEM: creates all replicas every 5 seconds
	for d := range slices.Values(deployments) {
		//temporary filter
		if d.State == deployment.Deleted {
			continue
		}

		for i := 0; i < d.Replicas; i++ {
			p := pod.Pod{
				ID:           fmt.Sprintf("%s-%s", d.Name, shortuuid.New()),
				DeploymentID: d.ID,
				State:        "created",
			}

			err := m.db.SavePod(p)
			if err != nil {
				slog.Error("save pod", slog.Any("error", err), slog.Any("deploymentID", d.ID))
				continue
			}

			m.scheduler.Add(Task{d, p})
		}

		//Temporary solution:
		d.State = deployment.Deleted
		_ = m.db.SaveDeployment(d)
	}

}
