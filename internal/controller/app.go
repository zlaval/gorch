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
	ctx := cmd.Context()

	db, err := NewDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	sc := NewScheduler(db)

	api := NewApi(db, sc)

	//Starts async scheduler
	go sc.Run(ctx)

	return api.Run(ctx)
}
