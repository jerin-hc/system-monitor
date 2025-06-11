package system

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/j3rryCodes/system-monitor/internal/collector"
	"github.com/j3rryCodes/system-monitor/internal/logger"
)

const prefix = "system"

func RegisterRoute(r *gin.Engine) {
	r.GET(fmt.Sprintf("%s/cpu", prefix), getCpuUsage)
	r.GET(fmt.Sprintf("%s/memory", prefix), getMemoryUsage)
	r.GET(fmt.Sprintf("%s/disk", prefix), getDiskUsage)
}

func getCpuUsage(c *gin.Context) {
	usage, err := collector.ExtractCPUMetics()
	if err != nil {
		c.JSON(503, gin.H{
			"message": "unable to fetch CPU usage",
			"status":  "failed",
		})
		logger.Logger().Error(err.Error())
	}
	c.JSON(200, gin.H{
		"cpuUsagePercentage": usage,
		"status":             "success",
	})
}

func getMemoryUsage(c *gin.Context) {
	memStats, err := collector.ExtractMemoryMetrics()
	if err != nil {
		c.JSON(503, gin.H{
			"message": "unable to fetch Memory usage",
			"status":  "failed",
		})
		logger.Logger().Error(err.Error())
	}
	c.JSON(200, gin.H{
		"metic":  memStats,
		"status": "success",
	})
}

func getDiskUsage(c *gin.Context) {
	diskStats, err := collector.ExtractDiskMetrics()
	if err != nil {
		c.JSON(503, gin.H{
			"message": "unable to fetch Disk metric",
			"status":  "failed",
		})
		logger.Logger().Error(err.Error())
	}
	c.JSON(200, gin.H{
		"metic":  diskStats,
		"status": "success",
	})
}
