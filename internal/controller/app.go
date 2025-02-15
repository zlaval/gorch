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

	wc := &WorkerClient{}

	db, err := NewDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	sc := NewScheduler(db, wc)
	mo := NewMonitor(db, sc)
	mc := NewMetricsCollector(db, wc)
	cp := NewCommandProcessor(db, sc, wc)
	api := NewApi(db, cp)

	go sc.Run(ctx)
	go mo.Run(ctx)
	go mc.Run(ctx)

	return api.Run(ctx)
}
