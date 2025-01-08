package worker

import (
	"context"
	"github.com/docker/docker/client"
	"log"
)

type workerCmd struct {
	name string
}

func Cmd(name string) workerCmd {
	return workerCmd{
		name: name,
	}
}

func (c *workerCmd) Run() error {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	w := NewWorker(cli)
	a := NewApi(w)

	return a.Run(ctx)
}
