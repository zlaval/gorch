package controller

import (
	"fmt"
	"github.com/renstrom/shortuuid"
	"gorch/pkg/command"
	"gorch/pkg/deployment"
	"time"
)

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
		return p.create(cmd)
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

func (p *CommandProcessor) create(cmd command.Command) ([]string, error) {
	createdAt := time.Now().UTC()
	d := deployment.Deployment{
		ID:        shortuuid.New(),
		Name:      cmd.DeploymentName,
		Replicas:  cmd.Replicas,
		State:     deployment.Created,
		CreatedAt: createdAt,
		Config:    cmd.Config,
		History: []string{
			fmt.Sprintf("created %s", createdAt.Format(time.DateTime)),
		},
	}

	if err := p.db.SaveDeployment(d); err != nil {
		return nil, fmt.Errorf("save deployment: %w", err)
	}

	return []string{fmt.Sprintf("Deployment has been created")}, nil
}
