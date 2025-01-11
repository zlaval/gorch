package worker

import (
	"Gorch/pkg/metrics"
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"log/slog"
	"slices"
)

type MetricsCollector struct {
	client *client.Client
}

func NewMetricsCollector(client *client.Client) MetricsCollector {
	return MetricsCollector{
		client: client,
	}
}

func (m MetricsCollector) CollectPodMetrics(ctx context.Context) []metrics.PodStats {
	containers, err := m.client.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		slog.Error("collect containers", slog.Any("error", err))
	}

	res := make([]metrics.PodStats, 0, len(containers))

	for c := range slices.Values(containers) {
		stats := metrics.PodStats{
			ContainerID: c.ID,
			State:       c.State,
		}

		i, err := m.client.ContainerInspect(ctx, c.ID)
		if err != nil {
			slog.Error("inspect pod", slog.String("id", c.ID), slog.Any("error", err))
		} else {
			stats.Name = i.Name
			stats.Error = i.State.Error
		}

		res = append(res, stats)
	}

	return res
}

func (m MetricsCollector) CollectNodeMetrics() metrics.Metrics {
	result := metrics.Metrics{}

	memory, err := mem.VirtualMemory()
	if err != nil {
		slog.Error("mem stat", slog.Any("error", err))
	} else {
		result.TotalAvailableMemoryMB = memory.Total / 1024 / 1024
		result.AvailableMemoryMB = memory.Available / 1024 / 1024
	}

	l, err := load.Avg()
	if err != nil {
		slog.Error("load stat", slog.Any("error", err))
	} else {
		result.Load1 = l.Load1
		result.Load5 = l.Load5
		result.Load15 = l.Load15
	}

	d, err := disk.Usage("/")
	if err != nil {
		slog.Error("disk stat", slog.Any("error", err))
	} else {
		result.TotalDiskMB = d.Total / 1024 / 1024
		result.FreeDiskMB = d.Free / 1024 / 1024
	}

	cpus, err := cpu.Info()
	if err != nil {
		slog.Error("cpu usage", slog.Any("error", err))
	} else {
		var cores int32 = 0
		for c := range slices.Values(cpus) {
			cores += c.Cores
		}
		result.TotalCPUCores = cores
	}

	return result
}
