# Hello, Server!

Go 기반 Linux 서버 모니터링 에이전트입니다. Agent가 주기적으로 CPU, memory, disk, network, system metric을 수집하고, Fiber 기반 HTTP server가 최신 snapshot을 REST API로 제공합니다.

## 특징

- Linux `/proc`와 `statfs` 기반 metric 수집
- 기본 수집 주기 1초
- `/api/period`로 수집 주기 변경
- `/proc/stat`의 `total/idle` delta 기반 CPU 사용률 계산
- `eth0` 고정 없이 loopback이 아닌 활성 네트워크 인터페이스 자동 선택
- mutex로 보호된 metric snapshot cache
- `context.Context`와 `time.Ticker` 기반 agent lifecycle

## 실행

```bash
go run .
```

기본 포트는 `9227`입니다.

```bash
curl http://localhost:9227/health
curl http://localhost:9227/api/metrics/all
```

## REST API

| Method | URL | 설명 |
|:--:|:---|:---|
| GET | `/health` | health check |
| GET | `/api/metrics/all` | 전체 metric snapshot |
| GET | `/api/metrics/cpu` | CPU metric |
| GET | `/api/metrics/disk` | Disk metric |
| GET | `/api/metrics/memory` | Memory metric |
| GET | `/api/metrics/network` | Network metric |
| GET | `/api/metrics/system` | System metric |
| GET | `/api/process` | `/proc` 기반 process 목록 |
| PUT | `/api/period` | agent 수집 주기 변경 |

## Response 예시

`GET /api/metrics/all`

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
      "DailyTime": "2026-09-02T00:00:00Z"
    },
    "system": {
      "osRelease": "Ubuntu 20.04.6 LTS",
      "uptime": "0 days 2 hour"
    }
  },
  "elapse": "1.6µs"
}
```

`PUT /api/period`

```bash
curl -X PUT http://localhost:9227/api/period \
  -H 'Content-Type: application/json' \
  -d '{"period":3}'
```

```json
{
  "success": true,
  "reason": "success",
  "elapse": "38.6µs"
}
```

`period` 값은 초 단위이며 1 이상이어야 합니다.

## 저장 파일

CPU/network CSV는 `store/` 디렉터리 아래에 생성됩니다.

- `store/metric_cpu.csv`
- `store/metric_network.csv`

## 테스트

```bash
go test ./...
go test -race ./...
go vet ./...
```

홈 디렉터리의 Go build cache에 쓸 수 없는 환경에서는 다음처럼 cache 위치를 지정할 수 있습니다.

```bash
GOCACHE=/tmp/helloserver-gocache go test ./...
```

`test/client_test.go`는 오래된 TCP client 테스트라 기본적으로 skip됩니다. 필요한 경우 `HELLOSERVER_TEST_ADDR` 환경변수에 테스트 대상 주소를 넣어 실행할 수 있습니다.

## 지원 환경

Linux의 `/proc` 파일 시스템을 전제로 합니다. macOS나 Windows native 환경에서는 metric 수집이 정상 동작하지 않을 수 있습니다.
