package measure

import (
	"fmt"
	"math"
	"time"
)

type Measure struct {
	Cpu     cpu     `json:"cpu"`
	Disk    disk    `json:"disk"`
	Memory  memory  `json:"memory"`
	Network network `json:"network"`
	System  system  `json:"system"`
}

type cpu struct {
	Usage int64 `json:"usage"`
}

type disk struct {
	Path  string  `json:"path"`
	All   float64 `json:"all"`
	Used  float64 `json:"used"`
	Avail float64 `json:"avail"`
	Usage float64 `json:"usage"`
}

type memory struct {
	Total  float64 `json:"total"`
	Used   float64 `json:"used"`
	Usage  float64 `json:"usage"`
	Cached float64 `json:"cached"`
}

type network struct {
	Iface     string  `json:"iface"`
	IPaddress string  `json:"ipAddress"`
	RxUsage   float64 `json:"rxUsage"`
	TxUsage   float64 `json:"txUsage"`
	DailyTime time.Time
}

type system struct {
	OsRelease string `json:"osRelease"`
	Uptime    string `json:"uptime"`
}

const (
	B  float64 = 1
	KB float64 = 1024 * B
	MB float64 = 1024 * KB
	GB float64 = 1024 * MB
)

func (m *Measure) Show() {
	fmt.Println("Uptime:\t\t", m.System.Uptime)
	fmt.Println("Os-Release:\t", m.System.OsRelease)
	fmt.Printf("Cpu:\t\t%d%%\n", m.Cpu.Usage)
	fmt.Printf("Memory.Total:\t%0.2fGB\n", math.Round(m.Memory.Total/MB))
	fmt.Printf("Memory.Used:\t%0.2fGB\n", m.Memory.Used/MB) // 정확한 공식 찾아서 수정
	fmt.Printf("Memory.Usage:\t%0.2f%%\n", m.Memory.Usage)
	fmt.Printf("Memory.Cached:\t%0.2fGB\n", m.Memory.Cached)
	fmt.Printf("Disk.All:\t%0.2fGB\n", m.Disk.All)
	fmt.Printf("Disk.Avail:\t%0.2fGB\n", m.Disk.Avail)
	fmt.Printf("Disk.Used:\t%0.2fGB\n", m.Disk.Used)
	fmt.Printf("Disk.Usage:\t%0.2f%%\n", m.Disk.Usage)
	fmt.Printf("%s\t\t%s\n", m.Network.Iface, m.Network.IPaddress)
	fmt.Printf("Rx:\t\t%0.2fKB\n", m.Network.RxUsage)
	fmt.Printf("Tx:\t\t%0.2fKB\n", m.Network.TxUsage)
}
