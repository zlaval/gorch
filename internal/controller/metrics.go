package controller

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"time"
)

type MetricsCollector struct {
	db     *Database
	client *WorkerClient
}

func NewMetricsCollector(db *Database, client *WorkerClient) *MetricsCollector {
	return &MetricsCollector{
		db:     db,
		client: client,
	}
}

func (m MetricsCollector) Run(ctx context.Context) {
	t := time.NewTicker(10 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			m.collectMetrics()
		case <-ctx.Done():
			return
		}
	}
}

func (m MetricsCollector) collectMetrics() {
	workers, err := m.db.LoadWorkers()
	if err != nil {
		slog.Error("loading workers", slog.Any("error", err))
		return
	}

	for w := range slices.Values(workers) {
		m.checkHealth(w)
		m.collectNodeMetrics(w)
		m.checkPods(w)
	}
}

func (m MetricsCollector) checkHealth(w WorkerEntity) {
	err := m.client.Health(w.Address)
	if err != nil {
		pods, err := m.db.LoadPods()
		if err != nil {
			slog.Error("load pods", slog.Any("error", err))
			return
		}

		//TASK: wait 3 health check grace period and only delete pods after that
		for p := range slices.Values(pods) {
			if p.Worker == w.Name {
				_ = m.db.DeletePod(p.ID)
				d, err := m.db.LoadDeployment(p.DeploymentID)
				if err != nil {
					slog.Error("loading deployment", slog.Any("error", err))
					continue
				}
				d.History = append(d.History,
					fmt.Sprintf("delete pod %s at %s", p.ID, time.Now().UTC().Format(time.DateTime)),
				)
				_ = m.db.SaveDeployment(*d)
			}
		}

		if w.Status == WorkerDown {
			_ = m.db.DeleteWorker(w)
		} else {
			w.Status = WorkerDown
			err = m.db.SaveWorker(w)
			if err != nil {
				slog.Error("save worker", slog.Any("error", err))
			}
		}

	}
}

func (m MetricsCollector) collectNodeMetrics(w WorkerEntity) {
	if w.Status == WorkerDown {
		return
	}

	mtr, err := m.client.NodeMetrics(w.Address)
	if err != nil {
		slog.Error("load node metrics", slog.Any("error", err))
		return
	}
	w.Metrics = mtr
	err = m.db.SaveWorker(w)
	if err != nil {
		slog.Error("save node metrics", slog.Any("error", err))
	}
}

func (m MetricsCollector) checkPods(w WorkerEntity) {
	if w.Status == WorkerDown {
		return
	}

	stats, err := m.client.PodMetrics(w.Address)
	if err != nil {
		slog.Error("pods metrics", slog.Any("error", err))
	}

	pods, err := m.db.LoadPods()
	if err != nil {
		slog.Error("loading pods", slog.Any("error", err))
	}

	for s := range slices.Values(stats) {
		for p := range slices.Values(pods) {
			if p.ContainerID == s.ContainerID {
				p.State = s.State
				_ = m.db.SavePod(p)
				break
			}
		}
	}

	for p := range slices.Values(pods) {
		if p.Worker == w.Name {
			exists := false
			//check if pod exists on the worker
			for s := range slices.Values(stats) {
				if p.ContainerID == s.ContainerID {
					exists = true
				}
			}
			if !exists {
				p.State = "dead"
				_ = m.db.SavePod(p)
			}
		}
	}
}
