package metrics

type Metrics struct {
	TotalAvailableMemoryMB uint64
	AvailableMemoryMB      uint64
	Load1                  float64
	Load5                  float64
	Load15                 float64
	TotalCPUCores          int32
	TotalDiskMB            uint64
	FreeDiskMB             uint64
}
