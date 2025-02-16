package controller

import (
	"context"
	"gorch/pkg/deployment"
	"gorch/pkg/pod"
	"log/slog"
	"math/rand"
	"slices"
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
		slog.Warn("no worker available")
		p := task.Pod
		p.State = "waiting"
		_ = s.db.SavePod(p)
		return
	}

	//silly select algorithm: randomly pick worker if it has enough memory
	//TASK: better scheduler algorithm
	//Problem: last metric was gathered 10s ago and multiple pods might create the worker
	//Solution: get metrics before a pod created or calculate the remaining resources from last metrics + deployment resources requirement
	var candidates []WorkerEntity
	for w := range slices.Values(workers) {
		if w.Status == WorkerUp &&
			w.Metrics.AvailableMemoryMB > uint64(task.Deployment.Config.MemoryRequest) {
			candidates = append(candidates, w)
		}
	}

	if len(candidates) == 0 {
		slog.Warn("no worker available")
		p := task.Pod
		p.State = "waiting"
		_ = s.db.SavePod(p)
		return
	}

	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	selectedWorker := candidates[0]

	p := task.Pod

	if p.MarkedForDelete {
		err := s.client.DeletePod(selectedWorker.Address, p.ContainerID)
		if err != nil {
			slog.Error("deleting pod", slog.Any("error", err), slog.Any("pod", p))
			return
		}
		p.Deleted = true
		err = s.db.SavePod(p)
		if err != nil {
			slog.Error("save pod", slog.Any("error", err), slog.Any("pod", p))
			return
		}
		return
	}

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
