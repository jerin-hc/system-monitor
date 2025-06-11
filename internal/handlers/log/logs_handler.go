package logs

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/j3rryCodes/system-monitor/internal/storage/influxdb"
)

const prefix = "logs"

func RegisterRoute(r *gin.Engine) {
	r.GET(prefix, getMetrics)
}

func getMetrics(c *gin.Context) {
	duration := c.Params.ByName("duration")
	if duration == "" {
		duration = "3600s"
	}
	d, err := time.ParseDuration(duration)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request",
			"status":  "failed",
			"discription": fmt.Sprintf(`Invalid parameter window: %s , A duration string is a possibly signed sequence of decimal numbers,
			 each with optional fraction and a unit suffix, such as "3600s", "1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"),
			  "ms", "s", "m", "h"`, duration),
		},
		)
		return
	}
	r, err := influxdb.GetPoints(d)
	if err != nil {
		c.JSON(503, gin.H{
			"message": "Internal server error",
			"status":  "failed",
		},
		)
		return
	}
	cleanedCSV := strings.ReplaceAll(r, `\r\n`, "\r\n")
	c.Header("Content-Type", "text/csv")
	c.String(200, cleanedCSV)
}
