package worker

import (
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"log"
)

type workerCmd struct {
	*cobra.Command

	name string
}

func Cmd() *cobra.Command {
	c := &workerCmd{}
	c.Command = &cobra.Command{
		Use:   "worker",
		Short: "start worker",
		RunE:  c.run,
	}

	c.init()

	return c.Command
}

func (c *workerCmd) init() {
	fs := c.Flags()
	fs.StringVarP(&c.name, "name", "n", "worker-1", "Unique name of the worker")
}

func (c *workerCmd) run(cmd *cobra.Command, _ []string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	w := NewWorker(cli)
	m := NewMetricsCollector(cli)
	a := NewApi(w, m)

	return a.Run(cmd.Context())
}
