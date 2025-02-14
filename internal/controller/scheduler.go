package controller

import (
	"context"
	"gorch/pkg/deployment"
	"gorch/pkg/pod"
	"log/slog"
)

type Task struct {
	Deployment deployment.Deployment
	Pod        pod.Pod
}

type Scheduler struct {
	ch     chan Task
	db     *Database
	client *WorkerClient
}

func NewScheduler(db *Database, client *WorkerClient) *Scheduler {
	return &Scheduler{
		ch:     make(chan Task),
		db:     db,
		client: client,
	}
}

func (s *Scheduler) Run(ctx context.Context) {
	for {
		select {
		case task := <-s.ch:
			s.schedule(task)
		case <-ctx.Done():
			return
		}
	}
}

func (s *Scheduler) Add(task Task) {
	s.ch <- task
}

func (s *Scheduler) schedule(task Task) {
	workers, err := s.db.LoadWorkers()
	if err != nil {
		slog.Error("loading workers", slog.Any("error", err))
		return
	}

	if len(workers) == 0 {
		slog.Warn("no worker")
		return
	}

	selectedWorker := workers[0]

	p := task.Pod

	res, err := s.client.CreatePod(selectedWorker.Address, p.ID, task.Deployment)
	if err != nil {
		slog.Error("creating pod", slog.Any("error", err))

		p.State = "dead"
		_ = s.db.SavePod(p)
		return
	}

	p.ContainerID = res.ContainerID
	p.StartedAt = res.StartedAt
	p.FinishedAt = res.FinishedAt
	p.Worker = selectedWorker.Name
	p.IP = res.IP
	p.EphemeralPorts = res.EphemeralPorts

	if err := s.db.SavePod(p); err != nil {
		slog.Error("save pod", slog.Any("error", err), slog.Any("id", p.ID))
	}

}
