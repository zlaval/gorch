package controller

import (
	"context"
	"github.com/renstrom/shortuuid"
	"gorch/pkg/command"
	"log/slog"
)

type Task struct {
	cmd command.Command
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

	_, err = s.client.CreatePod(selectedWorker.Address, shortuuid.New(), task.cmd)
	if err != nil {
		slog.Error("creating pod", slog.Any("error", err))
	}

}
