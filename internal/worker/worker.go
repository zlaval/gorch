package worker

import (
	"Gorch/pkg/pod"
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"io"
	"math"
	"os"
	"slices"
	"strings"
	"time"
)

type Worker struct {
	client *client.Client
}

func NewWorker(client *client.Client) Worker {
	return Worker{client}
}

func (w Worker) Start(ctx context.Context, request pod.CreateRequest) pod.ClientResponse {
	// pull image
	reader, err := w.client.ImagePull(ctx, request.Image, image.PullOptions{})
	if err != nil {
		return pod.ClientResponse{
			ID:       request.ID,
			ErrorMsg: fmt.Sprintf("pull image: %s", err),
		}
	}
	_, _ = io.Copy(os.Stdout, reader)

	//create container
	cc := &container.Config{
		Image:        request.Image,
		Cmd:          request.Cmd,
		Env:          request.Env,
		ExposedPorts: request.ExposedPorts,
	}

	pm := nat.PortMap{}
	for p := range request.ExposedPorts {
		pm[p] = []nat.PortBinding{}
	}

	hc := &container.HostConfig{
		Resources: container.Resources{
			Memory:   request.MemoryRequest * 1024 * 1024,
			NanoCPUs: int64(float64(request.CPURequest) * math.Pow(10, 6)),
		},
		RestartPolicy:   container.RestartPolicy{Name: container.RestartPolicyDisabled},
		PortBindings:    pm,
		PublishAllPorts: false,
	}

	name := strings.ReplaceAll(
		fmt.Sprintf("%s-%s", request.Name, request.ID),
		" ",
		"-",
	)

	resp, err := w.client.ContainerCreate(ctx, cc, hc, nil, nil, name)
	if err != nil {
		return pod.ClientResponse{
			ID:       request.ID,
			ErrorMsg: fmt.Sprintf("create container: %s", err),
		}
	}

	//start container
	if err := w.client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return pod.ClientResponse{
			ID:          request.ID,
			ContainerID: resp.ID,
			StartedAt:   time.Now().UTC(),
			ErrorMsg:    fmt.Sprintf("create container: %s", err),
		}
	}

	// retrieve container data
	cjson, err := w.client.ContainerInspect(ctx, resp.ID)
	if err != nil {
		return pod.ClientResponse{
			ID:          request.ID,
			ContainerID: resp.ID,
			StartedAt:   time.Now().UTC(),
		}
	}

	// get ephemeral ports
	ports := make([]string, 0, len(cjson.NetworkSettings.Ports))
	for _, pm := range cjson.NetworkSettings.Ports {
		for pb := range slices.Values(pm) {
			ports = append(ports, pb.HostPort)
		}
	}

	return pod.ClientResponse{
		ID:             request.ID,
		ContainerID:    resp.ID,
		StartedAt:      time.Now().UTC(),
		IP:             cjson.NetworkSettings.IPAddress,
		EphemeralPorts: ports,
	}
}

func (w *Worker) Stop(ctx context.Context, containerID string) pod.ClientResponse {
	//Stop container
	if err := w.client.ContainerStop(ctx, containerID, container.StopOptions{}); err != nil {
		return pod.ClientResponse{
			ContainerID: containerID,
			ErrorMsg:    fmt.Sprintf("stop container: %s", err),
		}
	}

	err := w.client.ContainerRemove(ctx, containerID, container.RemoveOptions{
		RemoveVolumes: true,
	})
	if err != nil {
		return pod.ClientResponse{
			ContainerID: containerID,
			ErrorMsg:    fmt.Sprintf("remove container: %s", err),
			FinishedAt:  time.Now().UTC(),
		}
	}

	return pod.ClientResponse{
		ContainerID: containerID,
		FinishedAt:  time.Now().UTC(),
	}

}
