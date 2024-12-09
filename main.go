package main

import (
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	"io"
	"log"
	"os"
)

func main() {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	reader, err := cli.ImagePull(ctx, "zalerix/webapp", image.PullOptions{})
	if err != nil {
		log.Fatal(err)
	}
	io.Copy(os.Stdout, reader)

	cc := &container.Config{
		Image: "zalerix/webapp",
		Cmd:   []string{},
		ExposedPorts: nat.PortSet{
			nat.Port("8000"): struct{}{},
		},
	}

	pm := nat.PortMap{}
	pm["8000"] = []nat.PortBinding{
		{HostPort: "9115"},
	}

	hc := &container.HostConfig{
		PortBindings: pm,
	}

	resp, err := cli.ContainerCreate(ctx, cc, hc, nil, nil, "mycontainer")
	if err != nil {
		log.Fatal(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		log.Fatal(err)
	}

	_, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	if err := <-errCh; err != nil {
		log.Fatal(err)
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		log.Fatal(err)
	}
	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
}
