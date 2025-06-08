package worker

import (
	"fmt"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"gorch/pkg/registry"
	"gorch/pkg/rest"
	"log"
	"log/slog"
	"net"
	"sync"
	"time"
)

type workerCmd struct {
	*cobra.Command

	name              string
	port              int
	controllerAddress string
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
	fs.StringVarP(&c.name, "name", "n", "", "Unique name of the worker")
	fs.IntVarP(&c.port, "port", "p", 8005, "Exposed port of the worker")

	fs.StringVarP(&c.controllerAddress, "controller", "c", "http://localhost:8000",
		"Url of the controller application",
	)

	_ = c.MarkFlagRequired("name")
}

func (c *workerCmd) run(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	// Simple registration flow, easy to understand, but does not work properly.
	// (It waits 10s in main goroutine even if the context was cancelled)

	/*for {
		if err := c.registerWorker(); err == nil {
			break
		}
		slog.Warn("Cannot register worker, keep trying")
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		time.Sleep(10 * time.Second)
	}*/

	// Proper worker registration flow
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		if err := c.registerWorker(); err == nil {
			slog.Info("Worker registration is done")
			return
		} else {
			slog.Info("Cannot register worker, start registration loop")
		}

		for {
			select {
			case <-ctx.Done():
				slog.Warn("Context canceled, stop worker registration process")
				return
			case <-ticker.C:
				if err := c.registerWorker(); err == nil {
					slog.Info("Worker registration is done")
					return
				}
				slog.Warn("Cannot register worker, keep trying")
			}
		}
	}()
	wg.Wait()

	w := NewWorker(cli)
	m := NewMetricsCollector(cli)
	a := NewApi(c.port, w, m)

	return a.Run(ctx)
}

func (c *workerCmd) registerWorker() error {
	ip := getIP()

	r := registry.WorkerRegRequest{
		Name:    c.name,
		Address: fmt.Sprintf("http://%s:%d", ip, c.port),
		IP:      ip,
		Port:    c.port,
	}

	ctrlURL := fmt.Sprintf("%s/register", c.controllerAddress)
	_, err := rest.Post(ctrlURL, r, rest.OmitBody)
	if err != nil {
		return err
	}
	return nil
}

func getIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}
