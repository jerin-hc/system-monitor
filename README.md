# system-monitor

A simple Go project to monitor system metrics (CPU, memory, and disk usage) using REST APIs.  
Built with [Gin](https://github.com/gin-gonic/gin), stores metrics in [InfluxDB](https://www.influxdata.com/), uses [Zap](https://github.com/uber-go/zap) for logging, and fetches system stats using [shirou/gopsutil](https://github.com/shirou/gopsutil).

### ✅ Run

```bash
go run . \
  --influxdbAddr="<host>:<port>" \
  --influxdbToken="<token>" \
  --influxdbBucket="<bucket-name>"
```

---

## 📦 Project Structure

```
system-monitoring/
├── go.mod
├── go.sum
├── internal
│   ├── collector
│   │   └── metric_collector.go
│   ├── handlers
│   │   ├── handler.go
│   │   ├── location/
│   │   ├── log/
│   │   │   └── logs_handler.go
│   │   ├── system/
│   │   │   └── system_handler.go
│   │   └── weather/
│   ├── logger/
│   │   └── logger.go
│   ├── scheduler/
│   │   └── sheduler.go
│   └── storage/
│       └── influxdb/
│           └── influxdb_handler.go
└── main.go
```

---

## ⚙️ Configuration

You can customize the following flags:

| Flag               | Description                                  | Default     |
|--------------------|----------------------------------------------|-------------|
| `--host`           | Application host address                     | `:8080`     |
| `--influxdbAddr`   | InfluxDB address (required)                  | `""`        |
| `--influxdbToken`  | InfluxDB authentication token (required)     | `""`        |
| `--influxdbOrg`    | InfluxDB organization                        | `my-org`    |
| `--influxdbBucket` | InfluxDB bucket (required)                   | `""`        |
| `--pollRate`       | Polling interval for metrics collection      | `60s`       |

Example usage:

```bash
go run main.go --influxdbAddr=http://localhost:8086 --influxdbToken=your-token --influxdbBucket=my-bucket
```

---

## 🔌 Endpoints

| Method | Endpoint           | Description                                                              |
|--------|--------------------|--------------------------------------------------------------------------|
| GET    | `/system/cpu`      | Returns current CPU usage in % (live from system)                        |
| GET    | `/system/memory`   | Returns current memory usage (live from system)                          |
| GET    | `/system/disk`     | Returns current disk usage (live from system)                            |
| GET    | `/logs`            | Returns all recorded metrics as CSV (from InfluxDB)                      |

**Note:** `/logs` supports an optional query param `duration`, e.g. `1h`, `30m`, default: `3600s`.

Example:

```bash
curl http://localhost:8080/system/cpu
curl "http://localhost:8080/logs?duration=2h"
```

---

## 🧠 How It Works

- A scheduler runs periodically (based on `pollRate`) to collect system metrics.
- Collected data is written to InfluxDB.
- `/system/*` endpoints fetch **live metrics** directly from the system.
- `/logs` endpoint retrieves **historical metrics** from InfluxDB and returns them as CSV.

---

## 🛠️ Dependencies

- [Gin](https://github.com/gin-gonic/gin) - REST framework
- [Zap](https://github.com/uber-go/zap) - Logging
- [gopsutil](https://github.com/shirou/gopsutil) - System metrics
- [InfluxDB Client](https://github.com/influxdata/influxdb-client-go) - Time-series storage

---

## 🚀 Future Enhancements

- Add support for weather/location data
- Dockerize the service
- Add authentication for endpoints
- Improve log filtering and CSV export

---

