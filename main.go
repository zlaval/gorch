package main

import (
	"Gorch/internal/worker"
	"Gorch/pkg/pod"
	"context"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/renstrom/shortuuid"
	"log"
	"time"
)

func main() {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	w := worker.NewWorker(cli)

	response := w.Start(ctx, pod.CreateRequest{
		ID:   shortuuid.New(),
		Name: "My Worker Pod",
		Config: pod.Config{
			Image: "zalerix/webapp",
			ExposedPorts: nat.PortSet{
				"8000": struct{}{},
			},
			CPURequest:    100,
			MemoryRequest: 300,
		},
	})

	time.Sleep(10 * time.Second)

	w.Stop(ctx, response.ContainerID)

}
