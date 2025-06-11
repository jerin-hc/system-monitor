package scheduler

import (
	"encoding/json"
	"time"

	"github.com/j3rryCodes/system-monitor/internal/collector"
	"github.com/j3rryCodes/system-monitor/internal/logger"
	"github.com/j3rryCodes/system-monitor/internal/storage/influxdb"
	"go.uber.org/zap"
)

func Schedule(sec string) {
	loadMetrics()
	d, _ := time.ParseDuration(sec)
	ticker := time.NewTicker(d)
	defer ticker.Stop()

	for range ticker.C {
		logger.Logger().Debug("Loading system metirics")
		loadMetrics()
	}
}

func loadMetrics() {
	fields := make(map[string]interface{})
	c, err := collector.ExtractCPUMetics()
	if err == nil {
		fields["cpu"] = c[0]
	}
	d, err := collector.ExtractDiskMetrics()
	if err == nil {
		jsonBytes, err := json.Marshal(d)
		if err == nil {
			jsonString := string(jsonBytes)
			fields["disk"] = jsonString
		}

	}
	m, err := collector.ExtractMemoryMetrics()
	if err == nil {
		jsonBytes, err := json.Marshal(m)
		if err == nil {
			jsonString := string(jsonBytes)
			fields["memory"] = jsonString
		}
	}
	logger.Logger().Info("Loaded system metics", zap.String("utc-time", time.Now().UTC().String()))

	influxdb.AddPoint("system_info", "metrics", fields)
}
