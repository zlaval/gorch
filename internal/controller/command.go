package controller

import "gorch/pkg/command"

type CommandProcessor struct {
	db        *Database
	scheduler *Scheduler
	client    *WorkerClient
}

func NewCommandProcessor(
	db *Database,
	scheduler *Scheduler,
	client *WorkerClient,
) *CommandProcessor {
	return &CommandProcessor{
		db:        db,
		scheduler: scheduler,
		client:    client,
	}
}

func (p *CommandProcessor) Execute(cmd command.Command) ([]string, error) {
	var res []string
	switch cmd.Action {
	case command.Create:
		return []string{"test"}, nil
	case command.ListPods:
		return []string{"test"}, nil
	case command.Delete:
		return []string{"test"}, nil
	case command.Scale:
		return []string{"test"}, nil
	case command.Inspect:
		return []string{"test"}, nil
	case command.Log:
		return []string{"test"}, nil
	case command.History:
		return []string{"test"}, nil
	case command.Workers:
		return []string{"test"}, nil
	case command.ListDeployments:
		return []string{"test"}, nil
	}

	return res, nil
}
