package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorch/internal/worker"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

type RootCommand struct {
	*cobra.Command
}

func newRootCommand() *RootCommand {
	c := &RootCommand{}
	c.Command = &cobra.Command{
		Use:   "gorch",
		Short: "gorch commands",
	}

	viper.AutomaticEnv()

	c.AddCommand(
		worker.Cmd(),
	)

	return c
}

func Execute() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		<-sig
		slog.Info("Shutdown signal has been received")
		cancel()
	}()

	c := newRootCommand()
	return c.Command.ExecuteContext(ctx)
}
