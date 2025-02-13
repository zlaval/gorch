package controller

import (
	"context"
	"fmt"
	"log/slog"
)

type Task struct {
}

type Scheduler struct {
	ch chan Task
	db *Database
}

func NewScheduler(db *Database) *Scheduler {
	return &Scheduler{
		ch: make(chan Task),
		db: db,
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

	fmt.Println(selectedWorker)
	fmt.Println(task)
}
