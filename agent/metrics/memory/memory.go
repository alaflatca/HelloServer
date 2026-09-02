package memory

import (
	"bufio"
	"fmt"
	"helloServer/measure"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type metric struct{}

func New() *metric {
	return &metric{}
}

// total  Total installed memory (MemTotal and SwapTotal in /proc/meminfo)
// used   Used memory (calculated as total - free - buffers - cache)
// free   Unused memory (MemFree and SwapFree in /proc/meminfo)
// avail  = free - reserved filesystem blocks(for root)
func (mt *metric) Process(ms *measure.Measure) error {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return errors.Wrap(err, "/proc/meminfo open error")
	}
	defer f.Close()

	memoryTotal, memoryAvailable, memoryCached, err := parseMemInfo(f)
	if err != nil {
		return err
	}

	// KB 단위
	ms.Memory.Total = memoryTotal
	ms.Memory.Used = memoryTotal - memoryAvailable
	ms.Memory.Usage = ((memoryTotal - memoryAvailable) / memoryTotal) * float64(100)
	ms.Memory.Cached = memoryCached / measure.MB

	return nil
}

func parseMemInfo(r io.Reader) (float64, float64, float64, error) {
	scanner := bufio.NewScanner(r)

	var memoryTotal, memoryAvailable, memoryCached float64
	var hasTotal, hasAvailable, hasCached bool
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		switch fields[0] {
		case "MemTotal:":
			value, err := parseMemInfoValue(fields)
			if err != nil {
				return 0, 0, 0, err
			}
			memoryTotal = value
			hasTotal = true
		case "MemAvailable:":
			value, err := parseMemInfoValue(fields)
			if err != nil {
				return 0, 0, 0, err
			}
			memoryAvailable = value
			hasAvailable = true
		case "Cached:":
			value, err := parseMemInfoValue(fields)
			if err != nil {
				return 0, 0, 0, err
			}
			memoryCached = value
			hasCached = true
		}
	}

	if scanner.Err() != nil {
		return 0, 0, 0, errors.Wrap(scanner.Err(), "Failed to scan /proc/meminfo")
	}
	if !hasTotal || !hasAvailable || !hasCached {
		return 0, 0, 0, fmt.Errorf("missing meminfo fields: total=%t available=%t cached=%t", hasTotal, hasAvailable, hasCached)
	}
	if memoryTotal <= 0 {
		return 0, 0, 0, errors.New("MemTotal must be greater than zero")
	}

	return memoryTotal, memoryAvailable, memoryCached, nil
}

func parseMemInfoValue(fields []string) (float64, error) {
	value, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return 0, errors.Wrapf(err, "parse meminfo value %q", fields[1])
	}
	return value, nil
}

func (mt *metric) Once(measure *measure.Measure) error {
	return nil
}
