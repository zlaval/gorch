package controller

import (
	"fmt"
	"github.com/renstrom/shortuuid"
	"gorch/pkg/command"
	"gorch/pkg/deployment"
	"slices"
	"strings"
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
		return p.listPods(cmd)
	case command.Delete:
		return []string{"test"}, nil
	case command.Scale:
		return []string{"test"}, nil
	case command.Inspect:
		return p.inspect(cmd)
	case command.Log:
		return p.logs(cmd)
	case command.History:
		return []string{"test"}, nil
	case command.Workers:
		return []string{"test"}, nil
	case command.ListDeployments:
		return p.deployments()
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

func (p *CommandProcessor) deployments() ([]string, error) {
	deployments, err := p.db.LoadDeployments()
	if err != nil {
		return nil, fmt.Errorf("loading deployments: %w", err)
	}

	res := make([]string, 0)
	res = append(res, "Name		Status		Replicas	Created			ID")
	for d := range slices.Values(deployments) {
		line := fmt.Sprintf("%s\t%s\t\t%d\t\t%s\t%s\t",
			d.Name, d.State, d.Replicas, d.CreatedAt.Format(time.DateTime), d.ID,
		)
		res = append(res, line)
	}
	return res, nil
}

func (p *CommandProcessor) listPods(cmd command.Command) ([]string, error) {
	res := make([]string, 0)
	res = append(res, "Name					Deployment	Status		Worker		IP		ContainerPorts		EpheremarPorts		Started")

	deployments, err := p.db.LoadDeployments()
	if err != nil {
		return nil, err
	}

	for d := range slices.Values(deployments) {
		pods, err := p.db.LoadPodsByDeploymentName([]byte(d.Name))
		if err != nil {
			return nil, err
		}
		for p := range slices.Values(pods) {

			var containerPorts []string
			for port := range d.Config.ExposedPorts {
				containerPorts = append(containerPorts, string(port))
			}

			line := fmt.Sprintf("%s\t%s\t%s\t\t%s\t\t%s\t%s\t\t\t%s\t\t\t%s",
				p.ID, d.Name, p.State, p.Worker, p.IP,
				strings.Join(containerPorts, ","),
				strings.Join(p.EphemeralPorts, ","),
				p.StartedAt.Format(time.DateTime),
			)
			res = append(res, line)
		}
	}
	return res, nil
}

func (p *CommandProcessor) inspect(cmd command.Command) ([]string, error) {
	podId := cmd.PodID
	pod, err := p.db.LoadPod(podId)
	if err != nil {
		return nil, err
	}

	worker, err := p.db.LoadWorker(pod.Worker)
	if err != nil {
		return nil, err
	}

	r, err := p.client.InspectPod(worker.Address, pod.ContainerID)
	if err != nil {
		return nil, err
	}

	return []string{r}, nil
}

func (p *CommandProcessor) logs(cmd command.Command) ([]string, error) {
	podId := cmd.PodID
	pod, err := p.db.LoadPod(podId)
	if err != nil {
		return nil, err
	}

	worker, err := p.db.LoadWorker(pod.Worker)
	if err != nil {
		return nil, err
	}

	logs, err := p.client.PodLogs(worker.Address, pod.ContainerID)
	if err != nil {
		return nil, err
	}

	return logs.Lines, nil
}
