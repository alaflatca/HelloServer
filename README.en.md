# Hello, Server!

[한국어](README.md)

A Go-based Linux server monitoring agent.

The agent periodically collects system metrics, and the Fiber HTTP server exposes the latest snapshot through a REST API.

## Features

- Collects CPU, memory, disk, network, and system metrics
- Calculates CPU usage from `/proc/stat`
- Selects an active non-loopback network interface automatically
- Provides an API for changing the collection period
- Uses a thread-safe metric snapshot cache
- Manages the collector lifecycle with `context.Context` and `time.Ticker`

## Run

```bash
go run .
```

The default port is `9227`.

```bash
curl http://localhost:9227/health
curl http://localhost:9227/api/metrics
```

## API

| Method | URL | Description |
|:--:|:---|:---|
| GET | `/health` | Health check |
| GET | `/api/metrics` | Get all metrics |
| GET | `/api/metrics/cpu` | Get CPU metrics |
| GET | `/api/metrics/disk` | Get disk metrics |
| GET | `/api/metrics/memory` | Get memory metrics |
| GET | `/api/metrics/network` | Get network metrics |
| GET | `/api/metrics/system` | Get system metrics |
| GET | `/api/processes` | Get process list |
| PUT | `/api/agent/period` | Change collection period |

The collection period is measured in seconds.

```bash
curl -X PUT http://localhost:9227/api/agent/period \
  -H 'Content-Type: application/json' \
  -d '{"period":3}'
```

## Response Examples

<details>
<summary>GET /health</summary>

```text
200 OK
```

</details>

<details>
<summary>GET /api/metrics</summary>

```json
{
  "success": true,
  "reason": "success",
  "data": {
    "cpu": {
      "usage": 4,
      "modelName": "Intel(R) Core(TM) CPU"
    },
    "disk": {
      "path": "/",
      "all": 250.92,
      "used": 58.37,
      "avail": 192.54,
      "usage": 23.26
    },
    "memory": {
      "total": 8089676,
      "used": 1416280,
      "usage": 17.5,
      "cached": 0.88
    },
    "network": {
      "iface": "ens5",
      "ipAddress": "192.168.0.3",
      "rxUsage": 12.5,
      "txUsage": 3.1,
      "DailyTime": "2026-09-04T00:00:00Z"
    },
    "system": {
      "osRelease": "Ubuntu 20.04.6 LTS",
      "uptime": "0 days 2 hour"
    }
  },
  "elapse": "1.6µs"
}
```

</details>

<details>
<summary>GET /api/metrics/cpu</summary>

```json
{
  "success": true,
  "reason": "success",
  "data": {
    "usage": 4,
    "modelName": "Intel(R) Core(TM) CPU"
  },
  "elapse": "1.2µs"
}
```

</details>

<details>
<summary>GET /api/metrics/disk</summary>

```json
{
  "success": true,
  "reason": "success",
  "data": {
    "path": "/",
    "all": 250.92,
    "used": 58.37,
    "avail": 192.54,
    "usage": 23.26
  },
  "elapse": "2.1µs"
}
```

</details>

<details>
<summary>GET /api/metrics/memory</summary>

```json
{
  "success": true,
  "reason": "success",
  "data": {
    "total": 8089676,
    "used": 1416280,
    "usage": 17.5,
    "cached": 0.88
  },
  "elapse": "1.9µs"
}
```

</details>

<details>
<summary>GET /api/metrics/network</summary>

```json
{
  "success": true,
  "reason": "success",
  "data": {
    "iface": "ens5",
    "ipAddress": "192.168.0.3",
    "rxUsage": 12.5,
    "txUsage": 3.1,
    "DailyTime": "2026-09-04T00:00:00Z"
  },
  "elapse": "2.6µs"
}
```

</details>

<details>
<summary>GET /api/metrics/system</summary>

```json
{
  "success": true,
  "reason": "success",
  "data": {
    "osRelease": "Ubuntu 20.04.6 LTS",
    "uptime": "0 days 2 hour"
  },
  "elapse": "2.3µs"
}
```

</details>

<details>
<summary>GET /api/processes</summary>

```json
{
  "success": true,
  "reason": "success",
  "data": [
    {
      "uid": "root",
      "pid": "1",
      "ppid": "0",
      "time": "",
      "cmd": "/sbin/init"
    }
  ],
  "elapse": "3.5ms"
}
```

</details>

<details>
<summary>PUT /api/agent/period</summary>

Request

```json
{
  "period": 3
}
```

Response

```json
{
  "success": true,
  "reason": "success",
  "elapse": "38.6µs"
}
```

</details>

## Test

```bash
go test ./...
go test -race ./...
go vet ./...
```

If the Go build cache is not writable, set a cache path explicitly.

```bash
GOCACHE=/tmp/helloserver-gocache go test ./...
```

## Notes

- This project targets Linux environments with `/proc`.
- CSV files are created under `store/`.
- `test/client_test.go` is a legacy TCP client test and is skipped by default.
