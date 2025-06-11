package collector

import (
	"time"

	"github.com/j3rryCodes/system-monitor/internal/logger"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"go.uber.org/zap"
)

func ExtractDiskMetrics() (diskStats map[string]any, err error) {
	usage, err := disk.Usage("/")
	if err != nil {
		logger.Logger().Error("error while fetching DISK metric", zap.Any("error", err))
	} else {
		diskStats = map[string]interface{}{
			"path":           usage.Path,
			"totalDiskSpace": usage.Total,
			"totalDiskUsage": usage.Used,
			"freeDiskSpace":  usage.Free,
		}
	}
	return diskStats, err
}

func ExtractMemoryMetrics() (memStats map[string]any, err error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		logger.Logger().Error("error while fetching DISK metric", zap.Any("error", err))
	} else {
		memStats = map[string]interface{}{
			"totalMemorySpace": vmStat.Total,
			"totalMemoryUsage": vmStat.Used,
			"freeMemorySpace":  vmStat.Free,
		}
	}
	return memStats, err
}

func ExtractCPUMetics() ([]float64, error) {
	usage, err := cpu.Percent(time.Second, false)
	if err != nil {
		logger.Logger().Error("error while fetching CPU metric", zap.Any("error", err))
	}
	return usage, err
}
