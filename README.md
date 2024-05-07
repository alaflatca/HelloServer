# Hello, Server!
**Simple Server Monitoring Agent + Web server**


## structure

### Agent
 : Collecting data from different types periodically ( cpu, memory, disk, network, system )
 
### Server
 : Web servers that provide data collected by agents

<br/>

## REST API

|Method | URL | 의미 |
|:--:|:---|:---|
|GET| `/api/metrics/all` | All Info |
|GET| `/api/metrics/cpu` | CPU Info |
|GET| `/api/metrics/disk` | Disk Info |
|GET| `/api/metrics/memory` | Memory Info |
|GET| `/api/metrics/network` | Network Info |
|GET| `/api/metrics/system` | System Info |
|PUT| `/api/period` |Agent period modify |

<br/>
<br/>


`/api/metrics/all`
<details>
 <summary> response </summary>
 
```json
{
    "success": true,
    "reason": "success",
    "data": {
        "cpu": {
            "usage": 4
        },
        "disk": {
            "path": "",
            "all": 250.92389297485352,
            "used": 58.37794494628906,
            "avail": 192.54594802856445,
            "usage": 23.265199760048137
        },
        "memory": {
            "total": 8089676,
            "used": 1416280,
            "usage": 17.507252453621135,
            "cached": 0.8846664428710938
        },
        "network": {
            "iface": "eth0:",
            "ipAddress": "172.30.148.16",
            "rxUsage": 0,
            "txUsage": 0,
            "DailyTime": "2024-05-07T00:00:00Z"
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
<br/>


`/api/metrics/cpu`
<details>
 <summary> response </summary>
 
```json
{
    "success": true,
    "reason": "success",
    "data": {
        "usage": 1
    },
    "elapse": "1.2µs"
}
```
</details>
<br/>


`/api/metrics/disk`
<details>
 <summary> response </summary>
 
 ```json
{
    "success": true,
    "reason": "success",
    "data": {
        "path": "",
        "all": 250.92389297485352,
        "used": 58.37795639038086,
        "avail": 192.54593658447266,
        "usage": 23.265204320830158
    },
    "elapse": "4.2µs"
}
```

</details>
<br/>


`/api/metrics/memory`
<details>
 <summary> response </summary>
 
 ```json
{
    "success": true,
    "reason": "success",
    "data": {
        "total": 8089676,
        "used": 1420448,
        "usage": 17.55877491261702,
        "cached": 0.8846778869628906
    },
    "elapse": "1.9µs"
}
 ```

</details>
<br/>


`/api/metrics/network`
<details>
<summary> response </summary>
 
 ```json
{
    "success": true,
    "reason": "success",
    "data": {
        "iface": "eth0:",
        "ipAddress": "192.168.0.3",
        "rxUsage": 0,
        "txUsage": 0,
        "DailyTime": "2024-05-07T00:00:00Z"
    },
    "elapse": "2.6µs"
}
```

</details>
<br/>


`/api/metrics/system`
<details>
 <summary> response </summary>
 
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
<br/>

`/api/period`
<details>
 <summary> request </summary>
 
```json
{
  "period": 3    // 3 second
}
```
</details>

<details>
 <summary> response </summary>
 
```json
{
    "success": true,
    "reason": "success",
    "elapse": "38.6µs"
}
```

</details>
<br/>
