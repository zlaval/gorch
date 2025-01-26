package controller

import "github.com/spf13/cobra"

type controllerCmd struct {
	*cobra.Command
}

func Cmd() *cobra.Command {
	c := &controllerCmd{}
	c.Command = &cobra.Command{
		Use:   "controller",
		Short: "start controller",
		RunE:  c.run,
	}

	return c.Command
}

func (c controllerCmd) run(cmd *cobra.Command, _ []string) error {
	api := NewApi()
	return api.Run(cmd.Context())
}
