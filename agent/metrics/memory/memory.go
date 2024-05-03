package memory

import (
	"bufio"
	"helloServer/agent/metrics"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func init() {
	if _, err := os.Stat("/proc/meminfo"); os.IsNotExist(err) {
		log.Println("'/proc/meminfo' is not exist")
		os.Exit(1)
	}
}

type metric struct{}

func New() *metric {
	return &metric{}
}

// total  Total installed memory (MemTotal and SwapTotal in /proc/meminfo)
// used   Used memory (calculated as total - free - buffers - cache)
// free   Unused memory (MemFree and SwapFree in /proc/meminfo)
// avail  = free - reserved filesystem blocks(for root)
func (mt *metric) Process(measure *metrics.Measure) error {
	var err error
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return errors.Wrap(err, "/proc/meminfo open error")
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	/*
		MemTotal:        8089672 kB
		MemFree:         4679572 kB
		MemAvailable:    6354024 kB
		Buffers:          586580 kB
		Cached:          1174896 kB
	*/

	var memoryTotal, memoryAvailable, memoryCached float64
	for i := 0; i < 5; i++ {
		scanner.Scan()
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		switch fields[0] {
		case "MemTotal:":
			value := fields[len(fields)-2]

			memoryTotal, err = strconv.ParseFloat(value, 64)
			if err != nil {
				return errors.Wrap(err, "Failed to convert (string -> int)")
			}
		case "MemAvailable:":
			value := fields[len(fields)-2]

			memoryAvailable, err = strconv.ParseFloat(value, 64)
			if err != nil {
				return errors.Wrap(err, "Failed to convert (string -> int)")
			}
		case "Cached:":
			value := fields[len(fields)-2]

			memoryCached, err = strconv.ParseFloat(value, 64)
			if err != nil {
				return errors.Wrap(err, "Failed to convert (string -> int)")
			}
		default:
			continue
		}
	}

	if scanner.Err() != nil {
		return errors.Wrap(scanner.Err(), "Failed to scan /proc/meminfo")
	}

	// KB 단위
	measure.Memory.Total = memoryTotal
	measure.Memory.Used = memoryTotal - memoryAvailable
	measure.Memory.Usage = ((memoryTotal - memoryAvailable) / memoryTotal) * float64(100)
	measure.Memory.Cached = memoryCached / metrics.MB

	return nil
}

func (mt *metric) Once(measure *metrics.Measure) error {
	return nil
}
