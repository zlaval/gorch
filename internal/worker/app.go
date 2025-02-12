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
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	for {
		if err := c.registerWorker(); err == nil {
			break
		}
		slog.Warn("Cannot register worker, keep trying")
		time.Sleep(10 * time.Second)
	}

	w := NewWorker(cli)
	m := NewMetricsCollector(cli)
	a := NewApi(c.port, w, m)

	return a.Run(cmd.Context())
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
