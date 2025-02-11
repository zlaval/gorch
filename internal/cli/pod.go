package cli

import (
	"github.com/spf13/cobra"
	"gorch/pkg/command"
)

type podCmd struct {
	*cobra.Command
}

func NewPodCmd() *cobra.Command {
	c := &podCmd{
		Command: &cobra.Command{
			Use:   "pod",
			Short: "pod actions",
		},
	}

	c.AddCommand(
		newListPodsCmd(),
		newInspectCmd(),
		newLogsCmd(),
	)

	return c.Command
}

type listPodsCmd struct {
	*cobra.Command
}

func newListPodsCmd() *cobra.Command {
	c := &listPodsCmd{}
	c.Command = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "list pods",
		RunE:    c.run,
	}

	return c.Command
}

func (c *listPodsCmd) run(_ *cobra.Command, _ []string) error {
	cmd := command.Command{
		Action: command.ListPods,
	}
	return ExecuteCommand(cmd)
}

type inspectCmd struct {
	*cobra.Command
}

func newInspectCmd() *cobra.Command {
	c := &inspectCmd{}
	c.Command = &cobra.Command{
		Use:     "inspect (pod-id)",
		Aliases: []string{"i"},
		Short:   "inspects pod",
		Args:    cobra.ExactArgs(1),
		RunE:    c.run,
	}

	return c.Command
}

func (c *inspectCmd) run(_ *cobra.Command, args []string) error {
	cmd := command.Command{
		Action: command.Inspect,
		PodID:  args[0],
	}
	return ExecuteCommand(cmd)
}

type logsCmd struct {
	*cobra.Command
}

func newLogsCmd() *cobra.Command {
	c := &logsCmd{}
	c.Command = &cobra.Command{
		Use:   "logs (pod-id)",
		Short: "get pod logs",
		Args:  cobra.ExactArgs(1),
		RunE:  c.run,
	}

	return c.Command
}

func (c *logsCmd) run(_ *cobra.Command, args []string) error {
	cmd := command.Command{
		Action: command.Log,
		PodID:  args[0],
	}
	return ExecuteCommand(cmd)
}
