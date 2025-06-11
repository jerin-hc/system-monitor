package influxdb

import (
	"context"
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/domain"
	"github.com/j3rryCodes/system-monitor/internal/logger"
	"go.uber.org/zap"
)

var b string
var w api.WriteAPI
var q api.QueryAPI

func Init(address string, token string, org string, bucket string) {
	b = bucket
	ctx, cancel := context.WithTimeoutCause(
		context.Background(),
		5*time.Second,
		fmt.Errorf("connection timeout, unable to connect to InfluxDB[%s] after 5sec", address),
	)
	defer cancel()

	client := influxdb2.NewClient(address, token)

	r, err := client.Health(ctx)
	if r.Status != "pass" {
		logger.Logger().Fatal("InfluxDB health check failed", zap.String("address", address), zap.String("error", err.Error()))
	}
	organization := verifyOrg(ctx, client, org)
	verifyBucket(ctx, client, organization, bucket)

	w = client.WriteAPI(org, bucket)
	q = client.QueryAPI(org)

	logger.Logger().Info("InfluxDB is reachable and token is valid", zap.String("address", address))
}

func verifyOrg(ctx context.Context, client influxdb2.Client, org string) *domain.Organization {
	organizationAPI := client.OrganizationsAPI()
	orgs, err := organizationAPI.GetOrganizations(ctx)

	if err == nil && orgs != nil {
		for _, o := range *orgs {
			if o.Name == org {
				return &o
			}
		}
		logger.Logger().Fatal("Provided organization not found", zap.String("providedOrganization", org))
	} else {
		logger.Logger().Fatal("Failed to fetch organizations: invalid or unauthorized InfluxDB token", zap.Error(err))
	}
	return nil
}

func verifyBucket(ctx context.Context, client influxdb2.Client, org *domain.Organization, bucket string) {
	bucketsAPI := client.BucketsAPI()
	buckets, err := bucketsAPI.FindBucketsByOrgName(ctx, org.Name)

	if err == nil && buckets != nil {
		for _, b := range *buckets {
			if b.Name == bucket {
				return
			}
		}
		logger.Logger().Warn("Provided bucket not found", zap.String("organization", org.Name), zap.String("bucket", bucket))
		logger.Logger().Info("Creating new bucket", zap.String("organization", org.Name), zap.String("bucket", bucket))
		bucketsAPI.CreateBucketWithName(ctx, org, bucket)
	} else {
		logger.Logger().Fatal("Failed to fetch buckets: invalid organization name or unauthorized InfluxDB token", zap.Error(err), zap.String("organization", org.Name))
	}
}

func AddPoint(measurement string, tag string, fields map[string]any) {
	time := time.Now().UTC()
	t := map[string]string{
		"type": tag,
	}
	p := influxdb2.NewPoint(
		measurement,
		t,
		fields,
		time,
	)
	w.WritePoint(p)
	logger.Logger().Debug("Asynchronously added new point to InfluxDB", zap.Any("utc-time", time), zap.String("measurement", measurement), zap.Any("tag", tag), zap.Any("fields", fields))
}

func GetPoints(start time.Duration) (string, error) {
	ctx, cancel := context.WithTimeoutCause(
		context.Background(),
		5*time.Second,
		fmt.Errorf("connection timeout while quering, unable to connect to InfluxDB"),
	)
	defer cancel()

	query := fmt.Sprintf(`from(bucket: "%s")
  |> range(start: %s)
  |> filter(fn: (r) => r._measurement == "%s")`,
		b, time.Now().UTC().Add(-start).Format(time.RFC3339), "system_info")

	result, err := q.QueryRaw(ctx, query, api.DefaultDialect())

	if err != nil {
		logger.Logger().Error("Unable to fetch points from InfluxDB")
		logger.Logger().Debug("error", zap.Error(err))
	}
	return result, err
}
