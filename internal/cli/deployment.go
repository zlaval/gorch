package cli

import (
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"gorch/pkg/command"
	"gorch/pkg/pod"
	"os"
)

type deploymentCmd struct {
	*cobra.Command
}

func NewDeploymentCmd() *cobra.Command {
	c := &deploymentCmd{
		Command: &cobra.Command{
			Use:   "deployment",
			Short: "deployment commands",
		},
	}

	c.AddCommand(
		newDeployCmd(),
	)

	return c.Command
}

type deployCmd struct {
	*cobra.Command

	filePath string
}

func newDeployCmd() *cobra.Command {
	c := &deployCmd{}
	c.Command = &cobra.Command{
		Use:     "deploy",
		Aliases: []string{"d", "dpy"},
		Short:   "deploy pods from manifest file",
		RunE:    c.run,
	}

	c.init()
	return c.Command
}

func (c *deployCmd) init() {
	c.Flags().StringVarP(&c.filePath, "file", "f", "", "-f file_path")
	_ = c.MarkFlagRequired("file")
}

func (c *deployCmd) run(_ *cobra.Command, _ []string) error {
	b, err := os.ReadFile(c.filePath)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}
	var d DeploymentYaml
	if err := yaml.Unmarshal(b, &d); err != nil {
		return fmt.Errorf("parse file %w", err)
	}

	cfg := d.Spec.Container

	envs := make([]string, len(cfg.Env))
	for i, e := range cfg.Env {
		envs[i] = fmt.Sprintf("%s=%s", e.Name, e.Value)
	}

	cmd := command.Command{
		Action:         command.Create,
		Replicas:       d.Spec.Replicas,
		DeploymentName: d.Spec.Name,
		Config: pod.Config{
			Image: cfg.Image,
			Cmd:   cfg.Command,
			Env:   envs,
			ExposedPorts: nat.PortSet{
				nat.Port(cfg.ContainerPort): struct{}{},
			},
			CPURequest:    cfg.Resources.CPU,
			MemoryRequest: cfg.Resources.Memory,
		},
	}

	return ExecuteCommand(cmd)
}

type DeploymentYaml struct {
	Spec struct {
		Name      string `yaml:"name"`
		Replicas  int    `yaml:"replicas"`
		Container struct {
			Image   string   `yaml:"image"`
			Command []string `yaml:"command"`
			Env     []struct {
				Name  string `yaml:"name"`
				Value string `yaml:"value"`
			}
			Resources struct {
				Memory int64 `yaml:"memory"`
				CPU    int64 `yaml:"cpu"`
			}
			ContainerPort string `yaml:"containerPort"`
		}
	}
}
