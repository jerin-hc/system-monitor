package main

import (
	"flag"

	"github.com/gin-gonic/gin"
	"github.com/j3rryCodes/system-monitor/internal/handlers"
	"github.com/j3rryCodes/system-monitor/internal/logger"
	"github.com/j3rryCodes/system-monitor/internal/scheduler"
	"github.com/j3rryCodes/system-monitor/internal/storage/influxdb"
)

var arrs = flag.String("host", ":8080", "Application host address, default is :8080")
var influxdbAddr = flag.String("influxdbAddr", "", "influxdb Address")
var influxdbToken = flag.String("influxdbToken", "", "influxdb Token")
var influxdbOrg = flag.String("influxdbOrg", "my-org", "influxdb Organization, default is my-org")
var influxdbBucket = flag.String("influxdbBucket", "", "influxdb Bucket")

var pollRate = flag.String("pollRate", "60s", "Metric poll, default is 60s")

func main() {
	flag.Parse()
	logger.Init()

	if *influxdbAddr == "" || *influxdbToken == "" || *influxdbBucket == "" {
		logger.Logger().Fatal("Error: --influxdbAddr, --influxdbToken and --influxdbBucket are required")
	}

	influxdb.Init(*influxdbAddr, *influxdbToken, *influxdbOrg, *influxdbBucket)
	go scheduler.Schedule(*pollRate)

	r := gin.Default()
	handlers.Handle(r)

	r.Run(*arrs)
}
