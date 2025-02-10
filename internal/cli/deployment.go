package cli

import (
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"gorch/pkg/command"
	"gorch/pkg/pod"
	"os"
	"slices"
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

	name          string
	replicas      int
	image         string
	command       []string
	env           []string
	memory        int64
	cpu           int64
	containerPort string
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
	fs := c.Flags()

	fs.StringVarP(&c.filePath, "file", "f", "", "-f file_path")

	fs.StringVarP(&c.name, "name", "n", "", "name of the deployment")
	fs.IntVarP(&c.replicas, "replicas", "r", 1, "number of pods")
	fs.StringVarP(&c.image, "image", "i", "", "container image")
	fs.StringSliceVarP(&c.command, "cmd", "c", []string{}, "command array. format: -cmd=\"val1,val2\"")
	fs.StringSliceVarP(&c.env, "env", "e", []string{}, "environment array. format: -env=\"key1=val1,key2=val2\"")
	fs.Int64VarP(&c.memory, "memory", "m", 0, "memory limit of a pod. 0 means unlimited")
	fs.Int64VarP(&c.cpu, "cpu", "u", 0, "cpu limit of the pod. 0 means unlimited")
	fs.StringVarP(&c.containerPort, "port", "p", "80", "exposed port of the pod")

	c.MarkFlagsOneRequired("file", "name")
	c.MarkFlagsRequiredTogether("name", "image")
	c.markFlagsMutuallyExclusiveTo("file",
		"name", "replicas", "cmd", "env", "memory", "cpu", "port",
	)

}

func (c *deployCmd) markFlagsMutuallyExclusiveTo(flag string, mes ...string) {
	for m := range slices.Values(mes) {
		c.MarkFlagsMutuallyExclusive(flag, m)
	}
}

func (c *deployCmd) run(_ *cobra.Command, _ []string) error {
	var cmd command.Command

	if c.filePath != "" {
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

		cmd = command.Command{
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
	} else {
		cmd = command.Command{
			Action:         command.Create,
			Replicas:       c.replicas,
			DeploymentName: c.name,
			Config: pod.Config{
				Image: c.image,
				Cmd:   c.command,
				Env:   c.env,
				ExposedPorts: nat.PortSet{
					nat.Port(c.containerPort): struct{}{},
				},
				CPURequest:    c.cpu,
				MemoryRequest: c.memory,
			},
		}
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
