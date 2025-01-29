package cli

import (
	"github.com/spf13/cobra"
	"gorch/pkg/command"
)

type WorkerCmd struct {
	*cobra.Command
}

func NewWorkerCmd() *cobra.Command {
	c := &WorkerCmd{}
	c.Command = &cobra.Command{
		Use:   "workers",
		Short: "List of workers",
		RunE:  c.run,
	}
	return c.Command
}

func (c WorkerCmd) run(_ *cobra.Command, _ []string) error {
	cmd := command.Command{
		Action: command.Workers,
	}
	return ExecuteCommand(cmd)
}
