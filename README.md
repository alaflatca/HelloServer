# Hello, Server!

Agent 
 : Collecting multiple data
 
Server
 : Rest api



REST API

- /api/metrics/all
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


- /api/metrics/cpu
{
    "success": true,
    "reason": "success",
    "data": {
        "usage": 1
    },
    "elapse": "1.2µs"
}


- /api/metrics/memory
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


- /api/metrics/disk
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


- /api/metrics/network
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


- /api/metrics/system
{
    "success": true,
    "reason": "success",
    "data": {
        "osRelease": "Ubuntu 20.04.6 LTS",
        "uptime": "0 days 2 hour"
    },
    "elapse": "2.3µs"
}