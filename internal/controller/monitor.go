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

	for d := range slices.Values(deployments) {
		m.handleDeployment(d)
	}
}

func (m Monitor) handleDeployment(dpl deployment.Deployment) {
	pods, err := m.db.LoadPodsByDeploymentName([]byte(dpl.Name))
	if err != nil {
		slog.Error("load pods", slog.Any("error", err))
	}

	if dpl.State == deployment.Deleted && len(pods) == 0 {
		if err := m.db.DeleteDeployment(dpl.ID); err != nil {
			slog.Error("delete deployment", slog.Any("error", err), slog.String("id", dpl.ID))
		}
		return
	}

	m.markPodsOfDeletedDeployment(&dpl, pods)
	m.deleteRemovedPods(&dpl, pods)
	m.markFailedPods(&dpl, pods)
	m.balancePods(&dpl, pods)

	if err := m.db.SaveDeployment(dpl); err != nil {
		slog.Error("save deployment", slog.Any("error", err))
	}

}

func (m Monitor) markPodsOfDeletedDeployment(dpl *deployment.Deployment, pods []pod.Pod) {
	//Mark pods of a deleted deployment as deletable
	if dpl.State != deployment.Deleted {
		return
	}

	for i, p := range pods {
		if p.MarkedForDelete || p.Deleted {
			continue
		}
		p.MarkedForDelete = true
		if err := m.db.SavePod(p); err != nil {
			slog.Error("mark pod as deleted", slog.Any("error", err), slog.Any("pod", p))
			continue
		}
		pods[i] = p
		m.scheduler.Add(Task{*dpl, p})
		dpl.History = append(dpl.History,
			fmt.Sprintf("remove pod %s at %s", p.ID, time.Now().UTC().Format(time.DateTime)),
		)
	}
}

func (m Monitor) deleteRemovedPods(dpl *deployment.Deployment, pods []pod.Pod) {
	//Remove deleted pods from the database
	for p := range slices.Values(pods) {
		if p.Deleted {
			if err := m.db.DeletePod(p.ID); err != nil {
				slog.Error("delete pod", slog.Any("error", err))
			} else {
				dpl.History = append(dpl.History,
					fmt.Sprintf("delete pod %s at %s", p.ID, time.Now().UTC().Format(time.DateTime)),
				)
			}
		}
	}
}

func (m Monitor) markFailedPods(dpl *deployment.Deployment, pods []pod.Pod) {
	//Mark dead pods for deletion (killed manually, by system...)
	for i, p := range pods {
		if p.MarkedForDelete || p.Deleted {
			continue
		}

		switch p.State {
		case "exited", "dead":
			p.MarkedForDelete = true
			if err := m.db.SavePod(p); err != nil {
				slog.Error("save pod", slog.Any("error", err))
				continue
			}
			pods[i] = p
			m.scheduler.Add(Task{*dpl, p})
			dpl.History = append(dpl.History,
				fmt.Sprintf("remove pod %s at %s", p.ID, time.Now().UTC().Format(time.DateTime)),
			)
		}
	}

}

func (m Monitor) balancePods(dpl *deployment.Deployment, pods []pod.Pod) {
	// create/remove pods to match for the deployment scale
	if dpl.State == deployment.Deleted {
		return
	}

	runningPodsCount := 0
	for p := range slices.Values(pods) {
		if !p.MarkedForDelete && !p.Deleted {
			runningPodsCount++
		}
	}

	if runningPodsCount > dpl.Replicas {
		delta := runningPodsCount - dpl.Replicas
		for i := 0; i < delta; i++ {
			p := pods[i]
			p.MarkedForDelete = true
			if err := m.db.SavePod(p); err != nil {
				slog.Error("save pod", slog.Any("error", err))
				continue
			}
			pods[i] = p
			m.scheduler.Add(Task{*dpl, p})
			dpl.History = append(dpl.History,
				fmt.Sprintf("remove pod %s at %s", p.ID, time.Now().UTC().Format(time.DateTime)),
			)
		}
	} else if runningPodsCount < dpl.Replicas {
		delta := dpl.Replicas - runningPodsCount
		for i := 0; i < delta; i++ {
			p := pod.Pod{
				ID:           fmt.Sprintf("%s-%s", dpl.Name, shortuuid.New()),
				DeploymentID: dpl.ID,
				State:        "created",
			}
			if err := m.db.SavePod(p); err != nil {
				slog.Error("save pod", slog.Any("error", err))
				continue
			}
			pods = append(pods, p)
			m.scheduler.Add(Task{*dpl, p})
			dpl.History = append(dpl.History,
				fmt.Sprintf("create pod %s at %s", p.ID, time.Now().UTC().Format(time.DateTime)),
			)
		}
	}

}
