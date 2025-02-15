package controller

import (
	"context"
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
		m.collectNodeMetrics(w)
	}
}

func (m MetricsCollector) collectNodeMetrics(w WorkerEntity) {
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
