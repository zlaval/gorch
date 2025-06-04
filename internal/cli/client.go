package cli

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"gorch/pkg/command"
	"gorch/pkg/rest"
	"os"
	"slices"
	"text/tabwriter"
)

type cliConfig struct {
	Cluster string
}

type client struct {
	controllerAddress string
}

func newClient() (*client, error) {
	b, err := os.ReadFile("config")
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	var cf cliConfig
	if err := yaml.Unmarshal(b, &cf); err != nil {
		return nil, fmt.Errorf("parse file: %w", err)
	}

	return &client{controllerAddress: cf.Cluster}, nil
}

func ExecuteCommand(cmd command.Command) error {
	cli, err := newClient()
	if err != nil {
		return fmt.Errorf("'config' file is mandatory: %w", err)
	}

	res, err := rest.Post(
		fmt.Sprintf("%s/command", cli.controllerAddress),
		cmd,
		rest.ExtractBody[command.Response],
	)

	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 5, 0, 4, ' ', 0)
	for r := range slices.Values(res.Data) {
		_, _ = fmt.Fprintln(w, r)
	}
	_ = w.Flush()

	return nil
}
